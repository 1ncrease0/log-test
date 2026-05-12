package parser

import (
	"archive/zip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type ArchiveReader struct {
	log *slog.Logger
}

func NewArchiveReader(log *slog.Logger) *ArchiveReader {
	return &ArchiveReader{log: log}
}

func (z *ArchiveReader) ReadAll(archivePath string) (map[string][]byte, error) {
	if err := z.validatePath(archivePath); err != nil {
		z.log.Error("archive path invalid", "path", archivePath, "err", err)
		return nil, err
	}

	r, err := zip.OpenReader(archivePath)
	if err != nil {
		z.log.Error("open zip", "path", archivePath, "err", err)
		return nil, fmt.Errorf("open zip %s: %w", archivePath, err)
	}
	defer func() {
		if cerr := r.Close(); cerr != nil {
			z.log.Warn("close zip reader", "path", archivePath, "err", cerr)
		}
	}()

	result := make(map[string][]byte, len(r.File))
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		data, err := z.readZipEntry(f)
		if err != nil {
			z.log.Error("read zip entry", "entry", f.Name, "err", err)
			return nil, err
		}
		result[filepath.Base(f.Name)] = data
	}
	return result, nil
}

func (z *ArchiveReader) readZipEntry(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("open zip entry %s: %w", f.Name, err)
	}
	defer func() {
		if cerr := rc.Close(); cerr != nil {
			z.log.Warn("close zip entry reader", "entry", f.Name, "err", cerr)
		}
	}()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("read zip entry %s: %w", f.Name, err)
	}
	return data, nil
}

func (z *ArchiveReader) validatePath(archivePath string) error {
	if !strings.EqualFold(filepath.Ext(archivePath), ".zip") {
		return fmt.Errorf("only .zip archives are supported: %s", archivePath)
	}

	absArchive, err := filepath.Abs(archivePath)
	if err != nil {
		return fmt.Errorf("resolve archive path: %w", err)
	}
	absData, err := filepath.Abs("data")
	if err != nil {
		return fmt.Errorf("resolve data directory: %w", err)
	}

	rel, err := filepath.Rel(absData, absArchive)
	if err != nil {
		return fmt.Errorf("archive must be inside data directory: %w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("archive path must be inside data directory: %s", archivePath)
	}

	info, err := os.Stat(absArchive)
	if err != nil {
		return fmt.Errorf("stat archive: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("path is a directory, expected .zip file: %s", archivePath)
	}

	return nil
}
