package commands

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/chermehdi/egor/config"
	"github.com/stretchr/testify/assert"
)

func TestCopy_CanGenerateTaskFile(t *testing.T) {
	tests := []struct {
		CppLibraryLocation string
		TaskFile           string
		Lang               string
		Expected           string
	}{
		{
			CppLibraryLocation: "",
			TaskFile:           "Main.java",
			Lang:               "java",
			Expected:           "Main.java",
		},
		{
			CppLibraryLocation: "/temp/lib",
			TaskFile:           "Main.java",
			Lang:               "java",
			Expected:           "Main.java",
		},
		{
			CppLibraryLocation: "/temp/lib",
			TaskFile:           "main.cpp",
			Lang:               "cpp",
			Expected:           "main_gen.cpp",
		},
	}

	for _, test := range tests {
		conf := &config.Config{
			CppLibraryLocation: test.CppLibraryLocation,
		}
		meta := &config.EgorMeta{
			TaskFile: test.TaskFile,
			TaskLang: test.Lang,
		}
		taskFile := getGenFile(meta, conf)
		assert.Equal(t, test.Expected, taskFile)
	}
}

func TestCopy_CanCopyToClipboard(t *testing.T) {
	tmp := os.TempDir()
	sol := path.Join(tmp, "main.cc")
	data := []byte("This is a test for writing to the clipboard")
	err := ioutil.WriteFile(sol, data, 0777)
	assert.Nil(t, err)
	defer os.Remove(sol)

	err = copyToClipboard(sol)
	assert.Nil(t, err)

	cpdata, err := clipboard.ReadAll()
	assert.Nil(t, err)

	assert.Equal(t, string(data), cpdata)
}
