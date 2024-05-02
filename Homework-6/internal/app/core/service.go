package core

type Service struct {
	ordersService ordersService
	pvzService    pvzService
	logger        loggerOps
}

func NewCoreService(ordersService ordersService, pvzService pvzService, logger loggerOps) *Service {
	return &Service{
		ordersService: ordersService,
		pvzService:    pvzService,
		logger:        logger,
	}
}
