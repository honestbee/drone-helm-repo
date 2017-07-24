package main

import (
	"bytes"
	"log"
	"testing"

	"github.com/honestbee/drone-helm-repo/pkg/util"
)

func TestValidateConfig(t *testing.T) {

	testConfigs := []struct {
		config Config
		pass   bool
	}{
		{config: Config{StorageURL: "gcs://test-bucket/prefix"}, pass: false}, //unsupported protocol
		{config: Config{StorageURL: "s3://test-bucket/prefix"}, pass: false},  //missing aws-region
		{config: Config{StorageURL: "s3://test-bucket/prefix", AWSRegion: "ap-southeast-1"}, pass: true},
	}

	for _, c := range testConfigs {
		err := validateConfig(c.config)
		if (err == nil) != c.pass {
			if c.pass {
				t.Errorf("Expected %v to pass", c.config)
			} else {
				t.Errorf("Expected %v to fail", c.config)
			}
		}
	}

}

func TestStoreFiles(t *testing.T) {
	logger, buf := getTestLogger()

	testfiles := []string{
		"one",
		"two",
		"three",
	}

	filesChan := make(chan *util.FileStat)

	go func() {
		defer close(filesChan)
		for _, f := range testfiles {
			filesChan <- &util.FileStat{
				Name: f,
				Dir:  ".",
			}
		}
	}()

	storedFileCount, err := storeFiles(createStorageMock(), filesChan, logger)

	if err != nil {
		t.Error(err)
		return
	}
	if storedFileCount != len(testfiles) {
		t.Errorf("expected %d - got %d", len(testfiles), storedFileCount)
		t.Errorf("%s\n", buf)
		return
	}
}

type mockedObjectStore struct{}

func createStorageMock() *mockedObjectStore {
	return &mockedObjectStore{}
}

func (m *mockedObjectStore) StoreFile(file *util.FileStat, logger *util.Logger) error {
	//do nothing
	return nil
}

func getTestLogger() (*util.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	return &util.Logger{
		Out:   log.New(buf, "[Out] ", log.Lshortfile),
		Err:   log.New(buf, "[Err] ", log.Lshortfile),
		Debug: log.New(buf, "[DEBUG] ", log.Lshortfile),
	}, buf
}
