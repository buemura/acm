package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/buemura/acm/internal/provider"
	"github.com/buemura/acm/internal/scanner"
)

func prompt(reader *bufio.Reader, label string) string {
	fmt.Print(label)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptFilters(reader *bufio.Reader) (providerName, fileType, age string) {
	names := provider.Names()
	providerName = prompt(reader, fmt.Sprintf("Provider [%s] (default: claude): ", strings.Join(names, ", ")))
	if providerName == "" {
		providerName = "claude"
	}

	fileType = prompt(reader, "File type filter (e.g. json, log) or empty for all: ")
	age = prompt(reader, "Minimum file age (e.g. 20m, 24h, 7d) or empty for all: ")
	return
}

func RunMenu() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== Cache Manager ===")
		fmt.Println("1. List cache files")
		fmt.Println("2. Clean cache files")
		fmt.Println("3. Exit")
		fmt.Println()

		choice := prompt(reader, "Choose an option: ")

		switch choice {
		case "1":
			menuList(reader)
		case "2":
			menuClean(reader)
		case "3":
			fmt.Println("Bye!")
			return
		default:
			fmt.Println("Invalid option, try again.")
		}
	}
}

func menuList(reader *bufio.Reader) {
	providerName, fileType, age := promptFilters(reader)

	p, err := provider.Get(providerName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if missing := provider.MissingPathChecks(p); len(missing) > 0 {
		if len(missing) == len(p.Checks) {
			fmt.Printf("No %s cache found. Is %s installed?\n", p.Name, p.Name)
			return
		}
		fmt.Printf(
			"Note: %s provider is missing expected %s in standard locations.\n",
			p.Name,
			strings.Join(missing, ", "),
		)
	}

	minAge, err := scanner.ParseAge(age)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	files, err := scanner.ScanFiles(p.CachePaths, scanner.FilterOpts{
		MinAge:   minAge,
		FileType: fileType,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	PrintTable(files)
}

func menuClean(reader *bufio.Reader) {
	providerName, fileType, age := promptFilters(reader)

	p, err := provider.Get(providerName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if missing := provider.MissingPathChecks(p); len(missing) > 0 {
		if len(missing) == len(p.Checks) {
			fmt.Printf("No %s cache found. Is %s installed?\n", p.Name, p.Name)
			return
		}
		fmt.Printf(
			"Note: %s provider is missing expected %s in standard locations.\n",
			p.Name,
			strings.Join(missing, ", "),
		)
	}

	minAge, err := scanner.ParseAge(age)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	files, err := scanner.ScanFiles(p.CachePaths, scanner.FilterOpts{
		MinAge:   minAge,
		FileType: fileType,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No files found matching the criteria.")
		return
	}

	PrintTable(files)

	total := scanner.TotalSize(files)
	fmt.Printf("\nSpace to be freed: %s\n", FormatSize(total))

	confirm := prompt(reader, "Delete these files? (y/N): ")
	if strings.ToLower(confirm) != "y" {
		fmt.Println("Aborted.")
		return
	}

	deleted := 0
	var freedBytes int64
	for _, f := range files {
		if err := os.Remove(f.Path); err != nil {
			fmt.Printf("Failed to delete %s: %v\n", f.Path, err)
			continue
		}
		deleted++
		freedBytes += f.Size
	}

	fmt.Printf("Deleted %d file(s), freed %s\n", deleted, FormatSize(freedBytes))
}
