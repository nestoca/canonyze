package internal

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanonize(t *testing.T) {

	files, err := ioutil.ReadDir("testdata")
	assert.NoError(t, err)

	for _, file := range files {
		name := file.Name()

		t.Run(name, func(t *testing.T) {
			// Input
			input, err := os.Open(filepath.Join("testdata", name, "input.yaml"))
			assert.NoError(t, err)
			defer input.Close()

			// Actual output
			var actualBuffer bytes.Buffer
			canonize(input, &actualBuffer)
			actual := actualBuffer.String()

			// Expected output
			expectedBytes, err := ioutil.ReadFile(filepath.Join("testdata", name, "expected.yaml"))
			assert.NoError(t, err)
			expected := string(expectedBytes)

			assert.Equal(t, expected, actual)
		})
	}
}
