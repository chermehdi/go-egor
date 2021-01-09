package commands

import (
	"github.com/chermehdi/egor/config"
	"github.com/chermehdi/egor/utils"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestTimedOutExecution(t *testing.T) {
	cmd := exec.Command("sleep", "1.1")
	status, _, err := utils.ExecuteWithTimeout(cmd, 1000+config.TimeOutDelta)
	assert.Equal(t, err, nil)
	assert.Equal(t, status, config.TO)
}

func TestTimedExecution(t *testing.T) {
	cmd := exec.Command("sleep", "1")
	status, _, err := utils.ExecuteWithTimeout(cmd, 1000+config.TimeOutDelta)
	assert.Equal(t, err, nil)
	assert.Equal(t, status, config.OK)
}

func TestChecker_TokenBasedCheckerShouldIgnoreSpaces(t *testing.T) {
	checker := &TokenChecker{}
	assert.Nil(t, checker.Check("a b c", "a  b c"))
}

func TestChecker_TokenBasedCheckerShouldIgnoreNewLines(t *testing.T) {
	checker := &TokenChecker{}
	assert.Nil(t, checker.Check("a b c", "a  b \nc"))
}

func TestChecker_TokenBasedCheckerShouldFailOnMismatchedNumberOfTokens(t *testing.T) {
	checker := &TokenChecker{}
	assert.NotNil(t, checker.Check("a b c", "ab d"))
}

func TestChecker_TokenBasedCheckerShouldFailOnDifferentToken(t *testing.T) {
	checker := &TokenChecker{}
	assert.NotNil(t, checker.Check("a b c", "a b d"))
}
