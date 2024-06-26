package consts

import "time"

const (
	MyHost            = "192.168.123.166"
	CodeUrl           = "http://192.168.146.60:8080/code"
	ReceiveUrl        = "http://192.168.146.193:8081/receive"
	SegmentSize       = 100
	ReadTimeout       = 10
	WriteTimeout      = 10
	ReadHeaderTimeout = 10
	SegmentLostError  = "lost"
	KafkaAddr         = "localhost:29092"
	KafkaTopic        = "segments"
	KafkaReadPeriod   = 2 * time.Second
)
