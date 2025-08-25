package main

import (
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	tmpDir := t.TempDir()

	srcFile := tmpDir + "/test_src.txt"
	content := []byte("wefygwefuhwefuhewuh")

	err := os.WriteFile(srcFile, content, 0o644)
	if err != nil {
		t.Fatalf("ошибка создания исходного файла в тесте: %v", err)
	}

	t.Run("Full copy", func(t *testing.T) {
		dstFile := tmpDir + "/test_dst.txt"

		err = Copy(srcFile, dstFile, 0, 0)
		if err != nil {
			t.Errorf("Ошибка при копировании в Full copy: %v", err)
		}

		result, err := os.ReadFile(dstFile)
		if err != nil {
			t.Fatalf("Ошибка чтения файла: %v", err)
		}

		if string(result) != string(content) {
			t.Errorf("Ошибка при копировании: Получили %q, а ожидали %q", result, content)
		}
	})

	t.Run("Offset copy", func(t *testing.T) {
		dstOfFile := tmpDir + "/test_offsetdst.txt"
		offset := int64(5)
		err = Copy(srcFile, dstOfFile, offset, 0)
		if err != nil {
			t.Errorf("Ошибка при копировании в Offset copy: %v", err)
		}
		result, err := os.ReadFile(dstOfFile)
		if err != nil {
			t.Fatalf("Ошибка при чтении: %v", err)
		}
		if string(result) != string(content[offset:]) {
			t.Errorf("Ошибка при копировании: Получили %q, а ожидали %q", result, content[offset:])
		}
	})

	t.Run("Копирование из несуществующего файла", func(t *testing.T) {
		badSrc := tmpDir + "/nonexistent.txt"
		dst := tmpDir + "/bad_dst.txt"

		err := Copy(badSrc, dst, 0, 0)
		if err == nil {
			t.Errorf("Ожидалась ошибка при копировании из несуществующего файла, но её не было")
		}
	})
}
