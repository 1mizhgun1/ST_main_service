package utils

import (
	"math"
)

func SplitData(data string, segmentSize int) []string {
	result := make([]string, 0)

	length := len(data)
	segmentCount := int(math.Ceil(float64(length) / float64(segmentSize)))

	for i := 0; i < segmentCount; i++ {
		result = append(result, data[i*segmentSize:min((i+1)*segmentSize, length)])
	}

	return result
}
