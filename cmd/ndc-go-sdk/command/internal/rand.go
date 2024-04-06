package internal

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"encoding/base64"

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
	switch s := scalar.Representation.Interface().(type) {
	case *schema.TypeRepresentationBoolean:
		return random.Intn(2) == 1
	case *schema.TypeRepresentationInt8:
		return random.Intn(math.MaxInt8)
	case *schema.TypeRepresentationInt16:
		return random.Intn(math.MaxInt16)
	case *schema.TypeRepresentationInt32:
		return random.Intn(math.MaxInt32)
	case *schema.TypeRepresentationInt64:
		return strconv.FormatInt(random.Int63n(math.MaxInt64), 10)
	case *schema.TypeRepresentationFloat32:
		return random.Float32() * (10 ^ 4)
	case *schema.TypeRepresentationFloat64:
		return random.Float64() * (10 ^ 8)
	case *schema.TypeRepresentationString:
		return internal.GenRandomString(10)
	case *schema.TypeRepresentationBigDecimal:
		return fmt.Sprintf("%.2f", random.Float64()*(10^8))
	case *schema.TypeRepresentationDate:
		return time.Now().Format("2006-01-02")
	case *schema.TypeRepresentationTimestamp:
		return time.Now().Format("2006-01-02T15:04:05Z")
	case *schema.TypeRepresentationTimestampTZ:
		return time.Now().Format(time.RFC3339)
	case *schema.TypeRepresentationUUID:
		return uuid.NewString()
	case *schema.TypeRepresentationBytes:
		return base64.StdEncoding.EncodeToString([]byte(internal.GenRandomString(10)))
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
