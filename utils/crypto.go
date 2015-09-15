package utils

import (
	"crypto/rand"
	"github.com/cortunl/cortunl/constants"
	"github.com/dropbox/godropbox/errors"
	"math/big"
	mathrand "math/rand"
	"time"
)

func RandBytes(size int) (bytes []byte, err error) {
	bytes = make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		err = &constants.ReadError{
			errors.Wrap(err, "utils: Random read error"),
		}
		return
	}

	return
}

func seedRand() {
	n, err := rand.Int(rand.Reader, big.NewInt(9223372036854775806))
	if err != nil {
		mathrand.Seed(time.Now().UTC().UnixNano())
		return
	}

	mathrand.Seed(n.Int64())
	return
}

func init() {
	seedRand()
}
