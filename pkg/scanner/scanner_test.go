package scanner

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	file, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	fp := file.Name()

	cmd := exec.Command("touch", fp)
	err = runCmd(cmd, fp)
	assert.NoError(t, err)
}
