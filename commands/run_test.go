package commands

import (
	"os/exec"
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
