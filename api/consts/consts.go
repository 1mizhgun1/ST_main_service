package consts

import "time"

const (
	CodeUrl           = "http://127.0.0.1:8080/test_send"
	ReceiveUrl        = "http://127.0.0.1:8080/test_receive"
	SegmentSize       = 100
	ReadTimeout       = 10
	WriteTimeout      = 10
	ReadHeaderTimeout = 10
	SegmentLostError  = "lost"
	KafkaAddr         = "localhost:9092"
	KafkaTopic        = "segments"
	KafkaReadPeriod   = 2 * time.Second
)
