package api

import (
	"math"
)

func splitData(data string, segmentSize int) []string {
	result := []string{}

	length := len(data)
	segmentCount := int(math.Ceil(float64(length) / float64(segmentSize)))

	for i := 0; i < segmentCount; i++ {
		result = append(result, data[i*segmentSize:min((i+1)*segmentSize, length)])
	}

	return result
}
