package cli_test

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/fatt/cmd/fatt/cli"
)

func TestVersionCliText(t *testing.T) {
	assert := assert.New(t)

	expected := fmt.Sprintf(
		`GitVersion:    devel
GitCommit:     unknown
GitTreeState:  unknown
BuildDate:     unknown
GoVersion:     %s
Compiler:      %s
Platform:      %s/%s
`, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)

	output, err := executeCommand(cli.NewVersionCommand())
	assert.NoError(err)
	assert.Equal(expected, output)
}

func TestVersionCliJSON(t *testing.T) {
	assert := assert.New(t)

	expected := fmt.Sprintf(`{
  "git_version": "devel",
  "git_commit": "unknown",
  "git_tree_state": "unknown",
  "build_date": "unknown",
  "go_version": "%s",
  "compiler": "%s",
  "platform": "%s/%s"
}
`, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
	a := []string{"-o", "json"}
	output, err := executeCommand(cli.NewVersionCommand(), a...)
	assert.NoError(err)
	assert.Equal(expected, output)

	a = []string{"--output-format", "json"}
	output, err = executeCommand(cli.NewVersionCommand(), a...)
	assert.NoError(err)
	assert.Equal(expected, output)
}

func executeCommand(cmd *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(cmd, args...)
	return output, err
}

func executeCommandC(cmd *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs(args)

	c, err = cmd.ExecuteC()

	return c, buf.String(), err
}
