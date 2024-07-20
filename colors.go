package loggers

// Color represents an ANSI color escape code
type Color string

const (
	// Reset resets all text formatting
	Reset Color = "\033[0m"

	// Bold represents bold text
	Bold Color = "\033[1m"

	// Underline represents underlined text
	Underline Color = "\033[4m"

	// Reverse represents reversed text (background and foreground colors are swapped)
	Reverse Color = "\033[7m"

	// LightRed represents light red color
	LightRed Color = "\033[91m"

	// LightGreen represents light green color
	LightGreen Color = "\033[92m"

	// LightYellow represents light yellow color
	LightYellow Color = "\033[93m"

	// LightBlue represents light blue color
	LightBlue Color = "\033[94m"

	// LightMagenta represents light magenta color
	LightMagenta Color = "\033[95m"

	// LightCyan represents light cyan color
	LightCyan Color = "\033[96m"

	// Black represents black color
	Black Color = "\033[30m"

	// Red represents red color
	Red Color = "\033[31m"

	// Green represents green color
	Green Color = "\033[32m"

	// Yellow represents yellow color
	Yellow Color = "\033[33m"

	// Blue represents blue color
	Blue Color = "\033[34m"

	// Magenta represents magenta color
	Magenta Color = "\033[35m"

	// Cyan represents cyan color
	Cyan Color = "\033[36m"

	// White represents white color
	White Color = "\033[37m"
)
