package autorclone

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// isDirectory determines if a file represented
// by `path` is a directory or not
func isDirectory(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}
	return nil
}

// UploadT defines upload arguments
type UploadT struct {
	SourceDirectory string   `arg:"" help:"Local directory to be used as a source" type:"path"`
	Destinations    []string `arg:"" name:"destination1 [destination2] [...]" help:"Space separated rclone remotes or local directories"`
}

// Run executes the 'upload' command after validations
func (u *UploadT) Run(logger *log.Logger) error {
	err := isDirectory(u.SourceDirectory)
	if err != nil {
		return err
	}
	return RunBatchRclone(u, logger)
}

// DaemonT defines daemon structure
type DaemonT struct {
	Stop bool `optional:"" help:"Stops running daemon"`
}

// Run executes the 'daemon' command to place autorclone in background
func (u *DaemonT) Run() error {
	return fmt.Errorf("not implemented yet: running as daemon")
}

// CLI defines the command line arguments and their defaults for the autorclone program
var CLI struct {
	LogLevel log.Level `help:"Set log level to one of: panic, fatal, error, warn, info, debug, trace" default:"${defaultLogLevel}"`

	RclonePath     string `optional:"" help:"Path to rclone binary, if empty will use rclone from PATH env" default:"rclone"`
	RcloneSyncArgs string `optional:"" help:"Rclone default sync arguments" env:"AUTO_RCLONE_SYNC_ARGS" default:"-v --min-size 0.001 --multi-thread-streams 0 --retries 1 --human-readable --track-renames --log-format shortfile sync"`

	Upload UploadT `cmd:"" help:"Synchronize source to rclone remote destination(s). Use 'rclone config show' to list them."`

	Daemon DaemonT `cmd:"" help:"Run as a background program, executing schelduled jobs"`
}
