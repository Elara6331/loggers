package loggers

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"time"

	"go.elara.ws/loggers/internal/buffer"
	"go.elara.ws/loggers/internal/colors"
)

// Pretty is an slog handler with a human-readable output
type Pretty struct {
	mtx *sync.Mutex
	goa []groupOrAttr

	// Colorize indicates whether colors will be used
	// in log output.
	Colorize bool

	// Out is where log output will be written.
	Out io.Writer

	Options
}

// NewPretty creates and returns a new [Pretty] handler.
// If opts doesn't specify a time format, [time.Kitchen] will be used.
func NewPretty(wr io.Writer, opts Options) *Pretty {
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.Kitchen
	}

	return &Pretty{
		mtx:      &sync.Mutex{},
		Colorize: opts.ForceColors || isTerm(wr),
		Out:      wr,
		Options:  opts,
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (p *Pretty) Enabled(_ context.Context, level slog.Level) bool {
	return level >= p.Level
}

// WithGroup returns a new [Pretty] handler with the given group name.
func (p *Pretty) WithGroup(name string) slog.Handler {
	return &Pretty{
		mtx:      p.mtx,
		Out:      p.Out,
		Options:  p.Options,
		Colorize: p.Colorize,
		goa:      append(p.goa, groupOrAttr{group: name}),
	}
}

// WithAttrs returns a new [Pretty] handler whose attributes
// consists of p's attributes followed by attrs.
func (p *Pretty) WithAttrs(attrs []slog.Attr) slog.Handler {
	goa := make([]groupOrAttr, len(attrs))
	for i, attr := range attrs {
		goa[i] = groupOrAttr{attr: attr}
	}
	return &Pretty{
		mtx:      p.mtx,
		Out:      p.Out,
		Options:  p.Options,
		Colorize: p.Colorize,
		goa:      append(p.goa, goa...),
	}
}

// Handle formats the given [log/slog.Record] as a human-readable string on a single line.
func (p *Pretty) Handle(_ context.Context, rec slog.Record) error {
	buf := buffer.Alloc()
	defer buffer.Free(buf)

	colors.WriteString(p.Colorize, buf, colors.Grey, rec.Time.Format(p.TimeFormat))
	buf.WriteByte(' ')

	switch rec.Level {
	case slog.LevelInfo:
		colors.WriteString(p.Colorize, buf, colors.Green, "INF")
	case slog.LevelError:
		colors.WriteString(p.Colorize, buf, colors.Red, "ERR")
	case slog.LevelWarn:
		colors.WriteString(p.Colorize, buf, colors.Yellow, "WRN")
	case slog.LevelDebug:
		colors.WriteString(p.Colorize, buf, colors.Magenta, "DBG")
	}
	buf.WriteByte(' ')

	buf.WriteString(rec.Message)

	lastGroup := ""
	for _, goa := range p.goa {
		switch {
		case goa.group == "":
			if lastGroup != "" {
				goa.attr.Key = lastGroup + goa.attr.Key
			}
			buf.WriteByte(' ')
			writeAttr(p.Colorize, buf, goa.attr)
		default:
			lastGroup += goa.group + "."
		}
	}

	if rec.NumAttrs() > 0 {
		rec.Attrs(func(attr slog.Attr) bool {
			if lastGroup != "" {
				attr.Key = lastGroup + attr.Key
			}
			buf.WriteByte(' ')
			writeAttr(p.Colorize, buf, attr)
			return true
		})
	}

	if p.ShowCaller {
		writeCaller(p.Colorize, buf, rec.PC)
	}

	buf.WriteByte('\n')

	p.mtx.Lock()
	defer p.mtx.Unlock()

	_, err := buf.WriteTo(p.Out)
	return err
}
