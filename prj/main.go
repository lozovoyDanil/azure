package main

import (
	"filmlib/auth"
	"filmlib/endpoints"
	"filmlib/repo"
	"log"
)

func main() {
	var dbURI string = "mongodb://user:12345@mongo3:27017"

	mongoRepo, err := repo.NewMongoRepo(dbURI)
	if err != nil {
		log.Fatal(err)
	}

	svc := auth.NewService(mongoRepo)
	authHandler := endpoints.NewAuthHandler(svc)

	e := authHandler.InitRoutes()

	if err := e.Start(":5003"); err != nil {
		mongoRepo.MongoShutdown()
		log.Fatal(err)
	}
}
