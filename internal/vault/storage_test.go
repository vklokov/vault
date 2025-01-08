package vault

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newStorage(t *testing.T) {
	var tests = []struct {
		name     string
		input    Config
		expected Storable
	}{
		{"with file storage", Config{Storage: "file"}, &FileStorage{}},
		{"with memory storage", Config{Storage: "memory"}, &MemoryStorage{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStorage(tt.input)
			assert.Equal(t, reflect.TypeOf(s), reflect.TypeOf(tt.expected))
		})
	}
}
