package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	TestDir     = "testdata/emptydir"
	EmptyResult = make(Environment)
)

func TestReadDir(t *testing.T) {
	t.Run("empty dirt string", func(t *testing.T) {
		result, err := ReadDir("")

		require.NoError(t, err)
		require.Exactly(t, EmptyResult, result)
	})

	t.Run("dir path does not exist", func(t *testing.T) {
		result, err := ReadDir("qqq")

		require.True(t, os.IsNotExist(err))
		require.Nil(t, result)
	})

	t.Run("empty filelist", func(t *testing.T) {
		makeTestDir()
		defer os.RemoveAll(TestDir)

		result, err := ReadDir(TestDir)

		require.NoError(t, err)
		require.Exactly(t, EmptyResult, result)
	})

	t.Run("dirs and a dot file", func(t *testing.T) {
		result, err := ReadDir("testdata")

		require.NoError(t, err)
		require.Exactly(t, EmptyResult, result)
	})

	t.Run("correct keys", func(t *testing.T) {
		makeTestDir()
		defer os.RemoveAll(TestDir)

		expected := make(Environment)
		tests := []struct {
			filename string
			filedata string
			envname  string
			envdata  string
		}{
			{filename: "User", filedata: "OTUS", envname: "User", envdata: "OTUS"},
			{filename: "Bkk", filedata: "space  \t\t", envname: "Bkk", envdata: "space"},
			{filename: "double", filedata: "name\nsurname", envname: "double", envdata: "name"},
			{filename: "empty", filedata: "\nanything", envname: "empty", envdata: ""},
		}
		for _, ts := range tests {
			saveTestFile(ts.filename, ts.filedata)
			expected[ts.envname] = &EnvValue{Value: ts.envdata, NeedRemove: false}
		}

		result, err := ReadDir(TestDir)

		require.NoError(t, err)
		require.Exactly(t, expected, result)
	})

	t.Run("incorrect keys", func(t *testing.T) {
		makeTestDir()
		defer os.RemoveAll(TestDir)

		tests := []struct {
			filename string
			filedata string
		}{
			{filename: "eq=uality", filedata: "hmm"},
			{filename: "dot.txt", filedata: "anything"},
			{filename: "ping-pong", filedata: "ball"},
		}
		for _, ts := range tests {
			saveTestFile(ts.filename, ts.filedata)
		}

		result, err := ReadDir(TestDir)

		require.NoError(t, err)
		require.Exactly(t, EmptyResult, result)
	})
}

func saveTestFile(name, data string) {
	fullname := filepath.Join(TestDir, name)
	file, err := os.Create(fullname)
	if err != nil {
		panic(err)
	}

	file.Write([]byte(data))
}

func makeTestDir() {
	err := os.Mkdir(TestDir, 0o755)
	if err != nil {
		panic(err)
	}
}
