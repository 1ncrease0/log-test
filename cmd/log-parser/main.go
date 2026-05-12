package main

import (
	"fmt"
	"log/slog"
	"os"

	"log-parser/internal/infrastructure/parser"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	path := "data/log.zip"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	p := parser.New(log)
	res, err := p.Parse(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ok path=%q nodes=%d ports=%d switch_infos=%d system_infos=%d sharp_infos=%d\n",
		path, len(res.Nodes), len(res.Ports), len(res.SwitchInfos), len(res.SystemInfos), len(res.SharpInfos))

	fmt.Println(res.Ports)
}
