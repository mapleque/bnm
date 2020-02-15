package server

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GetSessionToken() string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	prefix := strconv.Itoa(rand.Intn(10000))
	surfix := strconv.Itoa(rand.Intn(10000))
	return Md5(fmt.Sprintf("%s%s%s", prefix, timestamp, surfix))
}

func Md5(src string) string {
	h := md5.New()
	h.Write([]byte(src))
	data := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return string(dst)
}
