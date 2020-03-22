package commands

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestTimedOutExecution(t *testing.T) {
	cmd := exec.Command("sleep", "1")
	status, err := timedExecution(cmd, 999.99)
	assert.Equal(t, err, nil)
	assert.Equal(t, status, TO)
}

func TestTimedExecution(t *testing.T) {
	cmd := exec.Command("sleep", "1")
	status, err := timedExecution(cmd, 1250)
	assert.Equal(t, err, nil)
	assert.Equal(t, status, OK)
}
