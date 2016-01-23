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

var listCommands []string

func (s *Step) executeCommand() error {
	fmt.Println(s.Command)
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

	buf.ReadFrom(stdout)
	bufErr.ReadFrom(stderr)

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
		Args:    []string{"clone", "git@github.com:hbpasti/heather-hugo.git", path + "/themes/heather-hugo"},
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

	return commands
}

// Build compiles a project
func Build(username string, name string, path string) error {
	var commands []Step
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Not exists
		commands = initCommandsNewSite(username, name, path)
		// commandsForNewSite
		// path/to/whatever does not exist

	} else {
		fmt.Println("repo exists")
		// os.Chdir(path)
		// commands = initCommandsExistingSite(name)
		// commandsForExisting site
	}

	for i := range commands {
		err := commands[i].executeCommand()
		fmt.Println("Command")
		fmt.Println(commands[i].Command, commands[i].Args)
		fmt.Println("Stdout")
		fmt.Println(commands[i].Stdout)
		fmt.Println("Stderr")
		fmt.Println(commands[i].Stderr)
		fmt.Println("-----")
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
