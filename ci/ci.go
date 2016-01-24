package ci

import (
	"bytes"
	"fmt"

	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gophergala2016/hugoku/store"
)

// Step is a task to perform during the deployment
type Step struct {
	Command string
	Args    []string
	Stdout  string
	Stderr  string
}

// ExecuteCommand ...
func (s *Step) ExecuteCommand() error {
	cmd := exec.Command(s.Command, s.Args...)
	stdout, err := cmd.StdoutPipe()
	buf := new(bytes.Buffer)
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	bufErr := new(bytes.Buffer)
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	_, err = buf.ReadFrom(stdout)
	if err != nil {
		return err
	}
	_, err = bufErr.ReadFrom(stderr)
	if err != nil {
		return err
	}
	_, err = bufErr.ReadFrom(stderr)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	s.Stdout = buf.String()
	s.Stderr = bufErr.String()
	return err
}

func initCommandsNewSite(username string, name string, path string) []Step {
	return []Step{
		Step{
			Command: "hugo",
			Args:    []string{"new", "site", path},
			Stdout:  "",
			Stderr:  "",
		},
		Step{
			Command: "git",
			Args:    []string{"clone", "https://github.com/hbpasti/heather-hugo.git", path + "/themes/heather-hugo"},
			Stdout:  "",
			Stderr:  "",
		},
		Step{
			Command: "hugo",
			Args:    []string{"-s", path, "--theme=heather-hugo"},
			Stdout:  "",
			Stderr:  "",
		},
	}
}

func initCommandsExistingSite(username string, name string, path string) []Step {
	return []Step{}
	/*
		commands = append(commands, Step{
			Command: "git",
			Args:    []string{"pull", "origin", "master"},
			Stdout:  "",
			Stderr:  "",
		})

		commands = append(commands, Step{
			Command: "hugo",
			Stdout:  "",
			Stderr:  "",
		})
	*/
}

// Build compiles a project
func Build(username string, name string, path string) (store.BuildInfo, error) {
	var commands []Step
	var buildInfo = store.BuildInfo{BuildPath: path}
	buildInfo.BuildTime = time.Now()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		commands = initCommandsNewSite(username, name, path)
	} else {
		commands = initCommandsExistingSite(username, name, path)
	}
	for i := range commands {
		err := commands[i].ExecuteCommand()
		log.Println("Command:", commands[i].Command, commands[i].Args)
		log.Println("Stdout")
		log.Println(commands[i].Stdout)
		buildInfo.BuildLog += commands[i].Stdout
		log.Println("Stderr")
		log.Println(commands[i].Stderr)
		buildInfo.BuildErrorLog += commands[i].Stderr
		log.Println("-----")
		if err != nil {
			buildInfo.BuildDuration = time.Since(buildInfo.BuildTime)
			buildInfo.BuildStatus = "fail"
			return buildInfo, err
		}
	}
	buildInfo.BuildDuration = time.Since(buildInfo.BuildTime)
	buildInfo.BuildStatus = "ok"
	return buildInfo, nil
}

// Deploy deploys a project
func Deploy(username string, name string) (store.BuildInfo, error) {
	path := fmt.Sprintf("./repos/%s/%s", username, name)
	buildInfo, err := Build(username, name, path)
	return buildInfo, err
}
