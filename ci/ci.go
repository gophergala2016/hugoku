package ci

import (
	"bytes"
	"fmt"

	"log"
	"os"
	"os/exec"
)

// Step is a task to perform during the deployment
type Step struct {
	Command string
	Args    []string
	Stdout  string
	Stderr  string
}

func (s *Step) executeCommand() error {
	cmd := exec.Command(s.Command, s.Args...)
	stdout, err := cmd.StdoutPipe()
	buf := new(bytes.Buffer)
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	bufErr := new(bytes.Buffer)
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	_, err = buf.ReadFrom(stdout)
	if err != nil {
		log.Fatal(err)
	}
	_, err = bufErr.ReadFrom(stderr)
	if err != nil {
		log.Fatal(err)
	}
	_, err = bufErr.ReadFrom(stderr)
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	s.Stdout = buf.String()
	s.Stderr = bufErr.String()

	return err
}

func initCommandsNewSite(username string, name string, path string) []Step {
	var commands []Step

	commands = append(commands, Step{
		Command: "hugo",
		Args:    []string{"new", "site", path},
		Stdout:  "",
		Stderr:  "",
	})

	commands = append(commands, Step{
		Command: "git",
		Args:    []string{"clone", "https://github.com/hbpasti/heather-hugo.git", path + "/themes/heather-hugo"},
		Stdout:  "",
		Stderr:  "",
	})

	commands = append(commands, Step{
		Command: "hugo",
		Args:    []string{"-s", path, "--theme=heather-hugo"},
		Stdout:  "",
		Stderr:  "",
	})

	return commands
}

func initCommandsExistingSite(username string, name string, path string) []Step {
	var commands []Step
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
	return commands
}

// Build compiles a project
func Build(username string, name string, path string) error {
	var commands []Step
	if _, err := os.Stat(path); os.IsNotExist(err) {
		commands = initCommandsNewSite(username, name, path)
	} else {
		commands = initCommandsExistingSite(username, name, path)
	}

	for i := range commands {
		err := commands[i].executeCommand()
		log.Println("Command:", commands[i].Command, commands[i].Args)
		log.Println("Stdout:")
		log.Println(commands[i].Stdout)
		log.Println("Stderr:")
		log.Println(commands[i].Stderr)
		log.Println("-----")
		if err != nil {
			return err
		}
	}

	return nil
}

// Deploy deploys a project
func Deploy(username string, name string) (path string, err error) {
	path = fmt.Sprintf("./repos/%s/%s", username, name)
	err = Build(username, name, path)
	return path, err
}
