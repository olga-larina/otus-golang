package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmdSuccess(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	cmd := []string{
		"/bin/bash",
		filepath.Join(pwd, "testdata", "echo.sh"),
		"arg1=1",
		"arg2=2",
	}
	env := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: "", NeedRemove: false},
		"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{NeedRemove: true},
	}
	os.Setenv("UNSET", "test")
	os.Setenv("ADDED", "from original env")

	// исходные потоки ввода / вывода
	origStdout := os.Stdout
	origStderr := os.Stderr

	// перехват вывода
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	// восстановление потоков
	defer func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()

	returnCode := RunCmd(cmd, env)

	w.Close()

	require.Equal(t, 0, returnCode)

	var actualOutput bytes.Buffer
	io.Copy(&actualOutput, r)

	expectedOutput := `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2
`
	require.Equal(t, expectedOutput, actualOutput.String())
}

func TestRunCmdErrors(t *testing.T) {
	t.Run("not existing command", func(t *testing.T) {
		returnCode := RunCmd([]string{"notexists"}, Environment{})
		require.Equal(t, 1, returnCode)
	})
}
