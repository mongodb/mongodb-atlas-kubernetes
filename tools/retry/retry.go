package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func retry(times int, pause time.Duration, command string, args ...string) (int, error) {
	var err error
	var cmd *exec.Cmd
	for i := 0; i <= times; i++ {
		cmd = exec.Command(command, args...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err = cmd.Run(); err == nil {
			break
		}
		time.Sleep(pause)
	}
	return cmd.ProcessState.ExitCode(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage:\n$ [RETRIES=x] [PAUSE=y] %s {command to retry}", os.Args[0])
		os.Exit(1)
	}
	retries := 7
	retriesArg := os.Getenv("RETRIES")
	if retriesArg != "" {
		var err error
		retries, err = strconv.Atoi(retriesArg)
		if err != nil {
			log.Fatalf("Failed to convert RETRIES from %s: %v", err)
		}
		if retries < 1 {
			log.Fatalf("RETRIES must be 1 or more but got %v", retries)
		}
	}
	pause := time.Second
	pauseArg := os.Getenv("PAUSE")
	if pauseArg != "" {
		var err error
		pause, err = time.ParseDuration(pauseArg)
		if err != nil {
			log.Fatalf("Failed to convert PAUSE from %s: %v", err)
		}
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	exitCode, err := retry(retries, pause, cmd, args...)
	if err != nil {
		log.Printf("Failed to retry command %s: %v", cmd, err)
	}
	os.Exit(exitCode)
}
