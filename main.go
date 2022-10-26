package main

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	cli "github.com/urfave/cli"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const DEFAULT_COUNT = 8

const flagCount = "count"

func main() {
	rand.Seed(time.Now().UnixNano())

	cliApp := cli.NewApp()
	cliApp.Name = "gr"
	cliApp.Usage = "git recent"
	cliApp.Description = "Find and checkout recently used git branches"
	cliApp.EnableBashCompletion = true
	cliApp.Commands = cli.Commands{
		{
			Name:        "checkout",
			Description: "Select recent branch to checkout",
			Aliases:     []string{"c"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     flagCount,
					Usage:    "Number of recent branches to show",
					Required: false,
				},
			},
			Action: func(context *cli.Context) error {
				count := context.Int(flagCount)
				if count < 1 {
					count = DEFAULT_COUNT
				}
				cmd := exec.Command("git", "branch", "--sort=committerdate")
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					return errors.WithStack(err)
				}

				branches := strings.Split(out.String(), "\n")
				branches = branches[len(branches)-count : len(branches)-1]

				prompt := &survey.Select{
					Message: "Pick branch to checkout",
					Options: branches,
				}
				var branch string
				err = survey.AskOne(prompt, &branch)
				if err != nil {
					return errors.WithStack(err)
				}

				fmt.Printf("selected branch %s", branch)

				return nil
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
