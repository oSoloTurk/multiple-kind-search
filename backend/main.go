package main

import (
	_ "github.com/oSoloTurk/multiple-kind-search/docs"
	"github.com/oSoloTurk/multiple-kind-search/internal/cmd"
	"github.com/oSoloTurk/multiple-kind-search/internal/logger"
)

func main() {
	logger.InitLogger()
	cmd.Execute()
}
