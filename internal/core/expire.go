package core

import (
	"ngvanthanggit/internal/constant"
	"time"
)

func ActiveDeleteExpiredKeys() {
	for {
		expiredCount := 0
		sampleCountRemain := constant.ActiveExpireSampleSize
		for key, expiredTime := range dictStore.GetExpiredDictStore() {
			sampleCountRemain--
			if sampleCountRemain < 0 {
				break
			}
			// check if the key is expired
			if time.Now().UnixMilli() > int64(expiredTime) {
				expiredCount++
				dictStore.Del(key)
			}
		}

		// if expiredCount >= 10% of the sample -> repeat
		if float64(expiredCount)/float64(constant.ActiveExpireSampleSize) < constant.ActiveExpireThreshold {
			break
		}
	}
}
