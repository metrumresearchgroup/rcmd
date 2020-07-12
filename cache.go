package rcmd

import (
	"fmt"
	"os"
	"path/filepath"
)

//NewPackageCache provides a PackageCache, optionally forcing that
// the package cache directory is definitely exist or be created
// the overall hierarchy of the cache is to provide the top level cache directory
func NewPackageCache(dir string, mustWork bool) PackageCache {
	if !filepath.IsAbs(dir) {
		wd, _ := os.Getwd()
		dir = filepath.Join(wd, dir)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0777)
		if mustWork && err != nil {
			panic(fmt.Sprintf("error creating cache at: %s", dir))
		}
	}
	return PackageCache{BaseDir: dir}
}
