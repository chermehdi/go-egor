package commands

import (
	"github.com/chermehdi/egor/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createSimpleDummyMetaData() config.EgorMeta {
	meteData := config.EgorMeta{
		TaskName: "Dummy Task",
		TaskLang: "cpp",
		Inputs: []config.IoFile{
			config.IoFile{
				Name:   "test-0",
				Path:   "inputs/test-0.in",
				Custom: true,
			},
		},
		Outputs: []config.IoFile{
			config.IoFile{
				Name:   "test-0",
				Path:   "outputs/test-0.ans",
				Custom: true,
			},
		},
	}

	return meteData
}

func TestGetTestCases(t *testing.T) {
	metaData := createSimpleDummyMetaData()
	testCases := GetTestCases(metaData)
	inputs := metaData.Inputs
	outputs := metaData.Outputs

	assert.Equal(t, len(testCases), len(metaData.Inputs))
	assert.Equal(t, testCases[0].Id, inputs[0].GetId())
	assert.Equal(t, testCases[0].Name, inputs[0].Name)
	assert.Equal(t, testCases[0].Name, outputs[0].Name)
	assert.Equal(t, testCases[0].InputPath, inputs[0].Path)
	assert.Equal(t, testCases[0].OutputPath, outputs[0].Path)
	assert.Equal(t, testCases[0].Custom, inputs[0].Custom)
}
