package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveTemplateByLanguageJava(t *testing.T) {
	temp, err := ResolveTemplateByLanguage("java")
	assert.NoError(t, err)
	assert.Equal(t, JavaTemplate, temp)
}

func TestResolveTemplateByLanguagePython(t *testing.T) {
	temp, err := ResolveTemplateByLanguage("python")
	assert.NoError(t, err)
	assert.Equal(t, PythonTemplate, temp)
}

func TestResolveTemplateByLanguageCpp(t *testing.T) {
	temp, err := ResolveTemplateByLanguage("cpp")
	assert.NoError(t, err)
	assert.Equal(t, CppTemplate, temp)
}

func TestResolveTemplateByLanguageGolang(t *testing.T) {
	temp, err := ResolveTemplateByLanguage("golang")
	assert.NoError(t, err)
	assert.Equal(t, GolangTemplate, temp)
}

func TestResolveTemplateByLanguageUnknow(t *testing.T) {
	temp, err := ResolveTemplateByLanguage("dummy")
	assert.Error(t, err)
	assert.Empty(t, temp)
}
