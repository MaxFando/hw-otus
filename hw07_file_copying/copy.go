package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o666)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrUnsupportedFile
		}
		return err
	}
	defer fromFile.Close()

	fileInfo, err := fromFile.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return ErrUnsupportedFile
	}

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	bytesToCopy := fileSize - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	log.Println("bytesToCopy:", bytesToCopy, "offset:", offset, "limit:", limit)

	_, err = fromFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	bar := pb.Start64(bytesToCopy)
	defer bar.Finish()

	barReader := bar.NewProxyReader(fromFile)
	_, err = io.CopyN(toFile, barReader, bytesToCopy)
	if err != nil && errors.Is(err, io.EOF) {
		_ = os.Remove(toPath)
		return err
	}

	return nil
}
