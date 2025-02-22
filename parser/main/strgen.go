package main

import (
	"math/rand"

	"time"
	"unsafe"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrcUnsafe(n int, conf *checkParams, chanLet chan string) {
	for {
		time.Sleep(20 * time.Millisecond)
		b := make([]byte, n)
		for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
			if remain == 0 {
				cache, remain = src.Int63(), letterIdxMax
			}
			if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
				b[i] = letterBytes[idx]
				i--
			}
			cache >>= letterIdxBits
			remain--
		}
		go conf.Check(*(*string)(unsafe.Pointer(&b)), chanLet)
	}
}
