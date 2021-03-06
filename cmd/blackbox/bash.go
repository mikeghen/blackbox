package blackbox

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/logrusorgru/aurora"
)

// type Tracer struct{}
//
// func (t *Tracer) Write(p []byte) (n int, err error) {
// 	os.Stdout.Write(append([]byte(">"), p...))
// 	return len(p), nil
// }

func RunSync(command string, cmdArgs []string, env map[string]string, debug bool) error {
	cmd := exec.Command(command, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	setEnv(env)

	if debug {
		debugCmd := fmt.Sprintf("%s %s %s", strings.Join(formatEnv(env), " "), command, strings.Join(cmdArgs, " "))
		cmd.Stdout.Write([]byte(aurora.Cyan(debugCmd).String() + "\n"))
	}

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Println(exitError.ExitCode())
			return exitError
		}
		return err
	}

	return nil
}

func Run(command string, cmdArgs []string, env map[string]string, debug bool) error {

	// cmdArgs := strings.Fields(cmdString)
	cmd := exec.Command(command, cmdArgs...)

	setEnv(env)

	// DEBUG
	if debug {
		debugCmd := fmt.Sprintf("%s %s %s", strings.Join(formatEnv(env), " "), command, strings.Join(cmdArgs, " "))
		Trace("debug", debugCmd)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	stdoutScanner := bufio.NewScanner(stdout)
	go func() {
		for stdoutScanner.Scan() {
			Trace("info", stdoutScanner.Text())
		}
	}()

	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stderrScanner.Scan() {
			Trace("error", stderrScanner.Text())
		}
	}()

	err = cmd.Wait()
	return err
}

// ExecCommand needs to be moved outside of this package
func ExecCommand(command string, args []string, env map[string]string, debug bool) cmd.Status {
	envCmd := cmd.NewCmdOptions(
		cmd.Options{Buffered: true},
		command,
		args...,
	)
	// Set the given environment variables
	envCmd.Env = formatEnv(env)

	// DEBUG
	if debug {
		debugCmd := fmt.Sprintf("%s %s %s", strings.Join(formatEnv(env), " "), command, strings.Join(args, " "))
		Trace(aurora.Cyan(debugCmd).String())
	}

	// Run and wait for Cmd to return, discard Status
	status := <-envCmd.Start()

	// fmt.Println("exit code", strconv.Itoa(status.Exit))
	// Cmd has finished but wait for goroutine to print all lines
	for len(envCmd.Stdout) > 0 || len(envCmd.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}

	return status
}

func setEnv(env map[string]string) {
	for k, v := range env {
		os.Setenv(k, v)
	}
}

func formatEnv(env map[string]string) []string {
	var output []string
	for k, v := range env {
		output = append(output, fmt.Sprintf(`%s=%s`, k, v))
	}
	return output
}
