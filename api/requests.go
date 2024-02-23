package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type sendRequest struct {
	Username string    `json:"username"`
	Text     string    `json:"text"`
	SendTime time.Time `json:"send_time"`
}

type codeRequest struct {
	MessageId     uuid.UUID `json:"message_id"`
	SegmentNumber int       `json:"segment_number"`
	TotalSegments int       `json:"total_segments"`
	Username      string    `json:"username"`
	SendTime      time.Time `json:"send_time"`
	Data          string    `json:"data"`
}

func getSendData(w http.ResponseWriter, r *http.Request) (sendRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return sendRequest{}, err
	}
	defer r.Body.Close()

	data := sendRequest{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return sendRequest{}, err
	}

	w.WriteHeader(http.StatusOK)
	return data, err
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
