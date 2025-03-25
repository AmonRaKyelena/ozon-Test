package post

import (
	"ozon-test-project/internal/pkg/storage/mocks"

	"github.com/gojuno/minimock/v3"
)

type mocksData struct {
	storageMock *mocks.StorageMock
}

func newMock(ctrl *minimock.Controller) mocksData {
	return mocksData{
		storageMock: mocks.NewStorageMock(ctrl),
	}
}
