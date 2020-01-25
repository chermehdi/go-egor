package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestCurrentDirectory(t *testing.T) {
	//CreateDirectoryStructure(config.Task{});
}
