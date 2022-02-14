package autorclone

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"dataflows.com/autorclone/internal/pkg/utils"
	"github.com/go-cmd/cmd"
)

// RunBatchRclone attempts to execute rclone in batch mode
func RunBatchRclone(u *SourceDestT) error {
	failedTasks := 0
	for _, dest := range u.Destinations {
		// Check if rclone binary exists
		rclonePath, errLookup := exec.LookPath(CLI.RclonePath)
		if errLookup != nil {
			Logger.Errorf("Rclone binary '%s' does not exist, trying to download release %s", CLI.RclonePath, CLI.RcloneVersion)
			errEnsureRclone := EnsureRclone()
			if errEnsureRclone != nil {
				return errEnsureRclone
			}
			// on Windows will look for any of {".com", ".exe", ".bat", ".cmd"}
			rclonePath, errLookup = exec.LookPath("./rclone")
			if errLookup != nil {
				return errLookup
			}
		}
		// TODO: improve this, perhaps paralelize?
		success, err := RunIndividualRclone(rclonePath, u.Source, dest)
		if !success {
			failedTasks++
			Logger.Errorf("Failed to sync '%s' to '%s' because %s\n", u.Source, dest, err)
		}
	}
	Logger.Infof("Finished. %v tasks successful. %v tasks failed.", len(u.Destinations)-failedTasks, failedTasks)
	return nil
}

// RunIndividualRclone attempts to execute an instance of rclone
func RunIndividualRclone(rclonePath string, source string, destination string) (bool, error) {
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
	commandLine := strings.Join(rcloneCmd.Args, " ")
	Logger.Infof("COMMAND: %s %s", rcloneCmd.Name, commandLine)
	// Check if rclone is already running
	pid, errIsProcessRunning := utils.IsProcessRunning(rcloneCmd.Name, commandLine)
	if errIsProcessRunning != nil {
		return false, errIsProcessRunning
	}
	if pid > 0 {
		return false, fmt.Errorf("'%s %s' is already running with PID '%v'", rcloneCmd.Name, commandLine, pid)
	}
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
				Logger.Println(line)
			case line, open := <-rcloneCmd.Stderr:
				if !open {
					rcloneCmd.Stderr = nil
					continue
				}
				Logger.Println(line)
			}
		}
	}()
	// Stop command after specified timeout
	go func() {
		<-time.After(CLI.JobTimeout)
		rcloneCmd.Stop()
		Logger.Errorf("Timeout running job after %v", CLI.JobTimeout)
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

// EnsureRclone will download and extract specified or default version of rclone to current working directory
func EnsureRclone() error {
	fileUrl := fmt.Sprintf("https://github.com/rclone/rclone/releases/download/%s/rclone-%s-%s-%s.zip", CLI.RcloneVersion, CLI.RcloneVersion, runtime.GOOS, runtime.GOARCH)
	errDownload := utils.DownloadFile(".", fileUrl)
	if errDownload != nil {
		return errDownload
	}
	extractedFiles, errDecompress := utils.DecompressFiles(path.Base(fileUrl), "", []string{
		fmt.Sprintf("rclone-%s-%s-%s/%s", CLI.RcloneVersion, runtime.GOOS, runtime.GOARCH, utils.AppendExtension("rclone")),
	}, true)
	if errDecompress != nil {
		return errDecompress
	}
	for _, f := range extractedFiles {
		Logger.Infof("Extracted %s", f)
		if path.Base(f) == utils.AppendExtension("rclone") {
			os.Chmod(path.Base(f), 0755)
		}
	}
	return nil
}
