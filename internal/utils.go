package internal

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	alphaDigits   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var src = rand.NewSource(time.Now().UnixNano())

// GenRandomString generate random string with fixed length
func GenRandomString(n int) string {
	allowedChars := alphaDigits
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(allowedChars) {
			sb.WriteByte(allowedChars[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// DeepEqual checks if both values are recursively equal
// used for testing purpose only
func DeepEqual(v1, v2 any) bool {
	if reflect.DeepEqual(v1, v2) {
		return true
	}

	bytesA, _ := json.Marshal(v1)
	bytesB, _ := json.Marshal(v2)
	if string(bytesA) == string(bytesB) {
		return true
	}
	var x1 any
	var x2 any
	_ = json.Unmarshal(bytesA, &x1)
	_ = json.Unmarshal(bytesB, &x2)
	if reflect.DeepEqual(x1, x2) {
		return true
	}

	map1, ok1 := x1.(map[string]any)
	map2, ok2 := x2.(map[string]any)
	if !ok1 || !ok2 {
		return false
	}

	for k, v1 := range map1 {
		v2, ok := map2[k]
		if !ok || !DeepEqual(v1, v2) {
			return false
		}
	}
	return true
}
