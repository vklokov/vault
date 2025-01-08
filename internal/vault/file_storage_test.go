package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tmp = filepath.Join(projectRoot(), "test_storage.enc")
var secret = "ffffffffffffffffffffffffffffffff"

func setupFileStorage(t *testing.T, records Records) *FileStorage {
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	defer f.Close()

	data, err := json.Marshal(records)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	body, err := encryptString(string(data), secret)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	f.Write([]byte(body))

	s := FileStorage{
		filename: tmp,
		secret:   secret,
	}

	err = s.Init()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	return &s
}

func Test_FileStorage_Init(t *testing.T) {
	t.Run("when the file does not exist", func(t *testing.T) {
		defer os.Remove(tmp)

		s := FileStorage{
			secret:      "dummy",
			filename:    tmp,
			initialized: false,
		}

		err := s.Init()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		assert.True(t, s.initialized)
		assert.NotNil(t, s.records)

		_, err = os.Stat(tmp)
		assert.NoError(t, err)
	})

	t.Run("when the file exists", func(t *testing.T) {
		s := setupFileStorage(t, Records{"FOO": "bar"})
		defer os.Remove(tmp)

		assert.True(t, s.initialized)
		assert.NotNil(t, s.records)
		assert.Equal(t, s.records, Records{"FOO": "bar"})

		_, err := os.Stat(tmp)
		assert.NoError(t, err)
	})
}

func Test_FileStorage_Exists(t *testing.T) {
	s := setupFileStorage(t, Records{"FOO": "bar"})
	defer os.Remove(tmp)

	var tests = []struct {
		name     string
		input    string
		expected bool
	}{
		{"when the key exists", "FOO", true},
		{"when the key does not exist", "BAR", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, s.Exists(tt.input), tt.expected)
		})
	}
}

func Test_FileStorage_All(t *testing.T) {
	s := setupFileStorage(t, Records{"FOO": "bar"})
	defer os.Remove(tmp)

	assert.Equal(t, s.All(), Records{"FOO": "bar"})
}

func Test_FileStorage_Get(t *testing.T) {
	s := setupFileStorage(t, Records{"FOO": "bar"})
	defer os.Remove(tmp)

	var tests = []struct {
		name     string
		input    string
		expected string
	}{
		{"when the key exist", "FOO", "bar"},
		{"when the key does not exist", "BAR", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, s.Get(tt.input), tt.expected)
		})
	}
}

func Test_FileStorage_Upsert(t *testing.T) {
	s := setupFileStorage(t, Records{})
	defer os.Remove(tmp)

	s.Upsert("BAZ", "baz")
	assert.Equal(t, s.records, Records{"BAZ": "baz"})
}

func Test_FileStorage_Destroy(t *testing.T) {
	s := setupFileStorage(t, Records{"FOO": "bar"})
	defer os.Remove(tmp)

	s.Destroy("FOO")
	assert.Equal(t, s.records, Records{})
}

func Test_FileStorage_rewrite(t *testing.T) {
	var filesize = func(filename string) int64 {
		fi, err := os.Stat(filename)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		return fi.Size()
	}

	s := setupFileStorage(t, Records{"FOO": "bar"})
	defer os.Remove(tmp)

	sizeBefore := filesize(s.filename)

	s.records = Records{}
	s.rewrite()

	sizeAfter := filesize(s.filename)

	assert.NotEqual(t, sizeBefore, sizeAfter)
}
