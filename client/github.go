package client

import (
	"os"

	"github.com/Murray-LIANG/gorelease"
	"github.com/apex/log"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type githubClient struct {
	client *github.Client
}

func NewGitHub(ctx *gorelease.Context) (*githubClient, error) {
	token := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ctx.GitToken})

	client := github.NewClient(oauth2.NewClient(ctx, token))
	return &githubClient{client: client}, nil
}

func (c *githubClient) CreateRelease(
	ctx *gorelease.Context,
	body string,
) (int64, error) {

	releaseInfo := &github.RepositoryRelease{
		Name:    github.String(ctx.ReleaseSubject),
		TagName: github.String(ctx.ReleaseTag),
		Body:    github.String(body),
	}
	release, _, err := c.client.Repositories.GetReleaseByTag(
		ctx,
		ctx.GitOwner,
		ctx.GitRepo,
		ctx.ReleaseTag,
	)

	if err != nil {
		release, _, err = c.client.Repositories.CreateRelease(
			ctx,
			ctx.GitOwner,
			ctx.GitRepo,
			releaseInfo,
		)
	} else {
		release, _, err = c.client.Repositories.EditRelease(
			ctx,
			ctx.GitOwner,
			ctx.GitRepo,
			release.GetID(),
			releaseInfo,
		)
	}
	log.WithField("url", release.GetHTMLURL()).Info("release updated")
	return release.GetID(), err
}

func (c *githubClient) Upload(
	ctx *gorelease.Context,
	releaseID int64,
	name string,
	file *os.File,
) error {
	_, _, err := c.client.Repositories.UploadReleaseAsset(
		ctx,
		ctx.GitOwner,
		ctx.GitRepo,
		releaseID,
		&github.UploadOptions{
			Name: name,
		},
		file,
	)
	return err
}
