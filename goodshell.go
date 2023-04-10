package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type GoodShell struct {
	isRoot bool
}

type Command struct {
	Program, ProgramAbsPath string
	Argv                    []string
}

type Pipe struct {
	Left, Right *Command
}

func (g *GoodShell) ReadEvalPrintLoop() {
	for {
		g.printHeadline()
		line := g.read()
		cmdx := g.parse(line)
		g.eval(cmdx)
	}
}

func NewGoodShell() *GoodShell {
	g := &GoodShell{}
	u, err := user.Current()
	if err != nil {
		fatal("Couldn't get current user")
	}
	g.isRoot = u.Username == "root"
	return g
}

func fatal(msgf string, args ...any) {
	fmt.Fprintf(os.Stderr, msgf+"\n", args...)
	syscall.Exit(1)
}

func errorf(msgf string, args ...any) {
	fmt.Fprintf(os.Stderr, msgf+"\n", args...)
}

func (g *GoodShell) eval(cmdx []Command) {
	for _, c := range cmdx {
		exeLoc := g.getAbsoluteExePath(c.Program)
		if exeLoc == "" {
			continue
		}
		if !(g.isFileExecutable(exeLoc)) {
			fatal("Specified file %s is not an executable program", c.Program)
		}
		c.ProgramAbsPath = exeLoc
		// TODO Pipes
		g.run(c)
	}
}

func (g *GoodShell) run(cmd Command) {
	pid, err := syscall.ForkExec(cmd.ProgramAbsPath, cmd.Argv, &syscall.ProcAttr{
		Env: os.Environ(),
		Files: []uintptr{
			os.Stdout.Fd(),
			os.Stdin.Fd(),
			os.Stderr.Fd(),
		},
	})
	if err != nil {
		fatal("Failed to execute command '%s'", cmd.Program+strings.Join(cmd.Argv, " "))
	}

	var waitStatus syscall.WaitStatus
	_, err = syscall.Wait4(pid, &waitStatus, 0, nil)
	if err != nil {
		fatal("Failed to wait for process %d", pid)
	}

	if code := waitStatus.ExitStatus(); code > 0 {
		fmt.Printf("\033[1;31mProcess exited with status code %d\033[0m\n", waitStatus.ExitStatus())
	}
}

func (g *GoodShell) parse(line string) []Command {
	var res []Command
	parts := strings.Split(line, "|")
	for _, p := range parts {
		cmd := Command{}
		p = strings.TrimSpace(p)
		splitted := strings.SplitAfterN(p, " ", 2)
		cmd.Program = strings.TrimSpace(splitted[0])
		cmd.Argv = splitted[1:]
		res = append(res, cmd)
	}
	return res
}

func (g *GoodShell) read() string {
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		if errors.Is(io.EOF, err) {
			fatal("Unexpected end of file")
		}
		fatal("Unknown error: %s", err)
	}
	line = strings.TrimSuffix(line, "\n")
	return line
}

func (g *GoodShell) makeHeadlineHour() string {
	hour := time.Now().Hour()
	min := time.Now().Minute()
	if min < 10 {
		// add 0 prefix to min
		return fmt.Sprintf("\033[33m(%d:0%d)\033[0m", hour, min)
	}
	return fmt.Sprintf("\033[33m(%d:%d)\033[0m", hour, min)
}

func (g *GoodShell) makeHeadlineCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		fatal("Unable to get current working directory")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fatal("Unable to get $HOME")
	}
	if strings.HasPrefix(cwd, homeDir) {
		cwd = strings.Replace(cwd, homeDir, "~", 1)
	}
	return cwd
}

func (g *GoodShell) makeHeadline() string {
	hour := g.makeHeadlineHour()
	cwd := g.makeHeadlineCwd()
	const (
		prompt     = "gsh~ "
		rootPrompt = "gsh#~ "
	)
	var res strings.Builder
	res.WriteString(hour)
	res.WriteByte(' ')
	res.WriteString(cwd)
	res.WriteString(" • ")
	if g.isRoot {
		res.WriteString(rootPrompt)
		return res.String()
	}
	res.WriteString(prompt)
	return res.String()
}

func (g *GoodShell) printHeadline() {
	fmt.Print(g.makeHeadline())
}

func (g *GoodShell) getProgramLocFromPath(programName string) string {
	path := os.Getenv("PATH")
	locs := strings.Split(path, ":")
	for _, l := range locs {
		entries, err := os.ReadDir(l)
		if err != nil {
			fatal("Unknown error: %s", err.Error())
		}
		for _, e := range entries {
			potentialExe_Name := e.Name()
			if potentialExe_Name == programName {
				exePath := fmt.Sprintf("%s/%s", l, potentialExe_Name)
				return exePath
			}
		}
	}
	errorf("Program %s could not be found in $PATH", programName)
	return ""
}

func (g *GoodShell) isFileExecutable(loc string) bool {
	info, err := os.Stat(loc)
	if err != nil {
		if !(os.IsNotExist(err)) {
			fatal("Unknown error: %s", err.Error())
		}
		fatal("No such file or directory: %s", loc)
	}
	// check if file is executable by the owner, group, or other users.
	isExe := info.Mode()&0111 != 0
	if !(isExe) {
		fatal("%s is not an executable program", loc)
	}
	return true
}

func (g *GoodShell) getAbsoluteExePath(programName string) string {
	isRelativePath := strings.Contains(programName, "/")
	if isRelativePath {
		absPath, err := filepath.Abs(programName)
		if err != nil {
			fatal("Unknown error: %s", err.Error())
		}
		return absPath
	}

	exeLoc := g.getProgramLocFromPath(programName)
	return exeLoc
}
