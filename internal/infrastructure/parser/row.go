package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type rowParser struct {
	row []string
	idx map[string]int
	err error
}

func newRowParser(row []string, idx map[string]int) *rowParser {
	return &rowParser{row: row, idx: idx}
}

func (r *rowParser) readString(dst *string, col string) {
	if r.err != nil {
		return
	}
	*dst, r.err = getCol(r.row, r.idx, col)
}

func (r *rowParser) readInt(dst *int, col string) {
	if r.err != nil {
		return
	}
	*dst, r.err = intCol(r.row, r.idx, col)
}

func (r *rowParser) readInt64(dst *int64, col string) {
	if r.err != nil {
		return
	}
	*dst, r.err = int64Col(r.row, r.idx, col)
}

func (r *rowParser) readIntOpt(dst **int, col string) {
	if r.err != nil {
		return
	}
	*dst, r.err = optIntCol(r.row, r.idx, col)
}

func intCol(row []string, idx map[string]int, col string) (int, error) {
	s, err := getCol(row, idx, col)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("column %q: %w", col, err)
	}
	return n, nil
}

func int64Col(row []string, idx map[string]int, col string) (int64, error) {
	s, err := getCol(row, idx, col)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		u, err2 := strconv.ParseUint(s, 10, 64)
		if err2 != nil {
			return 0, fmt.Errorf("column %q: %w", col, err)
		}
		return int64(u), nil
	}
	return n, nil
}

func optIntCol(row []string, idx map[string]int, col string) (*int, error) {
	s, err := getCol(row, idx, col)
	if err != nil {
		return nil, err
	}
	if s == "N/A" {
		return nil, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("column %q: %w", col, err)
	}
	return &n, nil
}

func getCol(row []string, idx map[string]int, col string) (string, error) {
	i, ok := idx[col]
	if !ok {
		return "", fmt.Errorf("column %q not found", col)
	}
	if i >= len(row) {
		return "", fmt.Errorf("column %q: index %d out of range (row len %d)", col, i, len(row))
	}
	return strings.TrimSpace(row[i]), nil
}
