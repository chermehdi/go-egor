package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestGetHeader(t *testing.T) {
	header := GetHeader("my header", "==", 10);
	expected := fmt.Sprintf("my header\n====================");

	assert.Equal(t, header, expected)
}

