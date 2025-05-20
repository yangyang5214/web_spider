package pkg

import (
	"github.com/go-kratos/kratos/v2/errors"
	"os"
)

func FileExists(p string) bool {
	_, err := os.Stat(p)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true //ignore other errors
}
