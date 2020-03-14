package commands

import (
	"os"
	"path"
	"testing"

	"github.com/chermehdi/egor/config"
	"github.com/stretchr/testify/assert"
)

func getDummyTaskJson() string {
	return `{
		"group" : "Educational Codeforces Round 75 (Rated for Div. 2)",
			"interactive" : false,
			"output" : {
			"type" : "stdout"
		},
		"timeLimit" : 2000,
			"name" : "C. Minimize The Integer",
			"url" : "https://codeforces.com/contest/1251/problem/C",
			"testType" : "single",
			"languages" : {
			"java" : {
				"mainClass" : "Main",
					"taskClass" : "CMinimizeTheInteger"
			}
		},
		"input" : {
			"type" : "stdin"
		},
		"memoryLimit" : 256,
			"tests" : [{ "input" : "3\n0709\n1337\n246432\n","output" : "0079\n1337\n234642\n"}]
	}`
}

func TestExtractTaskFromString(t *testing.T) {
	task, _ := extractTaskFromJson(getDummyTaskJson())
	assert.Equal(t, task.Name, "C. Minimize The Integer")
	assert.Equal(t, task.Group, "Educational Codeforces Round 75 (Rated for Div. 2)")
	assert.Equal(t, task.Url, "https://codeforces.com/contest/1251/problem/C")
	assert.Equal(t, task.Interactive, false)
	assert.Equal(t, int(task.MemoryLimit), 256)
	assert.Equal(t, task.TestType, "single")

	assert.Equal(t, len(task.Tests), 1)
	assert.Equal(t, task.Tests[0].Input, "3\n0709\n1337\n246432\n")
	assert.Equal(t, task.Tests[0].Output, "0079\n1337\n234642\n")

	assert.Equal(t, task.Input.Type, "stdin")
	assert.Equal(t, task.Output.Type, "stdout")

	assert.Equal(t, len(task.Languages), 1)
	assert.Equal(t, task.Languages["java"].MainClass, "Main")
	assert.Equal(t, task.Languages["java"].TaskClass, "CMinimizeTheInteger")
}

func DeleteDir(dirPath string) {
	_ = os.RemoveAll(dirPath)
}

func TestCreateDirectoryStructure(t *testing.T) {
	task := config.Task{
		Name:        "Dummy task",
		Group:       "Codeforces",
		Url:         "http://codeforces.com/test",
		Interactive: false,
		MemoryLimit: 250,
		TimeLimit:   100,
		Tests: []config.TestCase{
			{Input: "1", Output: "2"},
		},
		TestType:  "",
		Input:     config.IOType{Type: "stdin"},
		Output:    config.IOType{Type: "stdout"},
		Languages: nil,
	}
	configuration := config.Config{
		Lang: struct {
			Default string `yaml:"default"`
		}{Default: "cpp"},
		ConfigFileName: "egor-meta.json",
		Version:        "1.0",
		Author:         "MaxHeap",
	}
	rootDir := path.Join(os.TempDir(), "egor")
	CreateDirectory(rootDir)
	defer DeleteDir(rootDir)

	_, err := CreateDirectoryStructure(task, configuration, rootDir)
	assert.NoError(t, err)

	taskDir := path.Join(rootDir, task.Name)

	assert.FileExists(t, path.Join(taskDir, configuration.ConfigFileName))
	assert.DirExists(t, path.Join(taskDir, "inputs"))
	assert.DirExists(t, path.Join(taskDir, "outputs"))
	assert.FileExists(t, path.Join(taskDir, "main.cpp"))
	assert.FileExists(t, path.Join(taskDir, "inputs", "test-0.in"))
	assert.FileExists(t, path.Join(taskDir, "outputs", "test-0.ans"))
}

func CreateDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0777); err != nil {
			return err
		}
	}
	return nil
} 
