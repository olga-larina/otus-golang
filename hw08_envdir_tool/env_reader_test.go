package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDirSuccess(t *testing.T) {
	dir := filepath.Join("testdata", "env")
	expected := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: "", NeedRemove: false},
		"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{NeedRemove: true},
	}
	env, err := ReadDir(dir)
	require.NoError(t, err)
	require.Equal(t, expected, env)
}

func TestReadDirErrors(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		_, err := ReadDir("")
		require.ErrorIs(t, err, ErrEmptyPath)
	})

	t.Run("not existing path", func(t *testing.T) {
		tmpDir := t.TempDir()
		_, err := ReadDir(filepath.Join(tmpDir, "test"))
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("file contains equal sign", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile, err := os.CreateTemp(tmpDir, "1=2")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = ReadDir(tmpDir)
		require.ErrorIs(t, err, ErrFileNameContainsEqualSign)
	})
}
