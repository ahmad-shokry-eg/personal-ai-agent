package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"help/ai"
	"help/ui"

	"github.com/Ingenimax/agent-sdk-go/pkg/interfaces"
)

var (
	lastCmd    string
	lastOutput string
)

func main() {
	// 1. Initialize AI Agent
	helperAgent, err := ai.NewAgent()
	if err != nil {
		fmt.Printf("Failed to initialize AI Agent: %v\n", err)
		os.Exit(1)
	}

	// 2. Open History file
	home, _ := os.UserHomeDir()
	historyPath := filepath.Join(home, ".help_history")
	historyFile, err := os.OpenFile(historyPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Warning: Failed to open history file: %v\n", err)
	} else {
		defer historyFile.Close()
	}

	reader := bufio.NewReader(os.Stdin)

	// REPL Loop
	for {
		fmt.Print("help> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)

		// Save to history
		if input != "" && historyFile != nil {
			if _, err := historyFile.WriteString(input + "\n"); err != nil {
				fmt.Printf("Warning: Failed to write to history file: %v\n", err)
			}
		}

		if input == "exit" || input == "quit" {
			break
		}

		if input == "" {
			// Launch the Bubbletea TUI
			res := ui.ShowMenu()
			handleMenuAction(res, helperAgent)
			continue
		}

		// Execute user command
		executeCommand(input)
	}
}

func executeCommand(input string) {
	lastCmd = input

	// using 'sh -c' to support pipes and redirects typed by user
	cmd := exec.Command("sh", "-c", input)
	
	// We want to capture the output for debugging, but also print it to terminal in real-time
	var outBuf bytes.Buffer
	multiOut := io.MultiWriter(os.Stdout, &outBuf)
	multiErr := io.MultiWriter(os.Stderr, &outBuf)
	
	cmd.Stdout = multiOut
	cmd.Stderr = multiErr
	
	// For standard input, we just bind it (e.g., if the user command requires input)
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		// Just store the error in the output buffer for debugger if needed
		fmt.Fprintf(&outBuf, "\nCommand Failed: %v\n", err)
	}

	lastOutput = outBuf.String()
}

func handleMenuAction(res ui.MenuResult, agent *ai.Agent) {
	switch res.Action {
	case ui.ActionQuit:
		os.Exit(0)
	case ui.ActionBack:
		return
	case ui.ActionDebug:
		if lastCmd == "" {
			fmt.Println("No command was previously executed to debug.")
			return
		}
		fmt.Printf("\nDebugging last command: `%s`\n", lastCmd)
		stream, err := agent.StreamDebug(context.Background(), lastCmd, lastOutput)
		if err != nil {
			fmt.Printf("Initial LLM error: %v\n", err)
			return
		}
		
		fmt.Println("\nAI Response:")
		for event := range stream {
			if event.Type == interfaces.AgentEventError {
				fmt.Printf("\nStream Error: %v\n", event.Error)
				break
			}
			// Just write the content chunk
			fmt.Print(event.Content)
		}
		fmt.Println()

	case ui.ActionPush:
		fmt.Println("\nGathering diff to generate commit message...")
		
		// Get diff using standard command
		diffCmd := exec.Command("git", "diff")
		diffOut, err := diffCmd.Output()
		if err != nil {
			fmt.Printf("Failed to get git diff: %v\n", err)
			return
		}

		if len(strings.TrimSpace(string(diffOut))) == 0 {
			// check cached diff
			diffCmd = exec.Command("git", "diff", "--cached")
			diffOut, _ = diffCmd.Output()
			if len(strings.TrimSpace(string(diffOut))) == 0 {
				fmt.Println("No unstaged or staged changes found. Please modify some files first.")
				return
			}
		}

		commitMsg, err := agent.GenerateCommitMessage(context.Background(), string(diffOut))
		if err != nil {
			fmt.Printf("Failed to generate commit message: %v\n", err)
			return
		}

		fmt.Printf("\nGenerated commit message: \n%s\n\n", commitMsg)
		
		// Run git add .
		fmt.Println("Running `git add .`...")
		if err := executeShell("git add ."); err != nil {
			return
		}

		// Run git commit
		fmt.Printf("Running `git commit -m \"%s\"`...\n", commitMsg)
		if err := executeShell(fmt.Sprintf("git commit -m \"%s\"", commitMsg)); err != nil {
			return
		}

		// Run git push
		fmt.Println("Running `git push`...")
		executeShell("git push")

	case ui.ActionBuild:
		if res.BuildCmd != "" {
			fmt.Printf("\nBuilding code: %s\n", res.BuildCmd)
			executeCommand(res.BuildCmd)
		}
	}
}

func executeShell(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}