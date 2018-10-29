package gorelease

import (
	"context"
	"time"
)

type ReleaseOptions struct {
	GitOwner       string
	GitRepo        string
	GitToken       string
	ReleaseSubject string
	ReleaseTag     string
	ReleaseAuthor  string
	ReleaseEmail   string
	ReleaseNotes   []string
	ReleaseAssets  []string
}

type Context struct {
	context.Context
	ReleaseOptions
}

func NewContext(opt ReleaseOptions, timeout time.Duration) (*Context, context.CancelFunc) {
	c, cancel := context.WithTimeout(context.Background(), timeout)
	return &Context{Context: c, ReleaseOptions: opt}, cancel
}
