package api

import (
	"net/http"

	"github.com/google/uuid"
)

func handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data := sendRequest{}
	err := getRequestData(w, r, &data)
	if err != nil {
		return
	}

	segments := splitData(data.Text, segmentSize)
	total := len(segments)

	client := &http.Client{}
	for i, segment := range segments {
		go sendCodeRequest(client, codeRequest{
			MessageId:     uuid.New(),
			SegmentNumber: i + 1,
			TotalSegments: total,
			Username:      data.Username,
			SendTime:      data.SendTime,
			Data:          segment,
		})
	}
}

func handleTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data := codeRequest{}
	err := getRequestData(w, r, &data)
	if err != nil {
		return
	}

	putSegmentToKafka(data.Data)
}
