package autorclone

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// SourceDestT defines source and destinations
type SourceDestT struct {
	Source       string   `arg:"" help:"source (rclone pre-configured remote or local directory/file)"`
	Destinations []string `arg:"" name:"destination1 [destination2] [...]" help:"Space separated rclone pre-configured remotes or local directories"`
}

// Run executes the 'upload' command after validations
func (u *SourceDestT) Run(logger *log.Logger) error {
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

	RclonePath     string        `optional:"" help:"Path to rclone binary, by default will try rclone from PATH env" default:"rclone"`
	RcloneSyncArgs string        `optional:"" help:"Rclone default sync arguments" env:"AUTO_RCLONE_SYNC_ARGS" default:"sync -v --min-size 0.001 --multi-thread-streams 0 --retries 1 --human-readable --track-renames --links --ignore-errors --log-format shortfile"`
	BackupSuffix   string        `help:"Backs up files with specified .suffix before deleting or replacing them. Existing backups will be overwritten. Set to empty to disable backup" default:"rclonebak"`
	JobTimeout     time.Duration `help:"Job timeout. Will terminate rclone after expired time" default:"10m"`

	Sync SourceDestT `cmd:"" help:"Synchronize source to rclone destination(s). Use 'rclone config show' to list them."`

	Daemon DaemonT `cmd:"" help:"Run as a background program, executing schelduled jobs"`
}
