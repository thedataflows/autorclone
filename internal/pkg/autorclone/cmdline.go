package autorclone

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var Logger *log.Logger
var version = "dev"

// SyncT defines source and destinations to be synchronized
type SyncT struct {
	Source       string        `arg:"" help:"source (rclone pre-configured remote or local directory/file)"`
	Destinations []string      `arg:"" name:"destination1 [destination2] [...]" help:"Space separated rclone pre-configured remotes or local directories"`
	BackupSuffix string        `optional:"" help:"Backs up files with specified .suffix before deleting or replacing them. Existing backups will be overwritten. Set to empty to disable backup" default:"rclonebak"`
	Timeout      time.Duration `help:"Job timeout. Will terminate job after expired time" default:"10m"`
}

// Run executes the 'sync' command after validations
func (s *SyncT) Run(logger *log.Logger) error {
	Logger = logger
	return RunBatchRclone(s)
}

// RunT defines jobs run structure
type RunT struct {
	Jobs               []string      `optional:"" help:"Run just this list of defined job names, either comma separated or repeating this flag"`
	List               bool          `optional:"" help:"List current jobs from jobs definition file"`
	Timeout            time.Duration `help:"Job timeout. Will terminate job after expired time" default:"10m"`
	JobsDefinitionFile string        `help:"Jobs definition yaml file" env:"AUTORCLONE_JOBS_FILE" default:"${defaultJobsDefinitionFile}"`
	jobsDefinition     *JobsDefinitionT
}

// Run executes the 'run' command to execute predefined jobs
func (r *RunT) Run(logger *log.Logger) error {
	Logger = logger
	jd, errRedJobs := ReadJobsDefinition()
	if errRedJobs != nil {
		return errRedJobs
	}
	r.jobsDefinition = jd
	if CLI.Run.List {
		y, _ := yaml.Marshal(r.jobsDefinition.Jobs)
		Logger.Printf("Jobs:\n%s", string(y))
		return nil
	}
	return r.jobsDefinition.RunJobs()
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

// VersionT defines version structure
type VersionT struct {
}

// Run executes the 'version' command to place autorclone in background
func (v *VersionT) Run() error {
	fmt.Printf("%s\n", version)
	return nil
}

// CLI defines the command line arguments and their defaults for the autorclone program
var CLI struct {
	LogLevel log.Level `help:"Set log level to one of: panic, fatal, error, warn, info, debug, trace" default:"${defaultLogLevel}"`

	RclonePath     string `optional:"" help:"Path to rclone binary, by default will try rclone from PATH env" default:"rclone"`
	RcloneVersion  string `optional:"" help:"Rclone release to be downloaded if not in PATH" default:"v1.58.0"`
	RcloneSyncArgs string `optional:"" help:"Rclone default sync arguments" env:"AUTORCLONE_SYNC_ARGS" default:"sync -v --min-size 0.001 --multi-thread-streams 0 --retries 1 --human-readable --track-renames --links --ignore-errors --log-format shortfile"`

	Sync SyncT `cmd:"" help:"Synchronize source to rclone destination(s). Use 'rclone config show' to list them."`

	Run RunT `cmd:"" help:"Manually run predefined sync jobs. Without any argument, will run all jobs in the predefined job definition file"`

	Daemon DaemonT `cmd:"" help:"Run as a background program, executing schelduled jobs"`

	Version VersionT `cmd:"" help:"Show version and exit"`
}
