package main

import (
	"github.com/buemura/agnt-cc/internal/cli"
	_ "github.com/buemura/agnt-cc/internal/provider"
)

func main() {
	cli.Execute()
}
