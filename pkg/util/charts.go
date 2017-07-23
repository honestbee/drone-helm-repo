package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//FindCharts feeds paths to directories which have a `Chart.yaml` at their root
func FindCharts(basePath string, exclude []string, logger *Logger) chan *FileStat {

	out := make(chan *FileStat)

	basePath = filepath.ToSlash(basePath)

	go func() {
		defer close(out)
		start := time.Now()
		logger.Debug.Printf("read local - start at %s", start)

		stat, err := os.Stat(basePath)
		if err != nil {
			logger.Err.Printf("%s\n", err)
			return
		}

		if !stat.IsDir() {
			out <- &FileStat{
				Err: fmt.Errorf("aborting, expecting directory to walk"),
			}
			return
		}

		err = filepath.Walk(basePath, func(filePath string, stat os.FileInfo, err error) error {
			relativePath := relativePath(basePath, filepath.ToSlash(filePath))
			for _, excluded := range exclude {
				if relativePath == excluded {
					logger.Debug.Printf("excluding %s\n", relativePath)
					if stat.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
			if stat == nil || stat.IsDir() {
				return nil
			}
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				out <- &FileStat{
					Err: err,
				}
			}
			//logger.Debug.Printf("I saw %s (%s)", relativePath, absPath)
			if stat.Name() == "Chart.yaml" {
				out <- &FileStat{
					Name:    relativePath,
					Dir:     filepath.Dir(filePath),
					Path:    absPath,
					ModTime: stat.ModTime(),
					Size:    stat.Size(),
				}
				//jump to next dir
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			logger.Err.Println(err)
		}

		logger.Debug.Printf("read local - end, it took %s", time.Since(start))
	}()

	return out
}

//PackageCharts takes a feed of Chart.yaml locations and returns a feed of packaged charts together with an Index.yaml
func PackageCharts(in chan *FileStat, tempDir string, logger *Logger, url string) chan *FileStat {

	helm := helmCommand{
		helmBin:     helmBin,
		logger:      logger,
		destination: tempDir,
	}
	// Make sure helm client is initialized
	helm.doInit()

	out := make(chan *FileStat)
	go func() {
		defer close(out)
		for c := range in {
			if c.Err != nil {
				logger.Err.Println(c.Err)
				continue
			}
			out <- helm.packageChart(c.Dir)
		}
		out <- helm.generateIndex(url)
	}()

	return out
}

func relativePath(path string, filePath string) string {
	if path == "." {
		return strings.TrimPrefix(filePath, "/")
	}
	path = strings.TrimPrefix(path, "./")
	a := strings.TrimPrefix(filePath, path)
	return strings.TrimPrefix(a, "/")
}
