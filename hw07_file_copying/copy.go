package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSameFile              = errors.New("source and target files are same")
	ErrNegativeOffset        = errors.New("offset value can't be negative")
	ErrNegativeLimit         = errors.New("limit value can't be negative")
	BuffLen                  = 1024
)

type infoStat struct {
	From os.FileInfo
	To   os.FileInfo
}

// Copy is the main function: copies content from fromPath to toPath at offset bytes for limit bytes.
func Copy(fromPath, toPath string, offset, limit int64) error {
	fs, err := checkInputData(fromPath, toPath, offset, limit)
	if err != nil {
		return err
	}

	fromFile, err := os.Open(fromPath) // open fromFile
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0o644) // open toFile
	if err != nil {
		return err
	}
	defer toFile.Close()

	buff := make([]byte, BuffLen)
	offsetCurrent := int64(0)

	pbCount := fs.From.Size()
	if limit > 0 && limit < pbCount {
		pbCount = limit
	}
	pb := NewProgressBar(pbCount)
	defer pb.Finish()

	for {
		readLen, err := fromFile.ReadAt(buff, offsetCurrent)
		offsetCurrent += int64(readLen)
		if limit > 0 && offsetCurrent >= limit { // for limited reading
			writeLen := limit - offsetCurrent + int64(readLen)
			pb.Add(writeLen)

			if _, err = toFile.Write(buff[:writeLen]); err != nil {
				return err
			}
			break
		}
		if err != nil {
			if errors.Is(err, io.EOF) { // eof - ok, reading is done
				pb.Add(int64(readLen))

				if _, err = toFile.Write(buff[:readLen]); err != nil {
					return err
				}
				break
			}
			return err
		}

		pb.Add(int64(readLen))
		if _, err = toFile.Write(buff); err != nil { // copy read buff to target file
			return err
		}
	}

	return nil
}

// checkInputData checks file paths for existing regular files and also checks offset and limit for correct values.
func checkInputData(fromPath, toPath string, offset, limit int64) (*infoStat, error) {
	fs := &infoStat{}

	if offset < 0 {
		return fs, ErrNegativeOffset
	}

	if limit < 0 {
		return fs, ErrNegativeLimit
	}

	fromFi, err := fileInfoFrom(fromPath, offset)
	if err != nil {
		return fs, err
	}
	fs.From = fromFi

	toFi, err := fileInfoTo(toPath)
	if err != nil {
		return fs, err
	}
	fs.To = toFi

	if os.SameFile(fromFi, toFi) {
		return fs, ErrSameFile
	}

	if toFi != nil && toFi.Size() > 0 {
		os.Truncate(toPath, 0)
	}

	return fs, nil
}

// fileInfoFrom is subfunc of checkInputData for fileFrom.
func fileInfoFrom(fromPath string, offset int64) (os.FileInfo, error) {
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return fileInfo, err
	}
	if fileInfo.Size() < offset {
		return fileInfo, ErrOffsetExceedsFileSize
	}
	if !fileInfo.Mode().IsRegular() {
		return fileInfo, ErrUnsupportedFile
	}
	return fileInfo, nil
}

// fileInfoFrom is subfunc of checkInputData for fileTo.
func fileInfoTo(toPath string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(toPath)
	if err != nil { // target file may not exist
		if os.IsNotExist(err) {
			return fileInfo, nil
		}
		return fileInfo, err
	}

	if !fileInfo.Mode().IsRegular() {
		return fileInfo, ErrUnsupportedFile
	}
	return fileInfo, nil
}
