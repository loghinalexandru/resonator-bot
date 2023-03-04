package logging

import (
	"log"
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

func TestStrLogLevel(t *testing.T) {
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
		t.Run("LogLevel", func(t *testing.T) {
			result := strLogLevel(test.lvl)

			if result != test.expected {
				t.Fatal("Invalid LogLevel!")
			}
		})
	}
}
