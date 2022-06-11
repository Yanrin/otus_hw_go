package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests check only by file size. Shall script checks content.

func TestCopyIncorrectInputData(t *testing.T) {
	fileFrom := "testdata/input.txt"
	fileTo := "/tmp/output_hw07.txt"
	offset := int64(0)
	limit := int64(0)

	t.Run("offset is negative", func(t *testing.T) {
		err := Copy(fileFrom, fileTo, -5, limit)
		require.True(t, errors.Is(err, ErrNegativeOffset))
	})

	t.Run("limit is negative", func(t *testing.T) {
		err := Copy(fileFrom, fileTo, offset, -10)
		require.True(t, errors.Is(err, ErrNegativeLimit))
	})

	t.Run("fileFrom is Dir", func(t *testing.T) {
		err := Copy("testdata", fileTo, offset, limit)

		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})

	t.Run("fileFrom does not exist", func(t *testing.T) {
		err := Copy("testdataqqq.txt", fileTo, offset, limit)

		require.Truef(t, os.IsNotExist(err), "actual err - %v", err)
	})

	t.Run("fileTo is Dir", func(t *testing.T) {
		err := Copy(fileFrom, "testdata", offset, limit)

		require.True(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})
}

func TestCopyMainProcess(t *testing.T) {
	smallfile := "testdata/small.txt"

	fileFrom := "testdata/input.txt"
	fileTo := "/tmp/output_hw07.txt"
	offset := int64(0)
	limit := int64(0)

	t.Run("common case", func(t *testing.T) {
		Copy(fileFrom, fileTo, offset, limit)

		require.True(t, filesEqual(fileFrom, fileTo))
	})

	t.Run("same files", func(t *testing.T) {
		fileFrom := fileTo
		err := Copy(fileFrom, fileTo, offset, limit)
		require.True(t, errors.Is(err, ErrSameFile))

		symlink := "/tmp/output_hw07_symlink.txt"
		err = os.Symlink(fileTo, symlink)
		if err != nil {
			panic(err)
		}
		err = Copy(fileFrom, symlink, offset, limit)
		require.True(t, errors.Is(err, ErrSameFile))

		os.Remove(symlink)
	})

	t.Run("rewriting: copy small file after big file to the same target", func(t *testing.T) {
		Copy(fileFrom, fileTo, offset, limit)

		fileFrom := smallfile
		Copy(fileFrom, fileTo, offset, limit)

		require.True(t, filesEqual(fileFrom, fileTo))
	})

	t.Run("limit < buff", func(t *testing.T) {
		limit := int64(5)

		Copy(fileFrom, fileTo, offset, limit)
		require.Equal(t, limit, fileSize(fileTo))
	})

	t.Run("size < offset ", func(t *testing.T) {
		fileFrom := smallfile
		err := Copy(fileFrom, fileTo, int64(2048), limit)

		require.True(t, errors.Is(err, ErrOffsetExceedsFileSize))
	})

	t.Run("buff < limit < size", func(t *testing.T) {
		limit := int64(5)
		oldBuff := BuffLen
		BuffLen = 2

		Copy(fileFrom, fileTo, offset, limit)

		require.Equal(t, limit, fileSize(fileTo))
		BuffLen = oldBuff
	})

	t.Run("size < limit", func(t *testing.T) {
		fileFrom := smallfile
		limit := int64(128)

		Copy(fileFrom, fileTo, offset, limit)

		require.True(t, filesEqual(fileFrom, fileTo))
	})

	t.Run("size = 0", func(t *testing.T) {
		fileFrom := "testdata/empty.txt"

		Copy(fileFrom, fileTo, offset, limit)

		require.True(t, filesEqual(fileFrom, fileTo))
	})
}

func filesEqual(path1, path2 string) bool {
	return fileSize(path1) == fileSize(path2)
}

func fileSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) { // target file may not exist
		panic(err)
	}

	return fileInfo.Size()
}
