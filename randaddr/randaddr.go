package randaddr

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandPort() int {
	return 10000 + int(r.Uint32()%20000)
}

func Resolve(s string) string {
	if strings.HasSuffix(s, ":rand") {
		s = strings.TrimSuffix(s, ":rand")
		s = fmt.Sprintf("%s:%d", s, RandPort())
	}
	return s
}

func Local() string {
	return fmt.Sprintf("localhost:%d", RandPort())
}
