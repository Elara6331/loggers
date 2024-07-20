package colors

import "io"

// Color represents an ANSI color escape code
type Color string

const (
	// Reset resets all text formatting
	Reset Color = "\x1b[0m"

	// Bold represents bold text
	Bold Color = "\x1b[1m"

	// Underline represents underlined text
	Underline Color = "\x1b[4m"

	// Reverse represents reversed text (background and foreground colors are swapped)
	Reverse Color = "\x1b[7m"

	// Grey represents the grey color
	Grey Color = "\x1b[38:5:240m"

	// LightRed represents the light red color
	LightRed Color = "\x1b[91m"

	// LightGreen represents the light green color
	LightGreen Color = "\x1b[92m"

	// LightYellow represents the light yellow color
	LightYellow Color = "\x1b[93m"

	// LightBlue represents the light blue color
	LightBlue Color = "\x1b[94m"

	// LightMagenta represents the light magenta color
	LightMagenta Color = "\x1b[95m"

	// LightCyan represents the light cyan color
	LightCyan Color = "\x1b[96m"

	// Black represents the black color
	Black Color = "\x1b[30m"

	// Red represents the red color
	Red Color = "\x1b[31m"

	// Green represents the green color
	Green Color = "\x1b[32m"

	// Yellow represents the yellow color
	Yellow Color = "\x1b[33m"

	// Blue represents the blue color
	Blue Color = "\x1b[34m"

	// Magenta represents the magenta color
	Magenta Color = "\x1b[35m"

	// Cyan represents the cyan color
	Cyan Color = "\x1b[36m"

	// White represents the white color
	White Color = "\x1b[38:5:15m"
)

// WriteCode writes a single col to w if uc is true
func WriteCode(uc bool, w io.Writer, col Color) (int, error) {
	if uc {
		return io.WriteString(w, string(col))
	}
	return 0, nil
}

// Write writes b to w. It prepends col and appends a reset if uc is true.
func Write(uc bool, w io.Writer, col Color, b []byte) (n int, err error) {
	if uc {
		nn, err := io.WriteString(w, string(col))
		if err != nil {
			return nn, err
		}
		n += nn
	}

	nn, err := w.Write(b)
	if err != nil {
		return nn, err
	}
	n += nn

	if uc {
		nn, err := io.WriteString(w, string(Reset))
		if err != nil {
			return nn, err
		}
		n += nn
	}

	return n, nil
}

// Write writes s to w. It prepends col and appends a reset if uc is true.
func WriteString(uc bool, w io.Writer, col Color, s string) (n int, err error) {
	if uc {
		nn, err := io.WriteString(w, string(col))
		if err != nil {
			return nn, err
		}
		n += nn
	}

	nn, err := io.WriteString(w, s)
	if err != nil {
		return nn, err
	}
	n += nn

	if uc {
		nn, err := io.WriteString(w, string(Reset))
		if err != nil {
			return nn, err
		}
		n += nn
	}

	return n, nil
}
