package commands

import (
	"github.com/chermehdi/egor/config"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
	
)

func DeleteDir(dirPath string) {
	_ = os.RemoveAll(dirPath)
}

func createDummyMetaData() (config.EgorMeta) {
	meta_data := config.EgorMeta {
		TaskName:	"Dummy Task",
		TaskLang: 	"cpp",
		Inputs: 	[]config.IoFile {
			config.IoFile {
				Name: 	"test-0",
				Path: 	"inputs/test-0.in",
				Custom:	false, 
			}, 
			config.IoFile {
				Name: 	"test-1",
				Path: 	"inputs/test-1.in",
				Custom:	true,
			}, 
		}, 
		Outputs: 	[]config.IoFile {
			config.IoFile {
				Name: 	"test-0",
				Path: 	"outputs/test-0.ans",
				Custom:	false,
			},
		},
	}

	return meta_data
}


func TestAddNewCaseInput(t *testing.T) {
	meta_data := createDummyMetaData()

	// create temp inputs directory
	_ = os.Mkdir("inputs", 0777)

	input_lines := [2]string{"Hello", "World"}
	case_name := "test-2"
	meta_data = AddNewCaseInput(input_lines, case_name, meta_data)

	assert.Equal(t, len(meta_data.Inputs), 3)
	assert.Equal(t, meta_data.Inputs[2].Name, case_name + ".in")
	assert.Equal(t, meta_data.Inputs[2].Custom, true)
	
	DeleteDir("inputs")
}


func TestAddNewCaseInput(t *testing.T) {
	meta_data := createDummyMetaData()

	// create temp outputs directory
	_ = os.Mkdir("outputs", 0777)

	input_lines := [2]string{"Hello", "World"}
	case_name := "test-2"
	meta_data = AddNewCaseOutput(input_lines, case_name, meta_data)

	assert.Equal(t, len(meta_data.Inputs), 2)
	assert.Equal(t, meta_data.Outputs[1].Name, case_name + ".ans")
	assert.Equal(t, meta_data.Outputs[1].Custom, true)

	DeleteDir("outputs")
}
