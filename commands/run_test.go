package commands

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimedOutExecution(t *testing.T) {
	cmd := exec.Command("sleep", "1.1")
	status, _, err := timedExecution(cmd, 1000+TimeOutDelta)
	assert.Equal(t, err, nil)
	assert.Equal(t, status, TO)
}

func TestTimedExecution(t *testing.T) {
	cmd := exec.Command("sleep", "1")
	status, _, err := timedExecution(cmd, 1000+TimeOutDelta)
	assert.Equal(t, err, nil)
	assert.Equal(t, status, OK)
}

func TestCompress(t *testing.T) {
	got1 := strings.Repeat("a", 64)
	got2 := strings.Repeat("a", 65)
	expected2 := strings.Repeat("a", 30) + "..." + strings.Repeat("a", 31)
	assert.Equal(t, Compress(got1), got1)
	assert.Equal(t, Compress(got2), expected2)
}
