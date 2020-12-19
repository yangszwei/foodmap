package registry

import (
	"foodmap/internal/infra/config"
	"foodmap/internal/infra/delivery"
	"foodmap/internal/infra/entity/store"
	"foodmap/internal/infra/persistence"
	"foodmap/internal/infra/validator"
	"log"
)

var v = validator.New()

// Start start app
func Start(configPath string) error {
	if configPath == "" {
		configPath = ".env"
	}
	cfg, err := config.Open(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	if err := config.Validate(cfg, *v); err != nil {
		log.Fatal(err.Error())
	}
	db, err := persistence.Open(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}
	r, err := setupRepos(&db)
	if err != nil {
		log.Fatalln(err)
	}
	u := setupUsecase(r)
	server := delivery.NewServer(cfg.Server.DataPath)
	delivery.SetupSession(server, []byte(cfg.Server.SessionSecret))
	server.SetupRouterGroups()
	store.SetupAPI(server, u.store)
	return server.Run(cfg.Server.Addr)
}
