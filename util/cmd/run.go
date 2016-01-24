package cmd

import (
	"log"
	"os/exec"
	"strings"
)

// Run runs a command and handles possible error
func Run(command string, args []string) {
	log.Println("Running command " + command)
	log.Println(strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
