package autorclone

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	ps "github.com/mitchellh/go-ps"
	log "github.com/sirupsen/logrus"
)

// RunBatchRclone attempts to execute rclone in batch mode
func RunBatchRclone(u *UploadT, logger *log.Logger) error {
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
	args = append(args, source)
	args = append(args, destination)
	cmd := exec.Command(rclonePath, args...)
	stdOut := new(strings.Builder)
	cmd.Stdout = stdOut
	stdErr := new(strings.Builder)
	cmd.Stderr = stdErr
	logger.Infof("COMMAND: %s\n", cmd)
	// Run rclone
	errRun := cmd.Run()
	if errRun != nil {
		return false, fmt.Errorf("\n%s\n%s", stdErr.String(), errRun)
	}
	logger.Infof("PID: %v\n", cmd.ProcessState.Pid())
	if stdOut.Len() > 0 {
		logger.Info("STDOUT:\n%s", stdOut.String())
	}
	if stdErr.Len() > 0 {
		logger.Infof("STDERR:\n%s", stdErr.String())
	}
	return true, nil
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
