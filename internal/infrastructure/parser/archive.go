package parser

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"log-parser/internal/application"
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
		entryData, readErr := z.readZipEntry(f)
		if readErr != nil {
			z.log.Error("read zip entry", "entry", f.Name, "err", readErr)
			return nil, fmt.Errorf("read entry %s: %w", f.Name, readErr)
		}
		result[filepath.Base(f.Name)] = entryData
	}
	return result, nil
}

func (z *ArchiveReader) ResolveRelative(rel string) (string, error) {
	rel = filepath.Clean(strings.TrimSpace(rel))
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: path escapes data directory", application.ErrInvalidPath)
	}
	full := filepath.Join("data", rel)
	absArchive, absErr := filepath.Abs(full)
	if absErr != nil {
		return "", fmt.Errorf("resolve archive path: %w", absErr)
	}
	if valErr := z.validatePath(absArchive); valErr != nil {
		return "", valErr
	}
	return absArchive, nil
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
		return fmt.Errorf("%w: only .zip archives are supported: %s", application.ErrInvalidPath, archivePath)
	}

	absArchive, err := filepath.Abs(archivePath)
	if err != nil {
		return fmt.Errorf("%w: resolve archive path: %w", application.ErrInvalidPath, err)
	}
	absData, err := filepath.Abs("data")
	if err != nil {
		return fmt.Errorf("%w: resolve data directory: %w", application.ErrInvalidPath, err)
	}

	rel, err := filepath.Rel(absData, absArchive)
	if err != nil {
		return fmt.Errorf("%w: archive must be inside data directory: %w", application.ErrInvalidPath, err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("%w: archive path must be inside data directory: %s", application.ErrInvalidPath, archivePath)
	}

	info, err := os.Stat(absArchive)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w: stat archive: %w", application.ErrArchiveNotFound, err)
		}
		return fmt.Errorf("%w: stat archive: %w", application.ErrInvalidPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%w: path is a directory, expected .zip file: %s", application.ErrInvalidPath, archivePath)
	}

	return nil
}
