package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/abiosoft/readline"
	"github.com/dolthub/ishell"
)

func main() {
	rlConf := readline.Config{
		Prompt:                 ">>> ",
		Stdout:                 os.Stdout,
		Stderr:                 os.Stderr,
		HistoryFile:            "./history.txt",
		HistoryLimit:           500,
		HistorySearchFold:      true,
		DisableAutoSaveHistory: true,
	}
	shellConf := ishell.UninterpretedConfig{
		ReadlineConfig: &rlConf,
		QuitKeywords: []string{
			"quit", "exit", "quit()", "exit()",
		},
		RunOnReturn:    true,
		LineTerminator: "\\",
	}

	shell := ishell.NewUninterpreted(&shellConf)
	shell.EOF(func(c *ishell.Context) {
		fmt.Println("goodbye!")
		c.Stop()
	})
	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		if count > 1 {
			c.Stop()
		} else {
			c.Println("Received SIGINT. Interrupt again to exit, or use ^D, quit, or exit")
		}
	})

	shell.Uninterpreted(func(c *ishell.Context) {
		// The entire input line is provided as the single element in c.Args
		query := c.Args[0]
		if len(strings.TrimSpace(query)) == 0 {
			return
		}

		singleLine := strings.ReplaceAll(query, "\n", " ")

		// Add this query to our command history
		if err := shell.AddHistory(singleLine); err != nil {
			shell.Println(color.RedString(err.Error()))
		}

		query = strings.TrimSuffix(query, shell.LineTerminator())

		// Execute the query on the database, then either print the query results or an error if there was one
		// print query for testing
		fmt.Println(query)

		// Update the prompts with the current database and branch name
		// shell.SetPrompt(">>> ")
		// shell.SetMultiPrompt("... ")
	})

	// Run the shell. This blocks until the user exits.
	shell.Run()
}
