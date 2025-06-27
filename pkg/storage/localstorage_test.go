package storage

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"passline/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	// Return os.FileInfo (can be nil) and error
	fileInfo, _ := args.Get(0).(os.FileInfo)
	return fileInfo, args.Error(1)
}

func (m *MockFileSystem) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	data, _ := args.Get(0).([]byte)
	return data, args.Error(1)
}

func (m *MockFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	args := m.Called(name, data, perm)
	return args.Error(0)
}

func (m *MockFileSystem) IsNotExist(err error) bool {
	args := m.Called(err)
	return args.Bool(0)
}

func TestGetAllItems(t *testing.T) {
	mockFS := new(MockFileSystem)

	mockFS.On("Stat", "/root/.passline/storage").Return(nil, nil)

	mockedItems := []Item{{Name: "test", Credentials: []Credential{{Username: "tpoe", Password: "test"}}}}
	mockedData, _ := json.Marshal(mockedItems)
	mockFS.On("ReadFile", "/root/.passline/storage/storage").Return(mockedData, nil)

	s, err := NewLocalStorage(mockFS)
	if err != nil {
		t.Errorf("Unable to initialize storage")
	}

	items, err := s.GetAllItems(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, items)
	assert.Equal(t, 1, len(items))
}

func TestAddItem(t *testing.T) {
	rootDir := config.Directory()

	mockFS := new(MockFileSystem)

	mockFS.On("Stat", rootDir+"/storage").Return(nil, os.ErrNotExist)

	mockFS.On("ReadFile", rootDir+"/config.json").Return(nil, nil)

	mockedItems := []Item{{Name: "test", Credentials: []Credential{{Username: "tpoe", Password: "test"}}}}
	mockedData, _ := json.Marshal(mockedItems)
	mockFS.On("ReadFile", rootDir+"/storage/storage").Return(mockedData, nil)

	mockFS.On("WriteFile", rootDir+"/storage/storage", mock.Anything, mock.Anything).Return(nil)

	cfg, err := config.Get(mockFS)
	if err != nil {
		t.Errorf("Unable to initialize config")
	}
	ctx := cfg.WithContext(context.Background())

	s, err := NewLocalStorage(mockFS)
	if err != nil {
		t.Errorf("Unable to initialize storage")
	}

	credential := Credential{
		Username: "tpoe2",
		Password: "1234",
	}

	s.AddCredential(ctx, "test", credential)

	mockFS.AssertCalled(t, "WriteFile", rootDir+"/storage/storage", mock.Anything, mock.Anything)
}

func TestDeleteWholeItem(t *testing.T) {
	rootDir := config.Directory()

	mockFS := new(MockFileSystem)

	mockFS.On("Stat", rootDir+"/storage").Return(nil, os.ErrNotExist)

	mockFS.On("ReadFile", rootDir+"/config.json").Return(nil, nil)

	mockedItems := []Item{{Name: "test", Credentials: []Credential{{Username: "tpoe", Password: "test"}}}}
	mockedData, _ := json.Marshal(mockedItems)
	mockFS.On("ReadFile", rootDir+"/storage/storage").Return(mockedData, nil)

	mockFS.On("WriteFile", rootDir+"/storage/storage", mock.Anything, mock.Anything).Return(nil)

	cfg, err := config.Get(mockFS)
	if err != nil {
		t.Errorf("Unable to initialize config")
	}
	ctx := cfg.WithContext(context.Background())

	s, err := NewLocalStorage(mockFS)
	if err != nil {
		t.Errorf("Unable to initialize storage")
	}

	s.DeleteCredential(ctx, Item{Name: "test"}, "tpoe")

	expectedItems := []Item{}
	expectedData, _ := json.Marshal(expectedItems)
	mockFS.AssertCalled(t, "WriteFile", rootDir+"/storage/storage", expectedData, mock.Anything)
}

func TestDeleteOneCredential(t *testing.T) {
	rootDir := config.Directory()

	mockFS := new(MockFileSystem)

	mockFS.On("Stat", rootDir+"/storage").Return(nil, os.ErrNotExist)

	mockFS.On("ReadFile", rootDir+"/config.json").Return(nil, nil)

	mockedItems := []Item{{
		Name: "test",
		Credentials: []Credential{{
			Username: "tpoe1",
			Password: "password1",
		}, {
			Username: "tpoe2",
			Password: "password2",
		}},
	}}
	mockedData, _ := json.Marshal(mockedItems)
	mockFS.On("ReadFile", rootDir+"/storage/storage").Return(mockedData, nil)

	mockFS.On("WriteFile", rootDir+"/storage/storage", mock.Anything, mock.Anything).Return(nil)

	cfg, err := config.Get(mockFS)
	if err != nil {
		t.Errorf("Unable to initialize config")
	}
	ctx := cfg.WithContext(context.Background())

	s, err := NewLocalStorage(mockFS)
	if err != nil {
		t.Errorf("Unable to initialize storage")
	}

	s.DeleteCredential(ctx, Item{Name: "test"}, "tpoe2")

	expectedItems := []Item{{
		Name: "test",
		Credentials: []Credential{{
			Category: "default",
			Username: "tpoe1",
			Password: "password1",
		}},
	}}
	expectedData, _ := json.Marshal(expectedItems)
	mockFS.AssertCalled(t,
		"WriteFile",
		rootDir+"/storage/storage",
		mock.MatchedBy(func(data []byte) bool {
			return assert.ObjectsAreEqual(expectedData, data)
		}),
		mock.Anything)
}
