package main

import (
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "envdir_test")
	if err != nil {
		t.Fatal("Ошибка при создании временной директории:", err)
	}
	defer os.RemoveAll(tmpDir)

	file1 := tmpDir + "/FOO"
	err = os.WriteFile(file1, []byte("123"), 0o644)
	if err != nil {
		t.Fatal("Ошибка при создании файла:", err)
	}

	file2 := tmpDir + "/BAR"
	err = os.WriteFile(file2, []byte("value\n"), 0o644)
	if err != nil {
		t.Fatal("Ошибка при создании файла:", err)
	}

	cmd := []string{"echo", "$FOO", "$BAR"}
	env := Environment{
		"FOO": EnvValue{Value: "123", NeedRemove: false},
		"BAR": EnvValue{Value: "value", NeedRemove: false},
	}

	returnCode := RunCmd(cmd, env)
	if returnCode != 0 {
		t.Errorf("Ожидается код возврата 0, получено %d", returnCode)
	}
}
