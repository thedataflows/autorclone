package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"dataflows.com/autorclone/internal/pkg/autorclone"
	"dataflows.com/autorclone/internal/pkg/utils"
	cli "github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := new(log.Logger)
	logger.SetFormatter(&log.TextFormatter{ForceColors: true})
	buf := new(bytes.Buffer)
	// TODO improve capturing of stdout/stderr to logger because now is not working. Perhaps overwrite os.Stdout ?
	w := io.MultiWriter(buf, os.Stdout)
	logger.SetOutput(w)
	utils.InitUtils(logger)

	// Home directory and default jobs
	homeDir, _ := os.UserConfigDir()
	autoRcloneHome, _ := filepath.Abs(homeDir + "/autorclone")
	errMkdir := os.MkdirAll(autoRcloneHome, 0700)
	if errMkdir != nil {
		logger.Errorf("Failed to create home directory '%s': %v", autoRcloneHome, errMkdir)
	}
	defaultJobsDefinitionFile, _ := filepath.Abs(autoRcloneHome + "/default_jobs.yaml")

	// Parse CLI
	ctx := cli.Parse(&autorclone.CLI, cli.Bind(logger), cli.Vars{
		"defaultLogLevel":           log.InfoLevel.String(),
		"defaultJobsDefinitionFile": defaultJobsDefinitionFile,
	})
	logger.SetLevel(autorclone.CLI.LogLevel)
	err := ctx.Run(cli.Bind(logger))
	ctx.FatalIfErrorf(err)
}
