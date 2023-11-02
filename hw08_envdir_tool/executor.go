package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type RunParams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type RunOption func(*RunParams)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, opts ...RunOption) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Env = mergeWithEnviron(env)
	runParams := RunParams{
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
		Stdout: os.Stdout,
	}
	for _, opt := range opts {
		opt(&runParams)
	}
	command.Stdin = runParams.Stdin
	command.Stderr = runParams.Stderr
	command.Stdout = runParams.Stdout

	if err := command.Run(); err != nil {
		log.Fatal(err)
		return -1
	}
	return command.ProcessState.ExitCode()
}

func mergeWithEnviron(envs Environment) []string {
	envMap := getEnviron()
	for key, env := range envs {
		if env.NeedRemove {
			delete(envMap, key)
			continue
		}
		envMap[key] = fmt.Sprintf("%v=%v", key, env.Value)
	}

	result := make([]string, 0, len(envs))
	for _, value := range envMap {
		result = append(result, value)
	}
	return result
}

func getEnviron() map[string]string {
	envs := os.Environ()
	result := make(map[string]string, len(envs))
	for _, env := range os.Environ() {
		parts := strings.Split(env, "=")
		if len(parts) < 2 {
			continue
		}
		result[parts[0]] = env
	}
	return result
}
