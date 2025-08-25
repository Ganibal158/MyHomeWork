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

	dst, err := os.Create(toPath)
	if err != nil {
		slog.Info("Ошибка создания файла назначения", "err", err)
		return err
	}
	info, err := src.Stat()
	if err != nil {
		slog.Info("Ошибка получения информации", "err", err)
		return err
	}
	srcSize := info.Size()
	if offset > srcSize {
		slog.Info("Размер отступа превысил размер файла", "err", ErrOffsetExceedsFileSize)
		return err
	}

	if srcSize-offset < limit || limit == 0 {
		limit = srcSize - offset
	}

	sreader := io.NewSectionReader(src, offset, limit)
	buf := make([]byte, 32*1024) // 32KB буфер

	var copied int64
	const progressBarWidth = 50

	printProgress := func(copied, total int64) {
		percent := float64(copied) / float64(total)
		filled := int(percent * progressBarWidth)
		empty := progressBarWidth - filled
		fmt.Printf("\r[%s%s] %.2f%%",
			string(repeat('#', filled)),
			string(repeat('-', empty)),
			percent*100,
		)
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
