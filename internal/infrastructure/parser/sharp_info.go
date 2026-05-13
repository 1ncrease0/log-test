package parser

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"log-parser/internal/domain"
)

const sharpKeyValueParts = 2

func parseSharpInfo(data []byte) ([]domain.NodeSharpInfo, error) {
	var result []domain.NodeSharpInfo
	var current *domain.NodeSharpInfo

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "---") {
			continue
		}

		if strings.HasPrefix(line, "SW_GUID=") {
			if current != nil {
				result = append(result, *current)
			}
			guid := strings.TrimPrefix(line, "SW_GUID=")
			current = &domain.NodeSharpInfo{NodeGUID: normSharpNodeGUID(guid)}
			continue
		}

		if current == nil {
			continue
		}

		parts := strings.SplitN(line, "=", sharpKeyValueParts)
		if len(parts) != sharpKeyValueParts {
			return nil, fmt.Errorf("line %d: invalid key=value: %q", lineNum, line)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		n, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("line %d: key %q: invalid int %q", lineNum, key, val)
		}

		switch key {
		case "endianness":
			current.Endianness = n
		case "enable_endianness_per_job":
			current.EnableEndiannessPerJob = n
		case "reproducibility_disable":
			current.ReproducibilityDisable = n
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan sharp_an_info: %w", err)
	}
	if current != nil {
		result = append(result, *current)
	}
	return result, nil
}
