package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Использование: go-envdir <директория> <команда>")
		os.Exit(1)
	}

	dir := os.Args[1]
	cmdArgs := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println("Ошибка при чтении директории:", err)
		os.Exit(1)
	}

	exitCode := RunCmd(cmdArgs, env)

	os.Exit(exitCode)
}
