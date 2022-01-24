package autorclone

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// UploadT defines upload arguments
type UploadT struct {
	Source       string   `arg:"" help:"source (rclone config or local directory/file)"`
	Destinations []string `arg:"" name:"destination1 [destination2] [...]" help:"Space separated rclone remotes or local directories"`
}

// Run executes the 'upload' command after validations
func (u *UploadT) Run(logger *log.Logger) error {
	return RunBatchRclone(u, logger)
}

// DaemonT defines daemon structure
type DaemonT struct {
	Stop bool `optional:"" help:"Stops running daemon"`
}

// Run executes the 'daemon' command to place autorclone in background
func (u *DaemonT) Run() error {
	// TODO: implement running as daemon
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
