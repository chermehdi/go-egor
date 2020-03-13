package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTestCase(t *testing.T) {
	metaData := createSimpleDummyMetaData()
	testCase := GetTestCase(metaData, 0)

	assert.Equal(t, testCase.Id, 0)
	assert.Equal(t, testCase.Id, metaData.Inputs[0].GetId())
	assert.Equal(t, testCase.InputPath, metaData.Inputs[0].Path)
	assert.Equal(t, testCase.OutputPath, metaData.Outputs[0].Path)
	assert.Equal(t, testCase.Custom, metaData.Inputs[0].Custom)
}