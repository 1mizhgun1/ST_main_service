package api

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func handleSend(w http.ResponseWriter, r *http.Request) {
	data := sendRequest{}
	err := getRequestData(r, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

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
	data := codeRequest{}
	err := getRequestData(r, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = putSegmentToKafka(data.Data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
