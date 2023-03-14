package logging

import (
	"log"
	"math"
	"strings"
	"testing"
)

func TestMinLogLevel(t *testing.T) {
	t.Parallel()

	buffer := &strings.Builder{}
	logger := log.New(buffer, "", 0)
	target := New(Warning, logger)

	target.Debug()
	target.Info()
	target.Warning()
	target.Error()

	if buffer.Len() == 0 {
		t.Error("valid log level never triggered")
	}

	if !strings.Contains(buffer.String(), "WARNING:") {
		t.Error("valid log level never triggered")
	}

	if !strings.Contains(buffer.String(), "ERROR:") {
		t.Error("valid log level never triggered")
	}
}

func TestToStr(t *testing.T) {
	t.Parallel()

	tableTst := []struct {
		lvl      LogLevel
		expected string
	}{
		{
			lvl:      Debug,
			expected: "DEBUG",
		},
		{
			lvl:      Info,
			expected: "INFO",
		},
		{
			lvl:      Warning,
			expected: "WARNING",
		},
		{
			lvl:      Error,
			expected: "ERROR",
		},
	}

	for _, test := range tableTst {
		t.Run("ToStr", func(t *testing.T) {
			result := ToStr(test.lvl)

			if result != test.expected {
				t.Error("invalid LogLevel")
			}
		})
	}
}

func TestToStrWithPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("the code did not panic")
		}
	}()

	ToStr(LogLevel(math.MaxInt32))
}

func TestToLogLevel(t *testing.T) {
	t.Parallel()

	tableTst := []struct {
		lvl      string
		expected LogLevel
	}{
		{
			lvl:      "  deBug ",
			expected: Debug,
		},
		{
			lvl:      "info",
			expected: Info,
		},
		{
			lvl:      " warNing",
			expected: Warning,
		},
		{
			lvl:      "error",
			expected: Error,
		},
	}

	for _, test := range tableTst {
		t.Run("ToLogLevel", func(t *testing.T) {
			result := ToLogLevel(test.lvl)

			if result != test.expected {
				t.Error("invalid LogLevel")
			}
		})
	}
}

func TestToLogLevelWithPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("the code did not panic")
		}
	}()

	ToLogLevel("random_string")
}
