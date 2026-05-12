package parser

import (
	"fmt"
	"log/slog"

	"log-parser/internal/application"
)

const (
	fileDBCSV     = "ibdiagnet2.db_csv"
	fileSharpInfo = "ibdiagnet2.sharp_an_info"
)

type Parser struct {
	log      *slog.Logger
	archives *ArchiveReader
}

func New(log *slog.Logger) *Parser {
	if log == nil {
		log = slog.Default()
	}
	return &Parser{
		log:      log,
		archives: NewArchiveReader(log),
	}
}

func (p *Parser) Parse(archivePath string) (application.ParseResult, error) {
	var zero application.ParseResult

	files, err := p.archives.ReadAll(archivePath)
	if err != nil {
		return zero, err
	}

	dbCSVData, ok := files[fileDBCSV]
	if !ok {
		err := fmt.Errorf("archive missing required file: %s", fileDBCSV)
		p.log.Error("parse aborted", "err", err)
		return zero, err
	}

	sharpData, ok := files[fileSharpInfo]
	if !ok {
		err := fmt.Errorf("archive missing required file: %s", fileSharpInfo)
		p.log.Error("parse aborted", "err", err)
		return zero, err
	}

	csv, err := parseDBCSV(dbCSVData)
	if err != nil {
		p.log.Error("parse db csv", "file", fileDBCSV, "err", err)
		return zero, fmt.Errorf("%s: %w", fileDBCSV, err)
	}

	sharpInfos, err := parseSharpInfo(sharpData)
	if err != nil {
		p.log.Error("parse sharp info", "file", fileSharpInfo, "err", err)
		return zero, fmt.Errorf("%s: %w", fileSharpInfo, err)
	}

	p.log.Info("parse completed",
		"path", archivePath,
		"nodes", len(csv.nodes),
		"ports", len(csv.ports),
		"switch_infos", len(csv.switchInfos),
		"system_infos", len(csv.systemInfos),
		"sharp_infos", len(sharpInfos),
	)

	return application.ParseResult{
		Nodes:       csv.nodes,
		Ports:       csv.ports,
		SwitchInfos: csv.switchInfos,
		SystemInfos: csv.systemInfos,
		SharpInfos:  sharpInfos,
	}, nil
}
