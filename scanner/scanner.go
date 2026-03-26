package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FileInfo struct {
	Path    string
	Size    int64
	ModTime time.Time
}

type FilterOpts struct {
	MinAge   time.Duration
	FileType string
}

func ScanFiles(paths []string, opts FilterOpts) ([]FileInfo, error) {
	var files []FileInfo
	seen := map[string]struct{}{}
	now := time.Now()

	for _, pattern := range paths {
		roots, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid path pattern %q: %w", pattern, err)
		}
		if len(roots) == 0 {
			continue
		}

		for _, root := range roots {
			err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil // skip files we can't read
				}
				if info.IsDir() {
					return nil
				}

				if _, ok := seen[path]; ok {
					return nil
				}

				if opts.FileType != "" {
					ext := strings.TrimPrefix(filepath.Ext(path), ".")
					if !strings.EqualFold(ext, opts.FileType) {
						return nil
					}
				}

				if opts.MinAge > 0 {
					age := now.Sub(info.ModTime())
					if age < opts.MinAge {
						return nil
					}
				}

				seen[path] = struct{}{}
				files = append(files, FileInfo{
					Path:    path,
					Size:    info.Size(),
					ModTime: info.ModTime(),
				})
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("scanning %s: %w", root, err)
			}
		}
	}

	return files, nil
}

func TotalSize(files []FileInfo) int64 {
	var total int64
	for _, f := range files {
		total += f.Size
	}
	return total
}

func ParseAge(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}

	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid age format: %s", s)
	}

	unit := s[len(s)-1:]
	valueStr := s[:len(s)-1]

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid age value: %s", s)
	}

	switch unit {
	case "m":
		return time.Duration(value * float64(time.Minute)), nil
	case "h":
		return time.Duration(value * float64(time.Hour)), nil
	case "d":
		return time.Duration(value * float64(24) * float64(time.Hour)), nil
	default:
		return 0, fmt.Errorf("unknown age unit %q (use m, h, or d)", unit)
	}
}
