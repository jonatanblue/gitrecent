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
	"strconv"
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
				var err error

				cmd := exec.Command("git", "branch", "--sort=committerdate")
				var out bytes.Buffer
				cmd.Stdout = &out
				err = cmd.Run()
				if err != nil {
					return errors.WithStack(err)
				}

				var branches []string
				for _, b := range strings.Split(out.String(), "\n") {
					if b != "" {
						branches = append(branches, b)
					}
				}

				count := DEFAULT_COUNT
				// Adjust default down if less than number of branches
				if count > len(branches) {
					count = len(branches)
				}

				// Override default if user has set the count flag
				s := context.String(flagCount)
				if s != "" {
					count, err = strconv.Atoi(s)
					if err != nil {
						return errors.WithStack(err)
					}
				}

				if count < 1 || count > len(branches) {
					return errors.New(fmt.Sprintf("count must be between 1 and %d", len(branches)))
				}
				branches = branches[len(branches)-count:]

				// Reverse order to get most recent branch on top
				for i, j := 0, len(branches)-1; i < j; i, j = i+1, j-1 {
					branches[i], branches[j] = branches[j], branches[i]
				}

				prompt := &survey.Select{
					Message: "Pick branch to checkout",
					Options: branches,
				}
				var branch string
				err = survey.AskOne(prompt, &branch)
				if err != nil {
					return errors.WithStack(err)
				}

				// Sanitise branch names from special characters
				branch = strings.Trim(branch, "*+ ")

				cmd = exec.Command("git", "checkout", branch)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					return errors.WithStack(err)
				}

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
