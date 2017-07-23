package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

// shared test code
var url = "http://test-repo.example.com"

//testdata from main project directory...
const testdir = "../../_testdata"

func sink(in chan *FileStat) map[string]*FileStat {
	out := make(map[string]*FileStat)
	for f := range in {
		if f.Err != nil {
			continue
		}
		out[f.Name] = f
	}
	return out
}

func getTestLogger() (*Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	return &Logger{
		Out:   log.New(buf, "[Out] ", log.Lshortfile),
		Err:   log.New(buf, "[Err] ", log.Lshortfile),
		Debug: log.New(buf, "[DEBUG] ", log.Lshortfile),
	}, buf
}

// Tests on charts utilities

func TestFindCharts(t *testing.T) {
	logger, buf := getTestLogger()

	var exclude []string
	fileChan := FindCharts(testdir, exclude, logger)

	files := sink(fileChan)

	expected := 6
	actual := len(files)
	if actual != expected {
		t.Errorf("wanted %d files, got %d files", expected, actual)
		t.Errorf("%s\n", buf)
	}

	filename := "one/Chart.yaml"
	chartname := filepath.Join(testdir, "one")
	file, ok := files[filename]
	if !ok {
		t.Errorf("Couldn't find file '%s' in file list", filename)
		t.Errorf("%+v", files)
		return
	}

	if file.Err != nil {
		t.Errorf("Expected file.Err to be nil, got %v\n", file.Err)
	}
	if file.Dir != chartname {
		t.Errorf("Expected file.Dir to be %s, got %s\n", chartname, file.Dir)
	}
	if file.Size != 81 {
		t.Errorf("Expected file.Name to be %d, got %d\n", 81, file.Size)
	}

	if file.Name != filename {
		t.Errorf("expected file.Name ('%s') to be the same as the key ('%s') of the map", file.Name, filename)
	}
}

func TestFindChartsWithExclude(t *testing.T) {

	tests := []struct {
		in      string
		out     int
		exclude []string
	}{
		{in: testdir, out: 5, exclude: []string{".excludeme"}},
	}

	for _, test := range tests {
		logger, buf := getTestLogger()
		chartChan := FindCharts(test.in, test.exclude, logger)
		files := sink(chartChan)
		if len(files) != test.out {
			t.Errorf("wanted %d files, got %d files", test.out, len(files))
			t.Errorf("%s\n", buf)
		}
	}
}

func TestHelmPackageSink(t *testing.T) {
	logger, buf := getTestLogger()
	tempDir, err := ioutil.TempDir(testdir, "output")
	if err != nil {
		t.Errorf("%s\n", err)
	}
	if tempDir, err = filepath.Abs(tempDir); err != nil {
		t.Errorf("%s\n", err)
	}
	defer os.RemoveAll(tempDir) // clean up

	//expect following packages from _testdata dir
	tests := []struct {
		chart   string
		version string
	}{
		{chart: "one", version: "0.1.0"},
		{chart: "two", version: "0.1.0"},
		{chart: "three", version: "0.2.0"},
		{chart: "twentyone", version: "0.1.0"},
		{chart: "onetwo", version: "0.1.0"},
	}

	//this test does not exclude the invalid chart ".excludeme"
	chartChan := FindCharts(testdir, []string{}, logger)
	packageChan := PackageCharts(chartChan, tempDir, logger, url)
	packages := sink(packageChan)

	indexStats, ok := packages["index.yaml"]
	if !ok {
		t.Error("Couldn't find 'index.yaml' in package list")
		t.Errorf("%s\n", buf)
	}

	index, err := loadIndex(indexStats.Path)
	if err != nil {
		t.Errorf("Error loading index: %v", err)
	}

	if len(index.Entries) != len(tests) {
		t.Errorf("Expected index.yaml to have %d Entries, but got %d", len(tests), len(index.Entries))
	}

	for _, test := range tests {
		packageName := fmt.Sprintf("%s-%s.tgz", test.chart, test.version)
		f, ok := packages[packageName]
		if !ok {
			t.Errorf("Couldn't find '%s' in package list", packageName)
			t.Errorf("%s\n", buf)
			for _, packaged := range packages {
				t.Errorf("%+v", packaged.Name)
				if packaged.Err != nil {
					t.Errorf("%+v", packaged.Err)
				}
			}
		}

		chartVersions, ok := index.Entries[test.chart]
		if !ok {
			t.Errorf("Couldn't find '%s' in index.yaml", packageName)
			t.Errorf("%s\n", buf)
			t.Errorf("%v\n", index)
		}

		//currently this will always have a single version as we don't support index merging yet
		chartURL := fmt.Sprintf("%s/%s", url, packageName)
		if chartVersions[0].URLs[0] != chartURL {
			t.Errorf("Expected %q to match %q\n", chartVersions[0].URLs[0], chartURL)
			t.Errorf("%s\n", buf)
			return
		}

		if f.Err != nil {
			t.Errorf("Expected package.Err to be nil, got %v\n", f.Err)
			t.Errorf("%s\n", buf)
			return
		}
	}
}
