package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var EmptyEnv = make(Environment)

func TestRunCmd(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		code := RunCmd([]string{"echo", "call", "me"}, EmptyEnv)
		require.Equal(t, 0, code)
	})

	t.Run("without args", func(t *testing.T) {
		code := RunCmd([]string{"pwd"}, EmptyEnv)
		require.Equal(t, 0, code)
	})

	t.Run("exit code 2", func(t *testing.T) {
		code := RunCmd([]string{"ls", "qazaqazaq"}, EmptyEnv) // no such file or directory
		require.Equal(t, 2, code)
	})
}

func TestSetEnv(t *testing.T) {
	t.Run("several sets", func(t *testing.T) {
		expected := make(map[string]string)
		expected["homework"] = "07" // unexisted
		expected["OLDPWD"] = "/usr" // existed

		env := make(Environment)
		for name := range expected {
			env[name] = &EnvValue{Value: expected[name], NeedRemove: false}
		}

		SetEnv(env)

		for name := range expected {
			require.Equal(t, expected[name], os.Getenv(name))
		}
	})

	t.Run("remove", func(t *testing.T) {
		env := make(Environment)
		name := "kapuki"
		env[name] = &EnvValue{Value: "kanuki", NeedRemove: false}

		SetEnv(env)
		require.Equal(t, env[name].Value, os.Getenv(name)) // written

		env[name].NeedRemove = true
		SetEnv(env)
		require.Equal(t, "", os.Getenv(name)) // removed
	})
}
