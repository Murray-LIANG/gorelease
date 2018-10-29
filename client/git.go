package client

import (
	"context"
	"errors"
	"os/exec"
	"strings"

	"github.com/apex/log"
)

type gitCli struct{}

// Run executes the git commands and returns the output and errors if any.
func (c *gitCli) Run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	log.WithField("args", args).Debug("running git cmd")

	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	log.WithField("output", output).Debug("git cmd finished")
	if err != nil {
		return "", errors.New(outputStr)
	}

	return outputStr, nil
}

// RunAndOutput executes the git commands and returns the output and errors after
// trimming the tailing \n.
func (c *gitCli) RunAndOutput(ctx context.Context, args ...string) (string, error) {
	output, err := c.Run(ctx, args...)
	output = strings.Split(output, "\n")[0]
	if err != nil {
		err = errors.New(strings.Trim(err.Error(), "\n"))
	}
	return output, err
}

var Git = &gitCli{}
