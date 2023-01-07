package trivy

import (
	"os/exec"
)

func MakeLicenseCmd(digest, path string) *exec.Cmd {
	return exec.Command("trivy", "image", "--quiet", "--security-checks", "license", "--format", "json", "--no-progress", "--output", path, digest) //nolint
}

func MakeVulnerabilityCmd(digest, path string) *exec.Cmd {
	return exec.Command("trivy", "image", "--quiet", "--security-checks", "vuln", "--format", "json", "--no-progress", "--output", path, digest) //nolint
}

func MakePackageCmd(digest, path string) *exec.Cmd {
	return exec.Command("trivy", "image", "--quiet", "--format", "spdx-json", "--no-progress", "--output", path, digest) //nolint
}
