package main

import (
	"github.com/oSoloTurk/multiple-kind-search/cmd"
	_ "github.com/oSoloTurk/multiple-kind-search/docs"
	"github.com/oSoloTurk/multiple-kind-search/internal/logger"
)

func main() {
	logger.InitLogger()
	cmd.Execute()
}
