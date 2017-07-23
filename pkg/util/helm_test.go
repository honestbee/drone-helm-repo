package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ghodss/yaml"
)

// IndexFile represents the index file in a chart repository
type IndexFile struct {
	APIVersion string                     `json:"apiVersion"`
	Generated  time.Time                  `json:"generated"`
	Entries    map[string][]*ChartVersion `json:"entries"`
	PublicKeys []string                   `json:"publicKeys,omitempty"`
}

// ChartVersion represents a chart entry in the IndexFile
type ChartVersion struct {
	URLs    []string  `json:"urls"`
	Created time.Time `json:"created,omitempty"`
	Removed bool      `json:"removed,omitempty"`
	Digest  string    `json:"digest,omitempty"`
}

func loadIndex(path string) (*IndexFile, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	i := &IndexFile{}
	if err := yaml.Unmarshal(data, i); err != nil {
		return i, err
	}
	return i, nil
}

func TestHelmIndex(t *testing.T) {
	logger, _ := getTestLogger()

	helm := helmCommand{
		helmBin:     helmBin,
		logger:      logger,
		destination: filepath.Join(testdir, "samplechart"),
	}

	f := helm.generateIndex(url)
	defer os.Remove(f.Path) //clean up

	if f.Err != nil {
		t.Error(f.Err)
	}

	i, err := loadIndex(f.Path)
	if err != nil {
		t.Error(err)
	}

	if len(i.Entries) != 1 {
		t.Errorf("Expected Entries to be %d but got %d", 1, len(i.Entries))
	}
}

func TestHelmPackage(t *testing.T) {
	logger, buf := getTestLogger()
	tempDir, err := ioutil.TempDir(testdir, "output")
	if err != nil {
		t.Errorf("%s\n", err)
	}
	defer os.RemoveAll(tempDir) // clean up
	helm := helmCommand{
		helmBin:     helmBin,
		logger:      logger,
		destination: tempDir,
	}

	//change tempdir to absolute path
	if tempDir, err = filepath.Abs(tempDir); err != nil {
		t.Errorf("%s\n", err)
	}

	tests := []struct {
		in  string
		out string
	}{
		{in: "one", out: "one-0.1.0.tgz"},
		{in: "two", out: "two-0.1.0.tgz"},
		{in: "three", out: "three-0.2.0.tgz"},
		{in: "twentyone", out: "twentyone-0.1.0.tgz"},
		{in: "onetwo", out: "onetwo-0.1.0.tgz"},
	}

	for _, test := range tests {
		buf.Reset()
		f := helm.packageChart(filepath.Join(testdir, test.in))

		if f.Err != nil {
			t.Errorf("Error: %v", f.Err)
		}
		if f.Path != filepath.Join(tempDir, test.out) {
			t.Errorf("Expected %q, got %q", filepath.Join(tempDir, test.out), f.Path)
			t.Errorf("%s\n", buf)
		}
	}
}
