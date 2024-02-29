package api

import (
	"github.com/satori/uuid"
	"net/http"
	"sync"
	"time"
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
type sendFunc func(client *http.Client, body receiveRequest)

var storage = Storage{}

func addMessage(request codeRequest) {
	storage[request.MessageId] = Message{
		Received: 0,
		Total:    request.TotalSegments,
		Last:     time.Now().UTC(),
		Username: request.Username,
		Data:     make([]string, request.TotalSegments),
		SendTime: request.SendTime,
	}
}

func addSegment(request codeRequest) {
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

func scanStorage(sender sendFunc) {
	mu := &sync.Mutex{}
	mu.Lock()

	client := &http.Client{}
	for id, message := range storage {
		if message.Received == message.Total {
			go sender(client, receiveRequest{
				Username: message.Username,
				Text:     getMessageText(id),
				SendTime: message.SendTime,
				Error:    "",
			})
			delete(storage, id)
		} else if time.Since(message.Last) > kafkaReadPeriod {
			go sender(client, receiveRequest{
				Username: message.Username,
				Text:     "",
				SendTime: message.SendTime,
				Error:    segmentLostError,
			})
			delete(storage, id)
		}
	}

	mu.Unlock()
}
