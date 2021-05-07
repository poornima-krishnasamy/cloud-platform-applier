package sysutil_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloud-platform-applier/pkg/sysutil"
)

func TestListFolders(t *testing.T) {
	repoPath := "somerepo"

	os.MkdirAll(repoPath, os.ModePerm)
	path := filepath.Join(repoPath, "namespace-01")
	folders, err := sysutil.ListFolders(path)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(repoPath)

}
