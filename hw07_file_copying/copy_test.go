package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyErrors(t *testing.T) {
	t.Run("empty from path", func(t *testing.T) {
		err := Copy("", "/tmp/out.txt", 0, 0)
		require.ErrorIs(t, err, ErrEmptyPath)
	})

	t.Run("empty to path", func(t *testing.T) {
		err := Copy("input.txt", "", 0, 0)
		require.ErrorIs(t, err, ErrEmptyPath)
	})

	t.Run("negative offset", func(t *testing.T) {
		err := Copy("input.txt", "/tmp/out.txt", -1, 0)
		require.ErrorIs(t, err, ErrNegativeOffset)
	})

	t.Run("negative limit", func(t *testing.T) {
		err := Copy("input.txt", "/tmp/out.txt", 0, -1)
		require.ErrorIs(t, err, ErrNegativeLimit)
	})

	t.Run("not existing input file", func(t *testing.T) {
		err := Copy("testdata/input1.txt", "/tmp/out.txt", 0, 0)
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("not regular input file (directory)", func(t *testing.T) {
		tmpDir := t.TempDir()
		err := Copy(tmpDir, "/tmp/out.txt", 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("not regular input file (/dev/urandom)", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp/out.txt", 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("offset larger than file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/out.txt", 10000000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
}

func TestCopySuccess(t *testing.T) {
	tests := []struct {
		limit        int64
		offset       int64
		expectedFile string
	}{
		{offset: 0, limit: 0, expectedFile: "testdata/out_offset0_limit0.txt"},
		{offset: 0, limit: 10, expectedFile: "testdata/out_offset0_limit10.txt"},
		{offset: 0, limit: 1000, expectedFile: "testdata/out_offset0_limit1000.txt"},
		{offset: 0, limit: 10000, expectedFile: "testdata/out_offset0_limit10000.txt"},
		{offset: 100, limit: 1000, expectedFile: "testdata/out_offset100_limit1000.txt"},
		{offset: 6000, limit: 1000, expectedFile: "testdata/out_offset6000_limit1000.txt"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("copy file offset=%d limit=%d", tc.offset, tc.limit), func(t *testing.T) {
			tmpDir := t.TempDir()
			outFile := filepath.Join(tmpDir, "out.txt")
			defer os.Remove(outFile)

			err := Copy("testdata/input.txt", outFile, tc.offset, tc.limit)
			require.NoError(t, err)

			eq, err := equalFiles(tc.expectedFile, outFile)
			require.NoError(t, err)
			require.True(t, eq)
		})
	}
}

func equalFiles(file1, file2 string) (bool, error) {
	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	const bufSize = 1024
	buf1 := make([]byte, bufSize)
	buf2 := make([]byte, bufSize)

	for {
		n1, err1 := f1.Read(buf1)
		n2, err2 := f2.Read(buf2)

		if n1 != n2 || !bytes.Equal(buf1[:n1], buf2[:n2]) {
			return false, nil
		}

		if err1 == io.EOF && err2 == io.EOF {
			return true, nil
		}

		if err1 == io.EOF || err2 == io.EOF {
			return false, nil
		}

		if err1 != nil {
			return false, err1
		}
		if err2 != nil {
			return false, err2
		}
	}
}
