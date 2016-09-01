package rand

import (
	"math/rand"
	"strconv"
)

func NumString(len_ int) string {

	if len_ < 1 {
		return ""
	}

	str := make([]byte, len_)
	for i := 0; i < len_; i++ {
		str[i] = strconv.Itoa(rand.Intn(10))[0]
	}

	return string(str)
}
