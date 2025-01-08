package vault

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Controller_New(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected Storable
	}{
		{"with file storage", "file", &FileStorage{}},
		{"with default storage", "memory", &MemoryStorage{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{Storage: tt.input}
			vault := New(config)

			assert.Equal(t, reflect.TypeOf(tt.expected), reflect.TypeOf(vault.storage))
		})
	}
}

func Test_Controller_All(t *testing.T) {
	var tests = []struct {
		name     string
		expected Records
	}{
		{"with records", Records{"FOO": "bar"}},
		{"without records", Records{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockStorage)
			mockStorage.On("All").Return(tt.expected)

			v := &Vault{storage: mockStorage}

			assert.Equal(t, tt.expected, v.All())
		})
	}
}

func Test_Controller_Get(t *testing.T) {
	mockStorage := new(MockStorage)
	mockStorage.On("Get", "FOO").Return("bar")
	mockStorage.On("Get", "BAR").Return("")

	v := &Vault{storage: mockStorage}

	var tests = []struct {
		name     string
		input    string
		expected string
	}{
		{"when the key does not exist", "FOO", "bar"},
		{"when the key exists", "BAR", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, v.Get(tt.input), tt.expected)
		})
	}
}

func Test_Controller_Upsert(t *testing.T) {
	var ErrUpsert = errors.New("upsert error")

	mockStorage := new(MockStorage)
	mockStorage.On("Upsert", "FOO", "bar").Return(nil)
	mockStorage.On("Upsert", "BAR", "baz").Return(ErrUpsert)

	v := &Vault{storage: mockStorage}

	var tests = []struct {
		name     string
		input    []string
		expected error
	}{
		{"when the key does not exist", []string{"FOO", "bar"}, nil},
		{"when the key exists", []string{"BAR", "baz"}, ErrUpsert},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, v.Upsert(tt.input[0], tt.input[1]), tt.expected)
		})
	}
}

func Test_Controller_Destroy(t *testing.T) {
	var ErrDestroy = errors.New("destroy error")

	mockStorage := new(MockStorage)
	mockStorage.On("Destroy", "FOO").Return(nil)
	mockStorage.On("Destroy", "BAR").Return(ErrDestroy)

	v := &Vault{storage: mockStorage}

	var tests = []struct {
		name     string
		input    string
		expected error
	}{
		{"when the key does not exist", "FOO", nil},
		{"when the key exists", "BAR", ErrDestroy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, v.Destroy(tt.input), tt.expected)
		})
	}
}
