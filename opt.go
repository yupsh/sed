package command

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

type Expression string
type ScriptFile string

type flags struct {
	InPlace       InPlaceFlag
	Quiet         QuietFlag
	ExtendedRegex ExtendedRegexFlag
	Expression    Expression
	ScriptFile    ScriptFile
}

func (f InPlaceFlag) Configure(flags *flags)       { flags.InPlace = f }
func (f QuietFlag) Configure(flags *flags)         { flags.Quiet = f }
func (f ExtendedRegexFlag) Configure(flags *flags) { flags.ExtendedRegex = f }
func (e Expression) Configure(flags *flags)        { flags.Expression = e }
func (s ScriptFile) Configure(flags *flags)        { flags.ScriptFile = s }
