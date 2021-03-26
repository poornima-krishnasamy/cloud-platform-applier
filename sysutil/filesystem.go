// +build admin

package sysutil

import (
	"io/fs"
	"log"
	"path/filepath"
)

type FileSystem struct{}

func (f *FileSystem) listFolders(path string) ([]string, error) {
	var folders []string
	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.Name() == ".terraform" || info.Name() == "resources" {
			return filepath.SkipDir
		}

		if info.IsDir() {
			folders = append(folders, path)
		}
		log.Printf("applying folder: %q\n", path)
		return nil
	})

	// filepath.Walk {
	// 	return func(path string, info os.FileInfo, err error) error {
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		if info.Name() == ".terraform" || info.Name() == "resources" {
	// 			return filepath.SkipDir
	// 		}

	// 		if info.IsDir() {
	// 			*folders = append(*folders, path)
	// 		}
	// 		return nil
	// 	}
	// }
	return folders, nil
}
