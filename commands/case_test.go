package commands

import (
	"github.com/chermehdi/egor/config"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func createDummyMetaData() config.EgorMeta {
	meteData := config.EgorMeta{
		TaskName: "Dummy Task",
		TaskLang: "cpp",
		Inputs: []config.IoFile{
			config.IoFile{
				Name:   "test-0",
				Path:   "inputs/test-0.in",
				Custom: false,
			},
			config.IoFile{
				Name:   "test-1",
				Path:   "inputs/test-1.in",
				Custom: true,
			},
		},
		Outputs: []config.IoFile{
			config.IoFile{
				Name:   "test-0",
				Path:   "outputs/test-0.ans",
				Custom: false,
			},
		},
	}

	return meteData
}

func TestAddNewCaseInput(t *testing.T) {
	meteData := createDummyMetaData()

	// create temp inputs directory
	_ = os.Mkdir("inputs", 0777)
	defer DeleteDir("inputs")

	inputLines := []string{"Hello", "World"}
	caseName := "test-2"
	meteData, err := AddNewCaseInput(inputLines, caseName, meteData)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(meteData.Inputs), 3)
	assert.Equal(t, meteData.Inputs[2].Name, caseName)
	assert.Equal(t, meteData.Inputs[2].Path, path.Join("inputs", caseName+".in"))
	assert.Equal(t, meteData.Inputs[2].Custom, true)

}

func TestAddNewCaseOutput(t *testing.T) {
	meteData := createDummyMetaData()

	// create temp outputs directory
	_ = os.Mkdir("outputs", 0777)
	defer DeleteDir("outputs")

	outputLines := []string{"Hello", "World"}
	caseName := "test-2"
	meteData, err := AddNewCaseOutput(outputLines, caseName, meteData)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(meteData.Outputs), 2)
	assert.Equal(t, meteData.Outputs[1].Name, caseName)
	assert.Equal(t, meteData.Outputs[1].Path, path.Join("outputs", caseName+".ans"))
	assert.Equal(t, meteData.Outputs[1].Custom, true)

}
