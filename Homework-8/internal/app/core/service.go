package core

// Service is a core object, that contains main app public functions
type Service struct {
	ordersService ordersService
	pvzService    pvzService
	logger        loggerOps
}

// NewCoreService creates core Service
func NewCoreService(ordersService ordersService, pvzService pvzService, logger loggerOps) *Service {
	return &Service{
		ordersService: ordersService,
		pvzService:    pvzService,
		logger:        logger,
	}
}
