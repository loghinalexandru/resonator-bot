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
		t.Fatal("Valid log level never triggered!")
	}

	if !strings.Contains(buffer.String(), "WARNING:") {
		t.Fatal("Valid log level never triggered!")
	}

	if !strings.Contains(buffer.String(), "ERROR:") {
		t.Fatal("Valid log level never triggered!")
	}
}

func TestToStr(t *testing.T) {
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
				t.Fatal("Invalid LogLevel!")
			}
		})
	}
}

func TestToStr_WithPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("The code did not panic")
		}
	}()

	ToStr(LogLevel(math.MaxInt32))
}

func TestToLogLevel(t *testing.T) {
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
				t.Fatal("Invalid LogLevel!")
			}
		})
	}
}

func TestToLogLevel_WithPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("The code did not panic")
		}
	}()

	ToLogLevel("random_string")
}
