package templates

import (
	"github.com/chermehdi/egor/config"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TempJavaTemplate = `
	import java.util.*;
	class Solution {
		public static void main(String[] args) {
		}
	}
`
const TempCppTemplate = `
	#include <bits/stdc++.h>
	using namespace std;

	int main() {
		ios_base::sync_with_stdio(false);
		cin.tie(0);
    }
`

func TestResolveTemplateByLanguageJava(t *testing.T) {
	testResolveTemplate(t, "java", JavaTemplate)
}

func TestResolveTemplateByLanguagePython(t *testing.T) {
	testResolveTemplate(t, "python", PythonTemplate)
}

func TestResolveTemplateByLanguageCpp(t *testing.T) {
	testResolveTemplate(t, "cpp", CppTemplate)
}

func TestResolveTemplateByLanguageUnknown(t *testing.T) {
	//testResolveTemplate(t, "dummy")
	configuration := config.Config{
		Lang: struct {
			Default string `yaml:"default"`
		}{Default: "dummy"},
	}
	temp, err := ResolveTemplateByLanguage(configuration)
	assert.Error(t, err)
	assert.Empty(t, temp)
}

func TestResolveTemplateByLanguage_CustomTemplateCpp(t *testing.T) {
	testResolveFromMap(t, "cpp", CppTemplate)
}

func TestResolveTemplateByLanguage_CustomTemplateJava(t *testing.T) {
	testResolveFromMap(t, "java", TempJavaTemplate)
}

func testResolveFromMap(t *testing.T, lang string, expectedTemplate string) {
	tempDir := os.TempDir()
	filePath := path.Join(tempDir, "config_template.template")
	err := ioutil.WriteFile(filePath, []byte(expectedTemplate), 0777)
	assert.NoError(t, err)
	defer os.Remove(filePath)

	configuration := config.Config{
		Lang: struct {
			Default string `yaml:"default"`
		}{Default: lang},
		CustomTemplate: map[string]string{lang: filePath},
	}
	fileTemplate, err := ResolveTemplateByLanguage(configuration)
	assert.NoError(t, err)
	assert.Equal(t, expectedTemplate, fileTemplate)
}

func testResolveTemplate(t *testing.T, lang, expected string) {
	configuration := config.Config{
		Lang: struct {
			Default string `yaml:"default"`
		}{Default: lang},
	}
	temp, err := ResolveTemplateByLanguage(configuration)
	assert.NoError(t, err)
	assert.Equal(t, temp, expected)
}
