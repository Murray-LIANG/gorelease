package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Murray-LIANG/gorelease"
	"github.com/Murray-LIANG/gorelease/step"
	"github.com/apex/log"
	"github.com/caarlos0/ctrlc"
	"github.com/fatih/color"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

type Steper interface {
	fmt.Stringer

	Run(ctx *gorelease.Context) error
}

var steps = []Steper{
	&step.Prerequisite{},
	&step.Changelog{},
	&step.Publish{},
}

var timeout, _ = time.ParseDuration("30m")
var highlight = color.New(color.Bold)

func main() {

	var (
		owner   = kingpin.Flag("owner", "The owner of the github repo.").Short('o').String()
		repo    = kingpin.Flag("repo", "The name of the github repo.").Short('r').String()
		token   = kingpin.Flag("token", "The token to login into the github.").Short('t').Required().String()
		subject = kingpin.Flag("subject", "The subject of the release.").Short('s').String()
		tag     = kingpin.Flag("tag", "The tag of the release.").Short('g').String()
		notes   = kingpin.Flag("release-notes", "The release notes of the release.").Short('n').String()
		asset   = kingpin.Flag("asset", "The assets to add in the release.").Short('a').Strings()
	)

	kingpin.Version(fmt.Sprintf("%v, commit %v, built at %v", version, commit, date))
	kingpin.VersionFlag.Short('V')
	kingpin.HelpFlag.Short('h')

	kingpin.Parse()

	startTime := time.Now()
	log.Info(highlight.Sprintf(
		"publishing github release using gorelease %s...", version))

	options := gorelease.ReleaseOptions{
		GitOwner:       *owner,
		GitRepo:        *repo,
		GitToken:       *token,
		ReleaseSubject: *subject,
		ReleaseTag:     *tag,
	}
	if *notes != "" {
		options.ReleaseNotes = []string{*notes}
	}
	options.ReleaseAssets = *asset

	ctx, cancel := gorelease.NewContext(options, timeout)
	defer cancel()

	var workflow = func() error {
		for _, s := range steps {
			log.Info(highlight.Sprint(strings.ToUpper(s.String())))
			if err := s.Run(ctx); err != nil {
				return err
			}
		}
		return nil
	}

	if err := ctrlc.Default.Run(ctx, workflow); err != nil {
		log.WithError(err).Error(highlight.Sprintf(
			"failed to publish github release after %0.2fs", time.Since(startTime)))
		os.Exit(1)
		return
	}
	log.Info(highlight.Sprintf("publish succeeded after %0.2fs", time.Since(startTime)))
}
