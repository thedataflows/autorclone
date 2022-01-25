package autorclone

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
	ps "github.com/mitchellh/go-ps"
	log "github.com/sirupsen/logrus"
)

// RunBatchRclone attempts to execute rclone in batch mode
func RunBatchRclone(u *SourceDestT, logger *log.Logger) error {
	failedTasks := 0
	for _, dest := range u.Destinations {
		// TODO: improve this, perhaps paralelize?
		success, err := RunIndividualRclone(u.Source, dest, logger)
		if !success {
			failedTasks++
			logger.Errorf("Failed to sync '%s' to '%s' because %s\n", u.Source, dest, err)
		}
	}
	logger.Infof("Finished. %v tasks successful. %v tasks failed.", len(u.Destinations)-failedTasks, failedTasks)
	return nil
}

// RunIndividualRclone attempts to execute an instance of rclone
func RunIndividualRclone(source string, destination string, logger *log.Logger) (bool, error) {
	// Check if rclone binary exists
	rclonePath, errLookup := exec.LookPath(CLI.RclonePath)
	if errLookup != nil {
		return false, fmt.Errorf("rclone binary '%s' does not exist", CLI.RclonePath)
	}
	// Check if rclone is already running
	pid, errProcRunning := ProcessRunning(rclonePath, logger)
	if errProcRunning != nil {
		return false, errProcRunning
	}
	if pid > 0 {
		return false, fmt.Errorf("'%s' is already running with PID '%v'", rclonePath, pid)
	}
	// Setup arguments
	args := strings.Split(CLI.RcloneSyncArgs, " ")
	if CLI.BackupSuffix != "" {
		args = append(args, "--suffix", "."+CLI.BackupSuffix, "--exclude", "*."+CLI.BackupSuffix)
	}
	args = append(args, source)
	args = append(args, destination)
	// Init command
	rcloneCmd := cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, rclonePath, args...)
	logger.Infof("COMMAND: %s %s", rcloneCmd.Name, strings.Join(rcloneCmd.Args, " "))
	// Print STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		// Done when both channels have been closed
		// https://dave.cheney.net/2013/04/30/curious-channels
		for rcloneCmd.Stdout != nil || rcloneCmd.Stderr != nil {
			select {
			case line, open := <-rcloneCmd.Stdout:
				if !open {
					rcloneCmd.Stdout = nil
					continue
				}
				logger.Println(line)
			case line, open := <-rcloneCmd.Stderr:
				if !open {
					rcloneCmd.Stderr = nil
					continue
				}
				logger.Println(line)
			}
		}
	}()
	// Stop command after specified timeout
	go func() {
		<-time.After(CLI.JobTimeout)
		rcloneCmd.Stop()
		logger.Errorf("Timeout running job after %v", CLI.JobTimeout)
	}()

	// Run and wait for Cmd to return, discard Status
	statusChan := <-rcloneCmd.Start()
	<-doneChan
	if statusChan.Exit > 0 {
		return false, fmt.Errorf("rclone terminated with code %v", statusChan.Exit)
	} else {
		return true, nil
	}
}

// ProcessRunning returns the PID of a running process matched by image name
func ProcessRunning(binaryPath string, logger *log.Logger) (int, error) {
	procs, err := ps.Processes()
	if err != nil {
		return 0, err
	}
	basePath := filepath.Base(binaryPath)
	for _, p := range procs {
		if p.Executable() == basePath {
			// TODO: filter also by process arguments, to be accurate about what instance of rclone is actually running
			return p.Pid(), nil
		}
	}
	return 0, nil
}
