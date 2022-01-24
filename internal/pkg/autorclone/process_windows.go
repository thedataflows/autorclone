//go:build windows

package autorclone

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func GetProcessArguments(pid int, logger *log.Logger) (string, error) {
	// TODO: this is a stub
	return "", fmt.Errorf("not yet implemented")
}
