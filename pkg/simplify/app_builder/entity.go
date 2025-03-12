package app_builder

type App interface {
	Run() error
}

type Builder interface {
	LoadConfig() Builder
	InitRepositories() Builder
	InitUseCases() Builder
	InitHandlers() Builder
	InitRoutes() Builder
	Build() App
}
