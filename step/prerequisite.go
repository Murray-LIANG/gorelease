package step

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Murray-LIANG/gorelease"
	"github.com/Murray-LIANG/gorelease/client"
)

type Prerequisite struct{}

var (
	ErrNotUnderGitFolder     = errors.New("current folder is not a git work tree")
	ErrNotFoundAnyTags       = errors.New("no tags found in the repo")
	ErrGitNotSetOriginRemote = errors.New("the repo has no `origin` remote")
)

func (s *Prerequisite) String() string {
	return "step of preparing prerequisites"
}

func (s *Prerequisite) Run(ctx *gorelease.Context) error {
	if isOptionsMissing(ctx) {
		if !isUnderGitFolder(ctx) {
			return ErrNotUnderGitFolder
		}
		if err := setRequiredOptions(ctx); err != nil {
			return err
		}
	}

	if isOptionsMissing(ctx) {
		return fmt.Errorf(
			"required options missing, GitOwner: %v, GitRepo: %v, ReleaseTag: %v",
			ctx.GitOwner, ctx.GitRepo, ctx.ReleaseTag)
	}

	if ctx.ReleaseSubject == "" {
		ctx.ReleaseSubject = strings.TrimLeft(ctx.ReleaseTag, "rv")
	}
	return nil
}

func isOptionsMissing(ctx *gorelease.Context) bool {
	return ctx.GitOwner == "" || ctx.GitRepo == "" || ctx.ReleaseTag == ""
}

func isUnderGitFolder(ctx *gorelease.Context) bool {
	output, err := client.Git.RunAndOutput(ctx, "rev-parse", "--is-inside-work-tree")
	return err == nil && output == "true"
}

func getCurrentTag(ctx *gorelease.Context) (string, error) {
	tag, err := client.Git.RunAndOutput(ctx, "describe", "--tags", "--abbrev=0")
	if err != nil {
		return "", ErrNotFoundAnyTags
	}
	return tag, nil
}

func getRepo(ctx *gorelease.Context) (string, string, error) {
	output, err := client.Git.RunAndOutput(ctx, "config", "--get", "remote.origin.url")
	if err != nil {
		return "", "", ErrGitNotSetOriginRemote
	}
	parts := strings.Split(output, "/")
	return parts[len(parts)-2], parts[len(parts)-1], nil
}

func setRequiredOptions(ctx *gorelease.Context) error {
	var err error
	if ctx.GitOwner == "" || ctx.GitRepo == "" {
		if ctx.GitOwner, ctx.GitRepo, err = getRepo(ctx); err != nil {
			return err
		}
	}
	if ctx.ReleaseTag == "" {
		if ctx.ReleaseTag, err = getCurrentTag(ctx); err != nil {
			return err
		}
	}
	return nil
}
