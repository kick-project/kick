// Package ansicodes is a collection of ANSI escape sequences used to format text output.
package ansicodes

// Codes is given its own type for safe function signatures
type Codes string

// Colour codes interpretted by the terminal
// NOTE: all codes must be of the same length or they will throw off the field alignment of tabwriter
const (
	Reset                   Codes = "\x1b[0000m"
	None                          = "\x1b[0000m"
	Bold                          = "\x1b[0001m"
	Faint                         = "\x1b[0002m"
	Underline                     = "\x1b[0004m"
	BlackText                     = "\x1b[0030m"
	RedText                       = "\x1b[0031m"
	GreenText                     = "\x1b[0032m"
	YellowText                    = "\x1b[0033m"
	BlueText                      = "\x1b[0034m"
	MagentaText                   = "\x1b[0035m"
	CyanText                      = "\x1b[0036m"
	WhiteText                     = "\x1b[0037m"
	DefaultText                   = "\x1b[0039m"
	DarkGreyText                  = "\x1b[1;30m"
	BrightRedText                 = "\x1b[1;31m"
	BrightGreenText               = "\x1b[1;32m"
	BrightYellowText              = "\x1b[1;33m"
	BrightBlueText                = "\x1b[1;34m"
	BrightMagentaText             = "\x1b[1;35m"
	BrightCyanText                = "\x1b[1;36m"
	BrightWhiteText               = "\x1b[1;37m"
	BlackBackground               = "\x1b[0040m"
	RedBackground                 = "\x1b[0041m"
	GreenBackground               = "\x1b[0042m"
	YellowBackground              = "\x1b[0043m"
	BlueBackground                = "\x1b[0044m"
	MagentaBackground             = "\x1b[0045m"
	CyanBackground                = "\x1b[0046m"
	WhiteBackground               = "\x1b[0047m"
	BrightBlackBackground         = "\x1b[0100m"
	BrightRedBackground           = "\x1b[0101m"
	BrightGreenBackground         = "\x1b[0102m"
	BrightYellowBackground        = "\x1b[0103m"
	BrightBlueBackground          = "\x1b[0104m"
	BrightMagentaBackground       = "\x1b[0105m"
	BrightCyanBackground          = "\x1b[0106m"
	BrightWhiteBackground         = "\x1b[0107m"
)
