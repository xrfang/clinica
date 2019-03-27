package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	_G_HASH string
	_G_REVS string
	_BUILT_ string
)

func verinfo() string {
	self := filepath.Base(os.Args[0])
	return fmt.Sprintf("%s V%s.%s", self, _G_REVS, _G_HASH)
}
