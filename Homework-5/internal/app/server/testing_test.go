package server

import (
	"testing"

	mock_core "homework/internal/app/core/mocks"

	"github.com/golang/mock/gomock"
)

type coreServiceFixtures struct {
	ctrl        *gomock.Controller
	srv         Server
	mockCoreOps *mock_core.MockCoreOps
}

func setUp(t *testing.T) coreServiceFixtures {
	ctrl := gomock.NewController(t)
	mockCoreOps := mock_core.NewMockCoreOps(ctrl)
	srv := Server{mockCoreOps}
	return coreServiceFixtures{
		ctrl:        ctrl,
		mockCoreOps: mockCoreOps,
		srv:         srv,
	}
}

func (a *coreServiceFixtures) tearDown() {
	a.ctrl.Finish()
}
