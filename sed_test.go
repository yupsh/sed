package sed_test

import (
	"context"
	"os"
	"strings"

	"github.com/yupsh/sed"
	"github.com/yupsh/sed/opt"
)

func ExampleSed() {
	ctx := context.Background()
	input := strings.NewReader("hello world\nhello universe\n")

	cmd := sed.Sed("s/hello/hi/")
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output: hi world
	// hi universe
}

func ExampleSed_withExpression() {
	ctx := context.Background()
	input := strings.NewReader("foo bar\nbaz foo\n")

	cmd := sed.Sed(opt.Expression("s/foo/FOO/"))
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output: FOO bar
	// baz FOO
}
