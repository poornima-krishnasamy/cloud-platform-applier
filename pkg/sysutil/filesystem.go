// Package sysutil provides utility functions needed for the cloud-platform-applier
package sysutil

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// FileSystemInterface allows for mocking out the functionality of FileSystem to avoid calls to the actual file system during testing.
type FileSystemInterface interface {
	ListFolderPaths(path string) ([]string, error)
}

// FileSystem provides utility functions for interacting with the file system.
type FileSystem struct{}

// ListFolders take the path as input, list all the folders in the give path and
// return a array of strings containing the list of folders
func ListFolderPaths(path string) ([]string, error) {
	var folders []string

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
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
			// log.Printf("applying folder: %q\n", folders)
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return folders, nil
}

func (fs *FileSystem) ChunkFolders(folders []string, nRoutines int) ([][]string, error) {

	nChunks := len(folders) / nRoutines

	fmt.Println("Number of folders per chunk", nChunks)

	var folderChunks [][]string
	for {
		if len(folders) == 0 {
			break
		}

		if len(folders) < nChunks {
			nChunks = len(folders)
		}

		folderChunks = append(folderChunks, folders[0:nChunks])
		folders = folders[nChunks:]
	}
	return folderChunks, nil
}

func listFiles(path string) ([]string, error) {

	var files []string

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatalf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}
			if info.IsDir() {
				fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
				return filepath.SkipDir
			}
			files = append(files, path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return files, nil

}
