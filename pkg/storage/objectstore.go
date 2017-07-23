package storage

import (
	"github.com/honestbee/drone-helm-repo/pkg/util"
)

// ObjectStore defines methods for storing objects
type (
	ObjectStore interface {
		// Store stores the file specified
		StoreFile(file *util.FileStat, logger *util.Logger) error
	}
)
