package zfs

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strings"
)

type OutputVersionT struct {
	Command string
	Major   int
	Minor   int
}

type ZFSVersionT struct {
	Userland string
	Kernel   string
}

type ZFSVersionOutputT struct {
	OutputVersion OutputVersionT
	ZFSVersion    ZFSVersionT
}

func (o ZFSVersionOutputT) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("output_version.command", o.OutputVersion.Command),
		slog.Int("output_version.major", o.OutputVersion.Major),
		slog.Int("output_version.minor", o.OutputVersion.Minor),
		slog.String("zfs_version.userland", o.ZFSVersion.Userland),
		slog.String("zfs_version.kernel", o.ZFSVersion.Kernel),
	)
}

func GetZFSVersion(logger *slog.Logger) (*string, error) {
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
	var o ZFSVersionOutputT
	if err := json.Unmarshal(stdo, &o); err != nil {
		return nil, fmt.Errorf("failed to read output of '%s'; output: (%w)", cmd.String(), err)
	}
	logger.Debug("ZFS Command Output", "output", o)
	return &o.ZFSVersion.Userland, nil
}
