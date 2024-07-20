package loggers

import (
	"context"
	"io"
	"log/slog"
	"sync"

	"go.elara.ws/loggers/internal/buffer"
	"go.elara.ws/loggers/internal/colors"
)

// CLI is an slog handler for command-line tools where users will view
// and read the logs throughout the application's execution.
type CLI struct {
	mtx *sync.Mutex
	goa []groupOrAttr

	// Colorize indicates whether colors will be used
	// in log output.
	Colorize bool

	// Out is where log output will be written.
	Out io.Writer

	Options
}

// NewPretty creates and returns a new [CLI] handler.
func NewCLI(wr io.Writer, opts Options) *CLI {
	return &CLI{
		mtx:      &sync.Mutex{},
		Colorize: opts.ForceColors || isTerm(wr),
		Out:      wr,
		Options:  opts,
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (c *CLI) Enabled(_ context.Context, level slog.Level) bool {
	return level >= c.Level
}

// WithGroup returns a new [CLI] handler with the given group name.
func (c *CLI) WithGroup(group string) slog.Handler {
	return &CLI{
		mtx:      c.mtx,
		Out:      c.Out,
		Options:  c.Options,
		Colorize: c.Colorize,
		goa:      append(c.goa, groupOrAttr{group: group}),
	}
}

// WithAttrs returns a new [CLI] handler whose attributes
// consists of c's attributes followed by attrs.
func (c *CLI) WithAttrs(attrs []slog.Attr) slog.Handler {
	goa := make([]groupOrAttr, len(attrs))
	for i, attr := range attrs {
		goa[i] = groupOrAttr{attr: attr}
	}
	return &CLI{
		mtx:      c.mtx,
		Out:      c.Out,
		Options:  c.Options,
		Colorize: c.Colorize,
		goa:      append(c.goa, goa...),
	}
}

// Handle formats the given [log/slog.Record] as a human-readable string on a single line.
func (c *CLI) Handle(_ context.Context, rec slog.Record) error {
	buf := buffer.Alloc()
	defer buffer.Free(buf)

	switch rec.Level {
	case slog.LevelInfo:
		colors.WriteString(c.Colorize, buf, colors.Green, "-->")
	case slog.LevelError:
		colors.WriteString(c.Colorize, buf, colors.Red, " ->")
	case slog.LevelWarn:
		colors.WriteString(c.Colorize, buf, colors.Yellow, " ->")
	case slog.LevelDebug:
		colors.WriteString(c.Colorize, buf, colors.Magenta, "[DBG]")
	}
	buf.WriteByte(' ')

	buf.WriteString(rec.Message)

	lastGroup := ""
	for _, goa := range c.goa {
		switch {
		case goa.group == "":
			if lastGroup != "" {
				goa.attr.Key = lastGroup + goa.attr.Key
			}
			buf.WriteByte(' ')
			writeAttr(c.Colorize, buf, goa.attr)
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
			writeAttr(c.Colorize, buf, attr)
			return true
		})
	}

	if c.ShowCaller {
		writeCaller(c.Colorize, buf, rec.PC)
	}

	buf.WriteByte('\n')

	c.mtx.Lock()
	defer c.mtx.Unlock()

	_, err := buf.WriteTo(c.Out)
	return err
}
