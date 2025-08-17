package sed

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"

	localopt "github.com/yupsh/sed/opt"
)


// Flags represents the configuration options for the sed command
type Flags = localopt.Flags
// Command implementation
type command opt.Inputs[string, Flags]

// Sed creates a new sed command with the given parameters
func Sed(parameters ...any) yup.Command {
	return command(opt.Args[string, Flags](parameters...))
}

func (c command) Execute(ctx context.Context, input io.Reader, output, stderr io.Writer) error {
	// Check for cancellation before starting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Get expression from flags or first positional argument
	expression := string(c.Flags.Expression)
	if expression == "" && len(c.Positional) > 0 {
		expression = c.Positional[0]
	}

	if expression == "" {
		fmt.Fprintln(stderr, "sed: missing expression")
		return fmt.Errorf("missing expression")
	}

	// Parse sed expression (simplified)
	sedCmd, err := c.parseExpression(expression)
	if err != nil {
		fmt.Fprintf(stderr, "sed: %v\n", err)
		return err
	}

	// Process files or stdin
	files := c.Positional
	if string(c.Flags.Expression) != "" {
		files = c.Positional // All positional args are files
	} else if len(c.Positional) > 1 {
		files = c.Positional[1:] // Skip first arg (expression)
	}

	if len(files) == 0 {
		return c.processReader(ctx, input, output, sedCmd)
	}

	for _, filename := range files {
		// Check for cancellation before each file
		if err := yup.CheckContextCancellation(ctx); err != nil {
			return err
		}

		if err := c.processFile(ctx, filename, output, stderr, sedCmd); err != nil {
			fmt.Fprintf(stderr, "sed: %s: %v\n", filename, err)
		}
	}

	return nil
}

type SedCommand struct {
	Command string
	Pattern string
	Replacement string
}

func (c command) parseExpression(expr string) (*SedCommand, error) {
	// Simplified sed expression parsing
	// Real implementation would be much more complex

	if strings.HasPrefix(expr, "s/") {
		// Substitute command: s/pattern/replacement/flags
		parts := strings.Split(expr[2:], "/")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid substitute expression")
		}

		return &SedCommand{
			Command: "s",
			Pattern: parts[0],
			Replacement: parts[1],
		}, nil
	}

	return nil, fmt.Errorf("unsupported sed expression: %s", expr)
}

func (c command) processReader(ctx context.Context, reader io.Reader, output io.Writer, sedCmd *SedCommand) error {
	scanner := bufio.NewScanner(reader)

	for yup.ScanWithContext(ctx, scanner) {
		line := scanner.Text()
		result := c.applyCommand(ctx, line, sedCmd)

		if !bool(c.Flags.Quiet) {
			fmt.Fprintln(output, result)
		}
	}

	// Check if context was cancelled
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	return scanner.Err()
}

func (c command) processFile(ctx context.Context, filename string, output, stderr io.Writer, sedCmd *SedCommand) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if bool(c.Flags.InPlace) {
		// In-place editing would require temporary file handling
		// This is a simplified version
		fmt.Fprintf(stderr, "sed: in-place editing not implemented in this stub\n")
	}

	return c.processReader(ctx, file, output, sedCmd)
}

func (c command) applyCommand(ctx context.Context, line string, sedCmd *SedCommand) string {
	// Check for cancellation before processing complex operations
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return line // Return original line on cancellation
	}

	switch sedCmd.Command {
	case "s":
		// Simple substitution
		regex, err := regexp.Compile(sedCmd.Pattern)
		if err != nil {
			return line // Return original on error
		}
		return regex.ReplaceAllString(line, sedCmd.Replacement)
	default:
		return line
	}
}
