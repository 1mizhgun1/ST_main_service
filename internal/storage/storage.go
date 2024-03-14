package storage

import (
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
		Username: request.Username,
		Data:     make([]string, request.TotalSegments),
		SendTime: request.SendTime,
	}
}

func AddSegment(request utils.CodeRequest) {
	mu := &sync.Mutex{}
	mu.Lock()

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

	mu.Unlock()
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

	for id, message := range storage {
		if message.Received == message.Total {
			go sender(utils.ReceiveRequest{
				Username: message.Username,
				Text:     getMessageText(id),
				SendTime: message.SendTime,
				Error:    "",
			})
			delete(storage, id)
		} else if time.Since(message.Last) > consts.KafkaReadPeriod+time.Second {
			go sender(utils.ReceiveRequest{
				Username: message.Username,
				Text:     "",
				SendTime: message.SendTime,
				Error:    consts.SegmentLostError,
			})
			delete(storage, id)
		}
	}

	mu.Unlock()
}
