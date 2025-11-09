package command_test

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/sed"
)

func TestSed_Substitute(t *testing.T) {
	result := run.Command(command.Sed("s/a/x/")).
		WithStdinLines("abc").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"xbc"})
}

func TestSed_GlobalSubstitute(t *testing.T) {
	result := run.Command(command.Sed("s/a/x/g")).
		WithStdinLines("aaa").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"xxx"})
}

func TestSed_Delete(t *testing.T) {
	result := run.Command(command.Sed("d")).
		WithStdinLines("a", "b", "c").Run()
	assertion.NoError(t, result.Err)
	assertion.Empty(t, result.Stdout)
}

func TestSed_Print(t *testing.T) {
	result := run.Command(command.Sed("p")).
		WithStdinLines("a").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 2) // Line is printed twice
}

func TestSed_EmptyInput(t *testing.T) {
	result := run.Quick(command.Sed("s/a/x/"))
	assertion.NoError(t, result.Err)
	assertion.Empty(t, result.Stdout)
}

func TestSed_InputError(t *testing.T) {
	result := run.Command(command.Sed("s/a/x/")).
		WithStdinError(errors.New("read failed")).Run()
	assertion.ErrorContains(t, result.Err, "read failed")
}

