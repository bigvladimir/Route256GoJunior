package server

import (
	"testing"

	"github.com/golang/mock/gomock"

	mock_server "homework/internal/app/server/mocks"
)

type coreServiceFixtures struct {
	ctrl        *gomock.Controller
	srv         Server
	mockCoreOps *mock_server.MockcoreOps
}

func setUp(t *testing.T) coreServiceFixtures {
	ctrl := gomock.NewController(t)
	mockCoreOps := mock_server.NewMockcoreOps(ctrl)
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
