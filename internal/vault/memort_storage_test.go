package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MemoryStorage_Init(t *testing.T) {
	t.Run("returns nil", func(t *testing.T) {
		s := MemoryStorage{}
		assert.Nil(t, s.Init())
	})

}

func Test_MemoryStorage_All(t *testing.T) {
	t.Run("returns a collection", func(t *testing.T) {
		expected := Records{"FOO": "bar"}
		s := MemoryStorage{
			records: expected,
		}

		assert.Equal(t, s.All(), expected)
	})
}

func Test_MemoryStorage_Exists(t *testing.T) {
	s := MemoryStorage{
		records: Records{"FOO": "bar"},
	}

	var tests = []struct {
		name     string
		input    string
		expected bool
	}{
		{"when a key does not exist", "Bar", false},
		{"when a key exists", "FOO", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, s.Exists(tt.input), tt.expected)
		})
	}
}

func Test_MemoryStorage_Get(t *testing.T) {
	s := MemoryStorage{
		records: Records{"FOO": "bar"},
	}

	var tests = []struct {
		name     string
		input    string
		expected string
	}{
		{"when a key does not exist", "BAR", ""},
		{"when a key exists", "FOO", "bar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, s.Get(tt.input), tt.expected)
		})
	}
}

func Test_MemoryStorage_Upsert(t *testing.T) {
	s := MemoryStorage{records: Records{}}
	s.Upsert("FOO", "BAR")
	assert.Equal(t, s.Get("FOO"), "BAR")
}

func Test_MemoryStorage_Destroy(t *testing.T) {
	s := MemoryStorage{records: Records{"FOO": "BAR"}}

	var tests = []struct {
		name     string
		input    string
		expected string
	}{
		{"when a key does not exists", "FOO", ""},
		{"when a key does not exists", "BAR", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Destroy(tt.input)
			assert.Equal(t, s.Get(tt.input), tt.expected)
		})
	}
}
