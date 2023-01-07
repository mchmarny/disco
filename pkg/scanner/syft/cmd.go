package syft

import (
	"os/exec"
)

func MakePackagesCmd(digest, path string) *exec.Cmd {
	return exec.Command("syft", "packages", "--quiet",
		"--output", "spdx-json",
		"--file", path,
		digest) //nolint
}
