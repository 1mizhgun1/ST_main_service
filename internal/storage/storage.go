package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/satori/uuid"

	"github.com/1mizhgun1/ST_main_service/internal/consts"
	"github.com/1mizhgun1/ST_main_service/internal/utils"
)

type Message struct {
	Received int
	Total    int
	Last     time.Time
	Id       int
	Username string
	Data     []string
	SendTime time.Time
}

type Storage map[uuid.UUID]Message
type sendFunc func(body utils.ReceiveRequest)

var storage = Storage{}

func addMessage(request utils.CodeRequest) {
	storage[request.MessageId] = Message{
		Received: 0,
		Total:    request.TotalSegments,
		Last:     time.Now().UTC(),
		Id:       request.Id,
		Username: request.Username,
		Data:     make([]string, request.TotalSegments),
		SendTime: request.SendTime,
	}
}

func AddSegment(request utils.CodeRequest) {
	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	id := request.MessageId
	_, found := storage[id]
	if !found {
		addMessage(request)
	}

	message, _ := storage[id]
	message.Received++
	message.Last = time.Now().UTC()
	message.Data[request.SegmentNumber-1] = request.Data
	storage[id] = message
}

func getMessageText(id uuid.UUID) string {
	result := ""
	message, _ := storage[id]
	for _, segment := range message.Data {
		result += segment
	}
	return result
}

func ScanStorage(sender sendFunc) {
	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	payload := utils.ReceiveRequest{}
	for id, message := range storage {
		if message.Received == message.Total {
			payload = utils.ReceiveRequest{
				Id:       message.Id,
				Username: message.Username,
				Text:     getMessageText(id),
				SendTime: message.SendTime,
				Error:    "",
			}
			fmt.Printf("sent message: %+v\n", payload)
			go sender(payload)
			delete(storage, id)
		} else if time.Since(message.Last) > consts.KafkaReadPeriod+time.Second {
			payload = utils.ReceiveRequest{
				Id:       message.Id,
				Username: message.Username,
				Text:     "",
				SendTime: message.SendTime,
				Error:    consts.SegmentLostError,
			}
			fmt.Printf("sent error: %+v\n", payload)
			go sender(payload)
			delete(storage, id)
		}
	}
}
