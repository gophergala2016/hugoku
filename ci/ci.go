package ci

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

type Step struct {
	Command string
	Stdout  string
	Stderr  string
}

var commands []Step
var listCommands []string

func (s *Step) executeCommand() error {
	fmt.Println(s.Command)

	cmd := exec.Command(s.Command)
	stdout, err := cmd.StdoutPipe()
	buf := new(bytes.Buffer)
	if err != nil {
		fmt.Println("Ah!, se siente")
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Ah!, se siente")
		log.Fatal(err)
	}
	buf.ReadFrom(stdout)
	if err := cmd.Wait(); err != nil {
		fmt.Println("Ah!, se siente")
		log.Fatal(err)
	}

	s.Stdout = buf.String()

	return err
}

func Build(name string) error {

	commands := append(commands, Step{
		Command: "ls",
		Stdout:  "",
		Stderr:  "",
	})

	for i := range commands {

		err := commands[i].executeCommand()
		fmt.Println(commands[i].Stdout)
		return err
	}

	return nil
}

func Deploy() {

}
