package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/satori/uuid"

	"github.com/1mizhgun1/ST_main_service/internal/consts"
)

type SendRequest struct {
	Id       int       `json:"id,omitempty"`
	Username string    `json:"username"`
	Text     string    `json:"data"`
	SendTime time.Time `json:"send_time"`
}

type ReceiveRequest struct {
	Id       int       `json:"id,omitempty"`
	Username string    `json:"username"`
	Text     string    `json:"data"`
	SendTime time.Time `json:"send_time"`
	Error    string    `json:"error"`
}

type CodeRequest struct {
	Id            int       `json:"id,omitempty"`
	MessageId     uuid.UUID `json:"message_id"`
	SegmentNumber int       `json:"segment_number"`
	TotalSegments int       `json:"total_segments"`
	Username      string    `json:"username"`
	SendTime      time.Time `json:"send_time"`
	Data          string    `json:"data"`
}

func GetRequestData(r *http.Request, requestData interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &requestData); err != nil {
		return err
	}

	return nil
}

func SendCodeRequest(body CodeRequest) {
	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", consts.CodeUrl, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func SendReceiveRequest(body ReceiveRequest) {
	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", consts.ReceiveUrl, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}
