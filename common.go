package loggers

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"go.elara.ws/loggers/internal/colors"
	"golang.org/x/term"
)

// Options represents common options for slog handlers
// in this package.
type Options struct {
	// TimeFormat represents the format of the timestamp
	// provided by the [Pretty] logger.
	TimeFormat string

	// Level is the log level above which the handler will
	// handle records.
	Level slog.Level

	// ShowCaller indicates whether the caller that created the
	// record should be provided in the output
	ShowCaller  bool
	
	// ForceColors prevents checking whether a handler outputs
	// to a valid tty and always enables colors if set to true.
	ForceColors bool
}

// groupOrAttr represents a group name or [log/slog.Attr].
type groupOrAttr struct {
	attr  slog.Attr
	group string
}

// writeGroup writes all the attributes in a group to buf.
func writeGroup(colorize bool, buf *bytes.Buffer, group slog.Attr) {
	attrs := group.Value.Group()
	for i, attr := range attrs {
		attr.Key = group.Key + "." + attr.Key
		writeAttr(colorize, buf, attr)
		if i < len(attrs)-1 {
			buf.WriteByte(' ')
		}
	}
}

// writeAttr writes a single attribute to buf.
//
// If the attribute is a [log/slog.Group], it calls [writeGroup].
// If the value of the attribute is an error, it will color it red.
func writeAttr(colorize bool, buf *bytes.Buffer, attr slog.Attr) {
	attr.Value = attr.Value.Resolve()

	if attr.Equal(slog.Attr{}) {
		return
	}

	if attr.Value.Kind() == slog.KindGroup {
		writeGroup(colorize, buf, attr)
		return
	}

	if _, ok := attr.Value.Any().(error); ok {
		colors.WriteString(colorize, buf, colors.Red, attr.Key+"=")
		colors.WriteCode(colorize, buf, colors.LightRed)
	} else {
		colors.WriteString(colorize, buf, colors.Cyan, attr.Key+"=")
		colors.WriteCode(colorize, buf, colors.White)
	}

	abuf := buf.AvailableBuffer()
	switch attr.Value.Kind() {
	case slog.KindInt64:
		abuf = strconv.AppendInt(abuf, attr.Value.Int64(), 10)
	case slog.KindUint64:
		abuf = strconv.AppendUint(abuf, attr.Value.Uint64(), 10)
	case slog.KindFloat64:
		abuf = strconv.AppendFloat(abuf, attr.Value.Float64(), 'g', -1, 64)
	case slog.KindBool:
		abuf = strconv.AppendBool(abuf, attr.Value.Bool())
	case slog.KindDuration:
		abuf = append(abuf, attr.Value.Duration().String()...)
	case slog.KindTime:
		abuf = attr.Value.Time().AppendFormat(abuf, time.RFC3339)
	default:
		abuf = strconv.AppendQuote(abuf, attr.Value.String())
	}
	buf.Write(abuf)
	colors.WriteCode(colorize, buf, colors.Reset)
}

// writeCaller extracts the caller from the given program counter
// and writes it to buf, enclosed in square brackets.
func writeCaller(colorize bool, buf *bytes.Buffer, pc uintptr) {
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	text := "[" + filepath.Base(frame.File) + ":" + strconv.Itoa(frame.Line) + "]"
	buf.WriteByte(' ')
	colors.WriteString(colorize, buf, colors.Bold+colors.LightBlue, text)
}

// isTerm checks if wr corresponds to a valid tty.
func isTerm(wr io.Writer) bool {
	if fl, ok := wr.(*os.File); ok {
		return term.IsTerminal(int(fl.Fd()))
	}
	return false
}
