package ui

import (
	"fmt"
	"time"

	"github.com/buemura/agnt-cc/internal/scanner"
)

func FormatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func RelativeTime(t time.Time) string {
	d := time.Since(t)

	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}

func PrintTable(files []scanner.FileInfo) {
	if len(files) == 0 {
		fmt.Println("No cache files found matching the criteria.")
		return
	}

	fmt.Printf("%-60s %10s %s\n", "FILE", "SIZE", "MODIFIED")
	fmt.Printf("%-60s %10s %s\n", "----", "----", "--------")

	for _, f := range files {
		name := f.Path
		if len(name) > 58 {
			name = "..." + name[len(name)-55:]
		}
		fmt.Printf("%-60s %10s %s\n", name, FormatSize(f.Size), RelativeTime(f.ModTime))
	}

	total := scanner.TotalSize(files)
	fmt.Printf("\nTotal: %d file(s), %s\n", len(files), FormatSize(total))
}
