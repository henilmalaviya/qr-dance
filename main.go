package main

import (
	"os"

	"github.com/henilmalaviya/qr-dance/util/logger"
)

func main() {
	if err := Run(os.Args); err != nil {
		logger.Error("%v", err)
		os.Exit(1)
	}
}
