package main

import (
	"flag"
	"log/slog"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	if from == "" || to == "" {
		slog.Info("Пути к исходному файлу и файлу в который выполняется запись не должны быть пустами")
		return
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		slog.Info("Возникла ошибка при копировании", "err", err)
	}
}
