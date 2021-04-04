package utils

import (
	"os"
	"path"
)

// CreateTempFile creates a file named `name` in a os temporary directory
// On success, the function will return the reference to the created file, otherwise
// an non-nil error is returned instead.
func CreateTempFile(name string) (*os.File, error) {
	tmp := os.TempDir()
	file, err := os.OpenFile(path.Join(tmp, name), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}
	return file, nil
}
