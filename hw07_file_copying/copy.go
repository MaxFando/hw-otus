package main

import (
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
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

	_, err = fromFile.Seek(offset, 0)
	if err != nil {
		return errors.Wrapf(err, "can't seek file %s", fromPath)
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return errors.Wrapf(err, "can't create file %s", toPath)
	}
	defer toFile.Close()

	return copyWithProgressBar(bytesToCopy, fromFile, toFile, fromPath, toPath)
}

func copyWithProgressBar(bytesToCopy int64, fromFile *os.File, toFile *os.File, fromPath, toPath string) error {
	bar := pb.Start64(bytesToCopy)
	defer bar.Finish()

	barReader := bar.NewProxyReader(fromFile)
	_, errCopy := io.CopyN(toFile, barReader, bytesToCopy)
	if errCopy != nil {
		errRemove := os.Remove(toPath)
		if errRemove != nil {
			return errors.Wrapf(errRemove, "can't remove file %s", toPath)
		}
		return errors.Wrapf(errCopy, "can't copy file %s to %s", fromPath, toPath)
	}

	return nil
}
