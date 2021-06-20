package main

import (
	"flag"
	"log"

	"github.com/jaimem88/zearch/internal/model"

	"github.com/jaimem88/zearch/internal/app"
	"github.com/jaimem88/zearch/internal/store"
)

var (
	usersFilename   = flag.String("users", "data/users.json", "Filename to load users from e.g. --users data/users.json")
	ticketsFilename = flag.String("tickets", "data/tickets.json", "Filename to load users from e.g. --users data/users.json")
	orgsFilename    = flag.String("organizations", "data/organizations.json", "Filename to load users from e.g. --users data/users.json")
)

func main() {
	flag.Parse()

	data, err := model.LoadData(*orgsFilename, *usersFilename, *ticketsFilename)
	if err != nil {
		log.Fatalf("load data: %+v\n", err)
	}

	c := app.New(store.New(data.Organizations, data.Users, data.Tickets))

	if err := c.Run(); err != nil {
		log.Fatalf("run: %+v\n", err)
	}
}
