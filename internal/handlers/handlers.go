package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/satori/uuid"

	"github.com/1mizhgun1/ST_main_service/internal/consts"
	"github.com/1mizhgun1/ST_main_service/internal/kafka"
	"github.com/1mizhgun1/ST_main_service/internal/utils"
)

// HandleSend godoc
// @Summary		Send
// @Description	Send message
// @Tags 		message
// @ID			send
// @Accept		json
// @Param		payload 	body	utils.SendRequest	true	"message data"
// @Success		200
// @Failure		400
// @Router		/send [post]
func HandleSend(w http.ResponseWriter, r *http.Request) {
	data := utils.SendRequest{}
	if err := utils.GetRequestData(r, &data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)

	segments := utils.SplitData(data.Text, consts.SegmentSize)
	total := len(segments)

	id := uuid.NewV4()
	for i, segment := range segments {
		payload := utils.CodeRequest{
			Id:            data.Id,
			MessageId:     id,
			SegmentNumber: i + 1,
			TotalSegments: total,
			Username:      data.Username,
			SendTime:      data.SendTime,
			Data:          segment,
		}
		go utils.SendCodeRequest(payload)
		fmt.Printf("sent segment: %+v\n", payload)
	}
}

// HandleTransfer godoc
// @Summary		Transfer
// @Description	Transfer message
// @Tags 		message
// @ID			transfer
// @Accept		json
// @Param		payload 	body	utils.CodeRequest	true	"segment data"
// @Success		200
// @Failure		400
// @Router		/transfer [post]
func HandleTransfer(w http.ResponseWriter, r *http.Request) {
	data := utils.CodeRequest{}
	if err := utils.GetRequestData(r, &data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	if err := kafka.PutSegmentToKafka(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
