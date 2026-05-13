package parser

import (
	"archive/zip"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeZip(t *testing.T, path string, files map[string][]byte) {
	t.Helper()

	f, err := os.Create(path)
	require.NoError(t, err)

	zw := zip.NewWriter(f)
	for name, data := range files {
		w, werr := zw.Create(name)
		require.NoError(t, werr)
		_, err = w.Write(data)
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	require.NoError(t, f.Close())
}

func TestArchiveReader_ResolveRelative(t *testing.T) {
	tmp := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, "data"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(tmp, "data", "a.zip"), []byte("x"), 0o644))
	t.Chdir(tmp)

	z := NewArchiveReader(slog.New(slog.NewTextHandler(io.Discard, nil)))
	got, err := z.ResolveRelative("a.zip")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(filepath.ToSlash(got), "data/a.zip"))
}

func TestArchiveReader_ResolveRelative_RejectsEscape(t *testing.T) {
	t.Parallel()

	z := NewArchiveReader(slog.New(slog.NewTextHandler(io.Discard, nil)))
	_, err := z.ResolveRelative("../secret.zip")
	require.Error(t, err)
}

func TestArchiveReader_ReadAll(t *testing.T) {
	tmp := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, "data"), 0o755))
	zipPath := filepath.Join(tmp, "data", "b.zip")
	writeZip(t, zipPath, map[string][]byte{"nested/x.txt": []byte("hello")})
	t.Chdir(tmp)

	z := NewArchiveReader(slog.New(slog.NewTextHandler(io.Discard, nil)))
	m, err := z.ReadAll(zipPath)
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), m["x.txt"])
}
