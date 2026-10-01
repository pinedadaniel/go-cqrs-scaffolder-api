package read

type Handlers struct {
	Health Health
}

func Register(router Router, handlers Handlers) {
	group := router.Group("/v1")
	{
		group.GET("/health", handlers.Health.Get)
	}
}
