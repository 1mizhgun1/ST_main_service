package handlers

import (
	"fmt"
	"net/http"

	"github.com/satori/uuid"

	"github.com/1mizhgun1/ST_main_service/api/consts"
	"github.com/1mizhgun1/ST_main_service/api/kafka"
	"github.com/1mizhgun1/ST_main_service/api/utils"
)

func HandleSend(w http.ResponseWriter, r *http.Request) {
	data := utils.SendRequest{}
	err := utils.GetRequestData(r, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	segments := utils.SplitData(data.Text, consts.SegmentSize)
	total := len(segments)

	for i, segment := range segments {
		go utils.SendCodeRequest(utils.CodeRequest{
			MessageId:     uuid.NewV4(),
			SegmentNumber: i + 1,
			TotalSegments: total,
			Username:      data.Username,
			SendTime:      data.SendTime,
			Data:          segment,
		})
	}
}

func HandleTransfer(w http.ResponseWriter, r *http.Request) {
	data := utils.CodeRequest{}
	err := utils.GetRequestData(r, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = kafka.PutSegmentToKafka(data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
