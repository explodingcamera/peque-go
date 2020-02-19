package sh

import (
	"context"
	"fmt"
	"os"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func exec(command string) error {
	r, err := interp.New(
		interp.Dir("./messagedb"),
		interp.Env(expand.ListEnviron()),
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
	)

	reader := strings.NewReader(command)
	prog, err := syntax.NewParser().Parse(reader, "")
	if err != nil {
		return err
	}
	r.Reset()
	ctx := context.Background()
	err = r.Run(ctx, prog)

	fmt.Println(err, r)
	return nil
}
