package opt

// Boolean flag types with constants
type InPlaceFlag bool
const (
	InPlace   InPlaceFlag = true
	NoInPlace InPlaceFlag = false
)

type QuietFlag bool
const (
	Quiet   QuietFlag = true
	NoQuiet QuietFlag = false
)

type ExtendedRegexFlag bool
const (
	ExtendedRegex   ExtendedRegexFlag = true
	NoExtendedRegex ExtendedRegexFlag = false
)

// Custom types for parameters
type Expression string
type ScriptFile string

// Flags represents the configuration options for the sed command
type Flags struct {
	InPlace       InPlaceFlag       // Edit files in place
	Quiet         QuietFlag         // Suppress default output
	ExtendedRegex ExtendedRegexFlag // Use extended regular expressions
	Expression    Expression        // Sed expression to execute
	ScriptFile    ScriptFile        // File containing sed script
}

// Configure methods for the opt system
func (f InPlaceFlag) Configure(flags *Flags) { flags.InPlace = f }
func (f QuietFlag) Configure(flags *Flags) { flags.Quiet = f }
func (f ExtendedRegexFlag) Configure(flags *Flags) { flags.ExtendedRegex = f }
func (e Expression) Configure(flags *Flags) { flags.Expression = e }
func (s ScriptFile) Configure(flags *Flags) { flags.ScriptFile = s }
