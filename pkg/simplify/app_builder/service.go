package app_builder

func Apply(builder Builder) App {
	return builder.
		LoadConfig().
		InitRepositories().
		InitUseCases().
		InitHandlers().
		InitRoutes().
		Build()
}
