package command

import (
	"regexp"
	"strings"

	gloo "github.com/gloo-foo/framework"
)

type command gloo.Inputs[string, flags]

func Sed(parameters ...any) gloo.Command {
	return command(gloo.Initialize[string, flags](parameters...))
}

func (p command) Executor() gloo.CommandExecutor {
	// Get script from first positional argument
	script := ""
	if len(p.Positional) > 0 {
		script = p.Positional[0]
	}

	return gloo.LineTransform(func(line string) (string, bool) {
		output := line

		// Parse sed command - support simple s/pattern/replacement/ syntax
		if strings.HasPrefix(script, "s/") || strings.HasPrefix(script, "s|") || strings.HasPrefix(script, "s,") {
			sep := script[1:2]
			parts := strings.Split(script[2:], sep)

			if len(parts) >= 2 {
				pattern := parts[0]
				replacement := parts[1]
				flags := ""
				if len(parts) > 2 {
					flags = parts[2]
				}

				// Compile regex
				re, compileErr := regexp.Compile(pattern)
				if compileErr != nil {
					return "", false
				}

				// Apply substitution
				if strings.Contains(flags, "g") {
					// Global replacement
					output = re.ReplaceAllString(line, replacement)
				} else {
					// Replace only first occurrence
					count := 0
					output = re.ReplaceAllStringFunc(line, func(match string) string {
						if count == 0 {
							count++
							return re.ReplaceAllString(match, replacement)
						}
						return match
					})
				}
			}
		} else if strings.HasPrefix(script, "d") {
			// Delete command
			return "", false
		} else if strings.HasPrefix(script, "p") {
			// Print command (print twice)
			return line + "\n" + line, true
		}

		return output, true
	}).Executor()
}
