package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(command []string, env Environment) (returnCode int) {
	SetEnv(env)

	cmd := exec.Command(command[0], command[1:]...) //nolint:gosec

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	return cmd.ProcessState.ExitCode()
}

func SetEnv(env Environment) {
	for name := range env {
		if env[name].NeedRemove {
			os.Unsetenv(name)
			continue
		}
		os.Setenv(name, env[name].Value)
	}
}
