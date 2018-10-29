package step

import (
	"errors"
	"fmt"

	"github.com/Murray-LIANG/gorelease"
	"github.com/Murray-LIANG/gorelease/client"
)

type Changelog struct{}

func (s *Changelog) String() string {
	return "step of getting changelog"
}

func (s *Changelog) Run(ctx *gorelease.Context) error {
	if ctx.ReleaseTag == "" {
		return errors.New("failed to get current tag from context")
	}
	logs, err := getChangelog(ctx, ctx.ReleaseTag)
	if err != nil {
		return err
	}
	ctx.ReleaseNotes = append(ctx.ReleaseNotes, fmt.Sprintf("## Changelog\n\n%v", logs))
	return nil
}

func getPreviousTag(ctx *gorelease.Context, tag string) (string, error) {
	prev, err := client.Git.RunAndOutput(ctx, "describe", "--tags", "--abbrev=0", tag+"^")
	if err != nil {
		prev, err = client.Git.Run(ctx, "rev-list", "--max-parents=0", "HEAD")
	}
	return prev, err
}

func getChangelog(ctx *gorelease.Context, tag string) (string, error) {
	preTag, err := getPreviousTag(ctx, tag)
	if err != nil {
		return "", err
	}
	return client.Git.Run(ctx, "log", "--pretty=oneline", "--abbrev-commit",
		"--no-decorate", "--no-color", fmt.Sprintf("%v..%v", preTag, tag))
}
