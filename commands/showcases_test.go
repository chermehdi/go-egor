package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetIoFilesMaps(t *testing.T) {
	metaData := createDummyMetaData()
	inputs, outputs := GetIoFilesMaps(metaData)

	assert.Equal(t, len(inputs), len(metaData.Inputs))
	assert.Equal(t, len(outputs), len(metaData.Outputs))
	assert.Equal(t, inputs["test-0"], metaData.Inputs[0])
	assert.Equal(t, inputs["test-1"], metaData.Inputs[1])
	assert.Equal(t, outputs["test-0"], metaData.Outputs[0])
}
