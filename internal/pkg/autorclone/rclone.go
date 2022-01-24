package autorclone

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// RunBatchRclone attempts to execute rclone in batch mode
func RunBatchRclone(u *UploadT, logger *log.Logger) error {
	for _, dest := range u.Destinations {
		// TODO: improve this, perhaps paralelize?
		success, err := RunIndividualRclone(u.SourceDirectory, dest, logger)
		if !success {
			logger.Errorf("Failed to sync %s to %s. Reason: %s\n", u.SourceDirectory, dest, err)
		}
	}
	return nil
}

// RunIndividualRclone attempts to execute an instance of rclone
func RunIndividualRclone(source string, destination string, logger *log.Logger) (bool, error) {
	rclonePath, err1 := exec.LookPath(CLI.RclonePath)
	if err1 != nil {
		return false, fmt.Errorf("rclone binary '%s' does not exist", CLI.RclonePath)
	}
	args := strings.Split(CLI.RcloneSyncArgs, " ")
	args = append(args, source)
	args = append(args, destination)
	cmd := exec.Command(rclonePath, args...)
	stdOut := new(strings.Builder)
	cmd.Stdout = stdOut
	stdErr := new(strings.Builder)
	cmd.Stderr = stdErr
	logger.Infof("COMMAND: %s\n", cmd)
	err2 := cmd.Run()
	if err2 != nil {
		return false, fmt.Errorf("\n%s\n%s", stdErr.String(), err2)
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
