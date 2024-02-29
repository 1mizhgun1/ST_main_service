package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/satori/uuid"
)

type sendRequest struct {
	Username string    `json:"username"`
	Text     string    `json:"data"`
	SendTime time.Time `json:"send_time"`
}
type receiveRequest struct {
	Username string    `json:"username"`
	Text     string    `json:"data"`
	SendTime time.Time `json:"send_time"`
	Error    string    `json:"error"`
}

type codeRequest struct {
	MessageId     uuid.UUID `json:"message_id"`
	SegmentNumber int       `json:"segment_number"`
	TotalSegments int       `json:"total_segments"`
	Username      string    `json:"username"`
	SendTime      time.Time `json:"send_time"`
	Data          string    `json:"data"`
}

func getRequestData(r *http.Request, requestData interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		return err
	}

	return nil
}

func sendCodeRequest(client *http.Client, body codeRequest) {
	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", codeUrl, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(reqBody)))

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func sendReceiveRequest(client *http.Client, body receiveRequest) {
	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", receiveUrl, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(reqBody)))

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}
