package sysutil

import (
	"testing"

	"github.com/cloud-platform-applier/sysutil"
	gomock "github.com/golang/mock/gomock"
)

func TestListFolders(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockFileSystem := sysutil.NewMockFileSystemInterface(mockCtrl)

}
