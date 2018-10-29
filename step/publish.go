package step

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Murray-LIANG/gorelease"
	"github.com/Murray-LIANG/gorelease/client"
)

type Publish struct{}

func (s *Publish) String() string {
	return "step of publishing a new release"
}

func (s *Publish) Run(ctx *gorelease.Context) error {
	githubClient, err := client.NewGitHub(ctx)
	if err != nil {
		return err
	}
	releaseID, err := githubClient.CreateRelease(
		ctx, strings.Join(ctx.ReleaseNotes, "\n"))
	if err != nil {
		return err
	}

	for _, asset := range ctx.ReleaseAssets {
		_, fileName := filepath.Split(asset)
		file, err := os.Open(asset)
		if err != nil {
			return err
		}
		defer file.Close()
		if err = githubClient.Upload(ctx, releaseID, fileName, file); err != nil {
			return err
		}
	}

	return nil
}
