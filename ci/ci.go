package ci

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Step struct {
	Command string
	Args    []string
	Stdout  string
	Stderr  string
}

var commands []Step
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

func initCommandsNewSite(name string) []Step {

	commands := append(commands, Step{
		Command: "git",
		Args:    []string{"pull", name},
		Stdout:  "",
		Stderr:  "",
	})

	commands = append(commands, Step{
		Command: "hugo",
		Args:    []string{"-s", name},
		Stdout:  "",
		Stderr:  "",
	})

	return commands
}

func initCommandsExistingSite(name string) []Step {
	commands := append(commands, Step{
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

func Build(name string) error {
	var commands []Step
	if _, err := os.Stat("/tmp/hugosites/" + name); os.IsNotExist(err) {

		commands = initCommandsNewSite(name)
		// commandsForNewSite
		// path/to/whatever does not exist
	} else {
		os.Chdir("/tmp/hugosites/" + name)
		commands = initCommandsExistingSite(name)

		// commandsForExisting site
	}
	wd, _ := os.Getwd()
	for i := range commands {

		err := commands[i].executeCommand()

		fmt.Println("Stdout")
		fmt.Println(commands[i].Stdout)
		fmt.Println("Stderr")
		fmt.Println(commands[i].Stderr)
		if err != nil {
			return err
		}
	}

	return nil
}

func Deploy() {

}
