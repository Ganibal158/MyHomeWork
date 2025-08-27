package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, err := os.Open(fromPath)
	if err != nil {
		slog.Info("Ошибка при открытии исходного файла", "err", ErrUnsupportedFile)
		return err
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		slog.Info("Ошибка получения информации", "err", err)
		return err
	}
	if !info.Mode().IsRegular() {
		slog.Info("Исходный файл не поддерживается", "err", err)
		return ErrUnsupportedFile
	}
	srcSize := info.Size()
	if offset > srcSize {
		slog.Info("Размер отступа превысил размер файла", "err", ErrOffsetExceedsFileSize)
		return err
	}

	if srcSize-offset < limit || limit == 0 {
		limit = srcSize - offset
	}

	dst, err := os.Create(toPath)
	if err != nil {
		slog.Info("Ошибка создания файла назначения", "err", err)
		return err
	}

	sreader := io.NewSectionReader(src, offset, limit)
	buf := make([]byte, 32*1024) // 32KB буфер

	var copied int64
	var lastPersent int
	printProgress := func(copied, total int64) {
		percent := int(float64(copied) * 100 / float64(total))
		if percent != lastPersent {
			lastPersent = percent
			fmt.Printf("\r[%s%s] %d%%",
				string(repeat('#', (percent/2))),
				string(repeat('-', (50-percent/2))),
				percent,
			)
		}
	}

	for {
		n, err := sreader.Read(buf)
		if n > 0 {
			written, werr := dst.Write(buf[:n])
			if werr != nil {
				return werr
			}
			copied += int64(written)
			printProgress(copied, limit)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	fmt.Println("\nКопирование завершено.")
	return nil
}

func repeat(char rune, count int) []rune {
	res := make([]rune, count)
	for i := range res {
		res[i] = char
	}
	return res
}
