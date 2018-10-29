package client

import (
	"context"
	"testing"

	"github.com/Murray-LIANG/gorelease"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ctx := &gorelease.Context{Context: context.Background()}
	output, err := Git.Run(ctx, "status")
	assert.NoError(t, err)
	assert.NotEmpty(t, output)

	output, err = Git.Run(ctx, "cmd-not-exist")
	assert.Error(t, err)
	assert.Empty(t, output)
	assert.Equal(
		t,
		"git: 'cmd-not-exist' is not a git command. See 'git --help'.\n",
		err.Error(),
	)
}

func TestRunAndOutput(t *testing.T) {
	ctx := &gorelease.Context{Context: context.Background()}
	output, err := Git.RunAndOutput(ctx, "rev-parse", "--is-inside-work-tree")
	assert.NoError(t, err)
	assert.Equal(t, "true", output)

	output, err = Git.RunAndOutput(ctx, "cmd-not-exist")
	assert.Error(t, err)
	assert.Empty(t, output)
	assert.Equal(
		t,
		"git: 'cmd-not-exist' is not a git command. See 'git --help'.",
		err.Error(),
	)

}
