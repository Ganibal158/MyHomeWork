package main

import (
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/cheggaaa/pb/v3"
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
	bar := pb.Full.Start64(limit)

	reader := bar.NewProxyReader(sreader)

	_, err = io.Copy(dst, reader)
	if err != nil {
		slog.Info("Ошибка копирования", "err", err)
		return err
	}

	bar.Finish()

	return nil
}
