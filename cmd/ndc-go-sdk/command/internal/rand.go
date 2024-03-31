package internal

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	alphaDigits   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// GenRandomScalarValue generates random scalar value depending on its representation type
func GenRandomScalarValue(random *rand.Rand, name string, scalar *schema.ScalarType) any {
	switch name {
	case "UUID":
		return uuid.New().String()
	case "DateTime":
		return time.Now().Format(time.RFC3339)
	}

	switch s := scalar.Representation.Interface().(type) {
	case *schema.TypeRepresentationBoolean:
		return random.Intn(2) == 1
	case *schema.TypeRepresentationInteger:
		return random.Intn(math.MaxInt16)
	case *schema.TypeRepresentationNumber:
		return random.Float32() * (10 ^ 4)
	case *schema.TypeRepresentationString:
		return internal.GenRandomString(10)
	case *schema.TypeRepresentationEnum:
		return s.OneOf[rand.Intn(len(s.OneOf))]
	default:
		return nil
	}
}

// GenRandomString generate random string with fixed length
func GenRandomString(src *rand.Rand, n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(alphaDigits) {
			sb.WriteByte(alphaDigits[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
