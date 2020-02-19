package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
)

func main() {

	r, err := interp.New(
		interp.Dir("./messagedb"),
		interp.Env(expand.ListEnviron()),

		interp.OpenHandler(func(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
			if path == "/dev/null" {
				return nil, nil
			}
			return interp.DefaultOpenHandler()(ctx, path, flag, perm)
		}),

		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
	)

	fmt.Println(r, err)
}
