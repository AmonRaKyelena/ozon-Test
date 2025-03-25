package post

import (
	"github.com/AmonRaKyelena/ozon-Test/internal/pkg/storage/mocks"

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
