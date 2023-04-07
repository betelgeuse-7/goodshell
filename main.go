package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	ROOT_PROMPT    = "gsh#~ "
	DEFAULT_PROMPT = "gsh~ "
)

type ErrNotInPath struct {
	msg string
}

func (e ErrNotInPath) Error() string { return e.msg }

type GoodShell struct {
	History *CommandChain
}

// return "gsh#~" if the user is root; else "gsh~".
// also include cwd.
func getPrompt() string {
	prompt := DEFAULT_PROMPT
	if os.Getuid() == 0 {
		prompt = ROOT_PROMPT
	}
	return prompt
}

// format the current hour and minute, and return it
func getTime() string {
	now := time.Now()
	hour, minute := now.Hour(), now.Minute()
	minuteStr := fmt.Sprintf("%d", minute)
	if minute < 10 {
		minuteStr = "0" + minuteStr
	}
	return fmt.Sprintf("\033[33m(%d:%s)\033[0m", hour, minuteStr)
}

func getCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		reportErrFatal("Unable to get current working directory\n")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		reportErrFatal("Unable to get $HOME\n")
	}
	if strings.HasPrefix(cwd, homeDir) {
		cwd = strings.Replace(cwd, homeDir, "~", 1)
	}
	return cwd
}

func makeHeadline() string {
	var headline strings.Builder
	headline.WriteString(getTime())
	headline.WriteByte(' ')
	headline.WriteString(getCwd())
	headline.WriteString(" • ")
	headline.WriteString(getPrompt())
	headline.WriteByte(' ')
	return headline.String()
}

func main() {
	gsh := &GoodShell{}
	gsh.REPL()
}

// Read-Evaluate-Print Loop
func (g *GoodShell) REPL() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(makeHeadline())
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		command = strings.TrimSuffix(command, "\n")
		if command == "exit" {
			break
		}

		parts := strings.Fields(command)
		if len(parts) == 0 {
			continue
		}

		programName, argv := parts[0], parts[1:]
		if programName == "cd" {
			if err := cd(argv); err != nil {
				reportErr(err.Error())
			}
			continue
		}

		exeLoc, err := getAbsoluteExeLoc(programName)
		if err != nil {
			if errors.Is(ErrNotInPath{}, err) {
				reportErr("The program %s could not be found in $PATH\n", programName)
			} else {
				reportErr("Could not get location of %s", programName)
			}
			continue
		}
		err = checkFileIsExecutable(exeLoc)
		if err != nil {
			reportErr(err.Error())
			continue
		}

		pid, err := syscall.ForkExec(exeLoc, argv, &syscall.ProcAttr{
			Env: os.Environ(),
			Files: []uintptr{
				os.Stdout.Fd(),
				os.Stdin.Fd(),
				os.Stderr.Fd(),
			},
		})
		if err != nil {
			reportErrFatal("Failed to execute command '%s'\n", strings.Join(parts, " "))
		}

		var waitStatus syscall.WaitStatus
		_, err = syscall.Wait4(pid, &waitStatus, 0, nil)
		if err != nil {
			reportErrFatal("Failed to wait for process %d\n", pid)
		}

		if code := waitStatus.ExitStatus(); code > 0 {
			fmt.Printf("Process exited with status code %d\n", waitStatus.ExitStatus())
		}
	}
}

// return the absolute file path to the program. check if the programName is a
// relative path; or if it should be searched for in $PATH.
func getAbsoluteExeLoc(programName string) (string, error) {
	isRelativePath := strings.Contains(programName, "/")
	if isRelativePath {
		return filepath.Abs(programName)
	}

	exeLoc, err := getExeLocationFromPath(programName)
	if err != nil {
		return "", err
	}
	return exeLoc, nil
}

// write formatted error message to stderr, and exit the process (the shell) with
// exit code 1.
func reportErrFatal(msgf string, args ...any) {
	fmt.Fprintf(os.Stderr, msgf, args...)
	os.Exit(1)
}

func reportErr(msgf string, args ...any) {
	fmt.Fprintf(os.Stderr, msgf, args...)
}

func checkFileIsExecutable(exeLoc string) error {
	info, err := os.Stat(exeLoc)
	if err != nil {
		if !(os.IsNotExist(err)) {
			return fmt.Errorf("This should not be printed. This is a bug.\n%s\n", err)
		}
		return fmt.Errorf("No such file or directory: %s\n", exeLoc)
	}
	// check if file is executable by the owner, group, or other users.
	isExe := info.Mode()&0111 != 0
	if !(isExe) {
		return fmt.Errorf("%s is not an executable program\n", exeLoc)
	}
	return nil
}

func getExeLocationFromPath(programName string) (string, error) {
	path := os.Getenv("PATH")
	locs := strings.Split(path, ":")
	for _, l := range locs {
		entries, err := os.ReadDir(l)
		if err != nil {
			return "", err
		}
		for _, e := range entries {
			potentialExe_Name := e.Name()
			if potentialExe_Name == programName {
				exePath := fmt.Sprintf("%s/%s", l, potentialExe_Name)
				return exePath, nil
			}
		}
	}
	return "", ErrNotInPath{}
}
