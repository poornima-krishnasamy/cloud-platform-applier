package sysutil_test

import (
	"log"
	"os"
	"testing"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/sysutil"
)

func TestListFolders(t *testing.T) {
	repoPath := "somerepo"

	os.MkdirAll(repoPath, os.ModePerm)

	ns_folder := "somerepo/somenamespace"

	os.Mkdir(ns_folder, os.ModePerm)

	folders, err := sysutil.ListFolderPaths(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, folder := range folders {
		t.Logf("Found directory %v\n", folder)
	}

	_, found := Find(folders, ns_folder)
	if !found {
		t.Errorf("Expected directory %v not found", ns_folder)
	}
	defer os.RemoveAll(repoPath)

}
