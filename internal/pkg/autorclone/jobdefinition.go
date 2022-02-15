package autorclone

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var JobsDefinitionVersions = "v1"

type JobsDefinitionT struct {
	Version string                `yaml:"version"`
	Flags   JobsDefinitionFlagsT  `yaml:"flags"`
	Jobs    []JobsDefinitionListT `yaml:"jobs"`
}

type JobsDefinitionFlagsT struct {
	BackupSuffix string        `yaml:"backupsuffix"`
	Timeout      time.Duration `yaml:"timeout,omitempty"`
}

type JobsDefinitionListT struct {
	Name         string        `yaml:"name"`
	Source       string        `yaml:"source"`
	Destinations []string      `yaml:"destinations"`
	Schedules    []string      `yaml:"schedules,omitempty"`
	Log          string        `yaml:"log,omitempty"`
	Timeout      time.Duration `yaml:"timeout,omitempty"`
}

// ReadJobsDefinition loads job definitions from yaml
func ReadJobsDefinition() (*JobsDefinitionT, error) {
	yamlFile, errRead := ioutil.ReadFile(CLI.Run.JobsDefinitionFile)
	if errRead != nil {
		return nil, errRead
	}
	data := &JobsDefinitionT{}
	errUnmarshal := yaml.Unmarshal(yamlFile, &data)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	if data.Version != JobsDefinitionVersions {
		return nil, fmt.Errorf("jobs definition version must be %s", JobsDefinitionVersions)
	}
	if CLI.LogLevel == logrus.DebugLevel {
		y, _ := yaml.Marshal(data)
		Logger.Debugf("%s\n%s", CLI.Run.JobsDefinitionFile, string(y))
	}
	return data, nil
}

func (j *JobsDefinitionT) RunOneJob(jobNumber int, job JobsDefinitionListT) error {
	Logger.Infof("Running job #%v '%s'", jobNumber, job.Name)
	timeout := j.Flags.Timeout
	if timeout == 0 {
		timeout = CLI.Run.Timeout
	}
	sync := &SyncT{job.Source, job.Destinations, j.Flags.BackupSuffix, timeout}
	err := RunBatchRclone(sync)
	if err != nil {
		Logger.Errorf("%s", err)
	}
	return nil
}

func (j *JobsDefinitionT) RunJobs() error {
	for ij, job := range j.Jobs {
		if len(CLI.Run.Jobs) == 0 {
			j.RunOneJob(ij, job)
			continue
		}
		for _, filterJob := range CLI.Run.Jobs {
			if job.Name == filterJob {
				j.RunOneJob(ij, job)
			}
		}
	}
	return nil
}
