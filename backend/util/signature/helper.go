package signature

import "time"

var startTime time.Time

func init() {
	bangkokTime, _ := time.LoadLocation("Asia/Bangkok")
	startTime = time.Date(2001, 9, 07, 17, 00, 00, 00, bangkokTime)
}

func extractPathSlice(path string, depth uint32) []byte {
	if depth == 0 {
		return []byte(path[:1])
	}

	count := int(depth)
	for index := 1; index < len(path); index++ {
		if path[index] == '/' {
			count--
			if count <= 0 {
				return []byte(path[:index])
			}
		}
	}

	return []byte(path[:1])
}