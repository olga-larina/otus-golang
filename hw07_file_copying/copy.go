package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeOffset        = errors.New("negative offset")
	ErrNegativeLimit         = errors.New("negative limit")
	ErrEmptyPath             = errors.New("empty path")
)

const bufferSize = int64(1024)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrNegativeOffset
	}
	if limit < 0 {
		return ErrNegativeLimit
	}
	if fromPath == "" || toPath == "" {
		return ErrEmptyPath
	}

	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		err := sourceFile.Close()
		if err != nil {
			log.Printf("error closing source file %v", err)
		}
	}()

	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	destinationFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		err := destinationFile.Close()
		if err != nil {
			log.Printf("error closing destination file %v", err)
		}
	}()

	var remaining int64
	if limit > 0 && limit < fileInfo.Size()-offset {
		remaining = limit
	} else {
		remaining = fileInfo.Size() - offset
	}

	sourceFile.Seek(offset, io.SeekStart)

	buffer := make([]byte, bufferSize)
	totalCopied := int64(0)

	bar := pb.Start64(remaining)
	defer bar.Finish()

	for remaining > 0 {
		bytesToRead := min(bufferSize, remaining)
		bytesRead, readErr := sourceFile.Read(buffer[:bytesToRead])
		if readErr != nil && readErr != io.EOF {
			return readErr
		}

		if bytesRead != 0 {
			_, writeErr := destinationFile.Write(buffer[:bytesRead])
			if writeErr != nil {
				return writeErr
			}
		}

		totalCopied += int64(bytesRead)
		remaining -= int64(bytesRead)

		bar.Add(bytesRead)

		if readErr == io.EOF || bytesRead == 0 {
			break
		}
		// time.Sleep(time.Second) // для проверки прогресс-бара
	}

	return nil
}
