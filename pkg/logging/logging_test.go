package logging

import (
	"log"
	"math"
	"strings"
	"testing"
)

func TestMinLogLevelLogsOnlyWarningAndAbove(t *testing.T) {
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

func TestToStrWithValidInput(t *testing.T) {
	t.Parallel()

	tableTst := []struct {
		lvl  LogLevel
		want string
	}{
		{
			lvl:  Debug,
			want: "DEBUG",
		},
		{
			lvl:  Info,
			want: "INFO",
		},
		{
			lvl:  Warning,
			want: "WARNING",
		},
		{
			lvl:  Error,
			want: "ERROR",
		},
	}

	for _, test := range tableTst {
		t.Run("ToStr", func(t *testing.T) {
			got := ToStr(test.lvl)

			if got != test.want {
				t.Error("invalid LogLevel")
			}
		})
	}
}

func TestToStrWithInvalidInputPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("the code did not panic")
		}
	}()

	ToStr(LogLevel(math.MaxInt32))
}

func TestToLogLevelWithValidInput(t *testing.T) {
	t.Parallel()

	tableTst := []struct {
		lvl  string
		want LogLevel
	}{
		{
			lvl:  "  deBug ",
			want: Debug,
		},
		{
			lvl:  "info",
			want: Info,
		},
		{
			lvl:  " warNing",
			want: Warning,
		},
		{
			lvl:  "error",
			want: Error,
		},
	}

	for _, test := range tableTst {
		t.Run("ToLogLevel", func(t *testing.T) {
			got := ToLogLevel(test.lvl)

			if got != test.want {
				t.Error("invalid LogLevel")
			}
		})
	}
}

func TestToLogLevelWithInvalidInputPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("the code did not panic")
		}
	}()

	ToLogLevel("random_string")
}
