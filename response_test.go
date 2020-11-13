package httpadapt

import (
	"strconv"
	"testing"
)

func TestResponse(t *testing.T) {
	for i, c := range []struct{}{} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_ = c
		})
	}
}
