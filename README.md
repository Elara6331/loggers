# Loggers

Loggers is a collection of [`log/slog`](https://pkg.go.dev/log/slog) handlers that pretty print logs, similar to Zerolog's ConsoleWriter, but without sacrificing any performance.

# Example

```go
package main

import (
	"log/slog"
	"os"
	"strconv"

	"go.elara.ws/loggers"
)

const input = "hello"

func main() {
	log := slog.New(loggers.NewPretty(os.Stdout, loggers.Options{
		Level: slog.LevelDebug,
		ShowCaller: true,
	}))

	i, err := strconv.Atoi(input)
	if err != nil {
		log.Error(
			"Couldn't convert to integer",
			slog.String("input", input),
			slog.Any("error", err),
		)
		return
	}

	log.Info(
		"Converted to integer!",
		slog.Int("output", i),
	)
}
```