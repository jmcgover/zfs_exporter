package zfs

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type OutputVersion struct {
	command string
	major   int64
	minor   int64
}

type ZFSVersion struct {
	userland string
	kernel   string
}

type ZFSVersionOutput struct {
	output_version OutputVersion
	zfs_version    ZFSVersion
}

func GetZFSVersion() (*string, error) {
	cmd := exec.Command(`zpool`, `--json`, `--json-int`, `list`, `-Ho`, `name`)

	// Setup pipes
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	// command begin
	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command '%s': %w", cmd.String(), err)
	}

	// stdout
	stdo, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to read output of '%s'; output: (%w)", cmd.String(), err)
	}

	// stderr
	stde, _ := io.ReadAll(stderr)
	if err = cmd.Wait(); err != nil {
		return nil, fmt.Errorf("failed to execute command '%s'; output: '%s' (%w)", cmd.String(), strings.TrimSpace(string(stde)), err)
	}

	// unmarshal JSON into Go objects
	var o ZFSVersionOutput
	if err := json.Unmarshal(stdo, &o); err != nil {
		return nil, fmt.Errorf("failed to read output of '%s'; output: (%w)", cmd.String(), err)
	}
	return &o.zfs_version.userland, nil
}
