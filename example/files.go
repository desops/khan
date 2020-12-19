package main

import (
	"fmt"

	"github.com/desops/khan"
)

func init() {
	for i := 0; i < 100; i++ {
		khan.Add(&khan.File{
			Path:    fmt.Sprintf("/tmp/file_%d", i),
			Content: fmt.Sprintf("Contents of file %d", i),
		})
	}
}
