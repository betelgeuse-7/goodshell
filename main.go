package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	ROOT_PROMPT    = "gsh#~ "
	DEFAULT_PROMPT = "gsh~ "
	ARROW_KEY_UP   = "\033[A"
	ARROW_KEY_DOWN = "\033[B"
)

type GoodShell struct {
	History *CommandChain
}

// return "gsh#~" if the user is root; else "gsh~".
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
	return fmt.Sprintf("\033[33m{%d:%s}\033[0m", hour, minuteStr)
}

func main() {
	gsh := &GoodShell{History: NewCommandChain("")}

	prompt := getPrompt()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s %s", getTime(), prompt)

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

		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		cmdChain := NewCommandChain(cmd.String())
		gsh.History.Add(cmdChain)
	}
}
