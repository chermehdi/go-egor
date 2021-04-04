package commands

import (
	"bytes"
	"os"
	"testing"

	"github.com/chermehdi/egor/config"
	"github.com/chermehdi/egor/utils"
	"github.com/stretchr/testify/assert"
)

type event struct {
	Type string
	Lang string
}

type mockProvider struct {
	events     []event
	CppRunner  utils.CodeRunner
	JavaRunner utils.CodeRunner
}

type execEvent struct {
	Step string
	Ctx  *utils.ExecutionContext
}

type mockCodeRunner struct {
	Lang       string
	compEvents []execEvent
	runEvents  []execEvent
}

func (m *mockCodeRunner) Compile(ctx *utils.ExecutionContext) (*utils.ExecutionResult, error) {
	m.compEvents = append(m.compEvents, execEvent{Step: "compile", Ctx: ctx})
	return &utils.ExecutionResult{}, nil
}

func (m *mockCodeRunner) Run(ctx *utils.ExecutionContext) (*utils.ExecutionResult, error) {
	m.runEvents = append(m.runEvents, execEvent{Step: "Run", Ctx: ctx})
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	return &utils.ExecutionResult{
		Stdout: stdout,
		Stderr: stderr,
	}, nil
}

func (m *mockProvider) resolve(lang string) (utils.CodeRunner, bool) {
	m.events = append(m.events, event{Type: "Resolve", Lang: lang})
	switch lang {
	case "java":
		return m.JavaRunner, true
	case "cpp":
		return m.CppRunner, true
	}
	panic("Should not get here")
}

func TestBatch_CanShouldSkipBatchTest(t *testing.T) {
	pr := &mockProvider{
		events: make([]event, 0),
	}
	err := runBatchInternal(1, &config.EgorMeta{}, pr.resolve)

	assert.Nil(t, err)
	assert.Len(t, pr.events, 0)
}

func TestBatch_CanRunBatchTest(t *testing.T) {
	crunner := &mockCodeRunner{
		compEvents: make([]execEvent, 0),
		runEvents:  make([]execEvent, 0),
	}
	jrunner := &mockCodeRunner{
		compEvents: make([]execEvent, 0),
		runEvents:  make([]execEvent, 0),
	}

	pr := &mockProvider{
		events:     make([]event, 0),
		CppRunner:  crunner,
		JavaRunner: jrunner,
	}

	file, err := utils.CreateTempFile("generate.cc")
	assert.Nil(t, err)
	defer os.Remove(file.Name())

	ntests := 10

	err = runBatchInternal(ntests, &config.EgorMeta{
		BatchFile: file.Name(),
		TaskLang:  "java",
	}, pr.resolve)

	assert.Nil(t, err)
	assert.Len(t, pr.events, 2)

	// two compilation events expected (generator, brute-force solution)
	assert.Len(t, crunner.compEvents, 2)

	// Only one compilation event, the main solution compilation.
	assert.Len(t, jrunner.compEvents, 1)

	// ran generator and brute solution
	assert.Len(t, crunner.runEvents, 2*ntests)

	// Only the main solution
	assert.Len(t, jrunner.runEvents, ntests)
}
