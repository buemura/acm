package main

import (
	"github.com/buemura/acm/cmd/cli"
	_ "github.com/buemura/acm/internal/provider"
)

func main() {
	cli.Execute()
}
