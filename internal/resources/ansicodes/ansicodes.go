// Package ansicodes is a collection of ANSI escape sequences used to format text output.
package ansicodes

// Codes is given its own type for safe function signatures
type Codes string

// Colour codes interpretted by the terminal
// NOTE: all codes must be of the same length or they will throw off the field alignment of tabwriter
const (
	Reset                   Codes = "\x1b[0000m"
	None                    Codes = "\x1b[0000m"
	Bold                    Codes = "\x1b[0001m"
	Faint                   Codes = "\x1b[0002m"
	Underline               Codes = "\x1b[0004m"
	BlackText               Codes = "\x1b[0030m"
	RedText                 Codes = "\x1b[0031m"
	GreenText               Codes = "\x1b[0032m"
	YellowText              Codes = "\x1b[0033m"
	BlueText                Codes = "\x1b[0034m"
	MagentaText             Codes = "\x1b[0035m"
	CyanText                Codes = "\x1b[0036m"
	WhiteText               Codes = "\x1b[0037m"
	DefaultText             Codes = "\x1b[0039m"
	DarkGreyText            Codes = "\x1b[1;30m"
	BrightRedText           Codes = "\x1b[1;31m"
	BrightGreenText         Codes = "\x1b[1;32m"
	BrightYellowText        Codes = "\x1b[1;33m"
	BrightBlueText          Codes = "\x1b[1;34m"
	BrightMagentaText       Codes = "\x1b[1;35m"
	BrightCyanText          Codes = "\x1b[1;36m"
	BrightWhiteText         Codes = "\x1b[1;37m"
	BlackBackground         Codes = "\x1b[0040m"
	RedBackground           Codes = "\x1b[0041m"
	GreenBackground         Codes = "\x1b[0042m"
	YellowBackground        Codes = "\x1b[0043m"
	BlueBackground          Codes = "\x1b[0044m"
	MagentaBackground       Codes = "\x1b[0045m"
	CyanBackground          Codes = "\x1b[0046m"
	WhiteBackground         Codes = "\x1b[0047m"
	BrightBlackBackground   Codes = "\x1b[0100m"
	BrightRedBackground     Codes = "\x1b[0101m"
	BrightGreenBackground   Codes = "\x1b[0102m"
	BrightYellowBackground  Codes = "\x1b[0103m"
	BrightBlueBackground    Codes = "\x1b[0104m"
	BrightMagentaBackground Codes = "\x1b[0105m"
	BrightCyanBackground    Codes = "\x1b[0106m"
	BrightWhiteBackground   Codes = "\x1b[0107m"
)
