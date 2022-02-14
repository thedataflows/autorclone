package main

import (
	"bytes"
	"io"
	"os"

	"dataflows.com/autorclone/internal/pkg/autorclone"
	"dataflows.com/autorclone/internal/pkg/utils"
	cli "github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := new(log.Logger)
	logger.SetFormatter(&log.TextFormatter{ForceColors: true})
	buf := new(bytes.Buffer)
	// TODO improve capturing of stdout/stderr to logger because now is not working. Perhaps overwrite os.Stderr ?
	w := io.MultiWriter(buf, os.Stdout)
	logger.SetOutput(w)
	utils.InitUtils(logger)

	ctx := cli.Parse(&autorclone.CLI, cli.Bind(logger), cli.Vars{
		"defaultLogLevel": log.InfoLevel.String(),
	})
	logger.SetLevel(autorclone.CLI.LogLevel)
	err := ctx.Run(cli.Bind(logger))
	ctx.FatalIfErrorf(err)
}
