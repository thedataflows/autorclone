package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/shirou/gopsutil/v3/process"
	log "github.com/sirupsen/logrus"
)

var Logger *log.Logger

func InitUtils(l *log.Logger) {
	Logger = l
}

// DownloadFile will download a url to a local file
func DownloadFile(filepath string, url string) error {
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(filepath, url)
	// start download
	Logger.Infof("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	if resp.HTTPResponse != nil {
		Logger.Infof("  %v\n", resp.HTTPResponse.Status)
	}
	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
Loop:
	for {
		select {
		case <-t.C:
			Logger.Infof("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())
		case <-resp.Done:
			// download is complete
			break Loop
		}
	}
	// check for errors
	if err := resp.Err(); err != nil {
		return err
	}
	Logger.Infof("Downloaded to '%v'\n", resp.Filename)
	return nil
}

// DecompressFile will extract a list of files from given archive
func DecompressFiles(archivePath, destination string, filesToDecompress []string, stripPath bool) ([]string, error) {
	// TODO maybe check based on mime type instead of simple extension?
	switch path.Ext(archivePath) {
	case ".zip":
		uz := NewUnzip()
		if destination == "" {
			destination, _ = os.Getwd()
		}
		files, err := uz.Extract(archivePath, destination, filesToDecompress, stripPath)
		if err != nil {
			return nil, err
		}
		return files, nil
	default:
		return nil, fmt.Errorf("decompression not yet implemented for '%s'", path.Ext(archivePath))
	}
}

// AppendExtension appends exe if OS is windows
func AppendExtension(fileName string) string {
	var ext = ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	return fileName + ext
}

// ProcessRunning returns the PID of a running process matched by image name and command line
func IsProcessRunning(binaryPath, cmdLine string) (int, error) {
	procs, err := process.Processes()
	if err != nil {
		return 0, err
	}
	command := filepath.Clean(binaryPath)
	if cmdLine != "" {
		command += " " + cmdLine
	}
	for _, p := range procs {
		processCmd, _ := p.Cmdline()
		if strings.Contains(processCmd, command) {
			return int(p.Pid), nil
		}
	}
	return 0, nil
}
