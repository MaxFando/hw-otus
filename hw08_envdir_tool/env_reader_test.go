package main

import (
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {
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

	env, err := ReadDir(tmpDir)
	if err != nil {
		t.Fatal("Ошибка при чтении директории:", err)
	}

	if len(env) != 2 {
		t.Errorf("Ожидается 2 переменных окружения, получено %d", len(env))
	}

	if env["FOO"].Value != "123" || env["FOO"].NeedRemove {
		t.Errorf("Неверное значение для переменной окружения FOO")
	}

	if env["BAR"].Value != "value" || env["BAR"].NeedRemove {
		t.Errorf("Неверное значение для переменной окружения BAR")
	}
}
