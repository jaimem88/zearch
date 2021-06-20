package model

import (
	"fmt"

	"github.com/jaimem88/zearch/internal/reader"
)

// Data holds all the parsed data per entity.
type Data struct {
	Organizations Organizations
	Users         Users
	Tickets       Tickets
}

// LoadData will read all filenames content and return the parsed data into their
// respective types.
func LoadData(orgsFilename, usersFilename, ticketsFilename string) (*Data, error) {
	var orgs Organizations
	var users Users
	var tickets Tickets

	err := reader.ReadJSONFile(orgsFilename, &orgs)
	if err != nil {
		return nil, fmt.Errorf("failed to load: %s %w", orgsFilename, err)
	}

	err = reader.ReadJSONFile(usersFilename, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to load: %s %w", usersFilename, err)
	}

	err = reader.ReadJSONFile(ticketsFilename, &tickets)
	if err != nil {
		return nil, fmt.Errorf("failed to load: %s %w", ticketsFilename, err)
	}

	return &Data{
		Organizations: orgs,
		Users:         users,
		Tickets:       tickets,
	}, nil
}
