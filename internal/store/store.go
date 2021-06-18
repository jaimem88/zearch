package store

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jaimem88/zearch/internal/model"
)

var ErrNotFound = errors.New("not found")

type Storage struct {
	// store the users, tickets and organizations by ID so that accessing them
	// is done in constant time
	UsersMap         map[model.UserID]model.User
	TicketsMap       map[model.TicketID]model.Ticket
	OrganizationsMap map[model.OrgID]model.Organization

	// Keep a list of users and tickets per orgID
	OrgsUsers   map[model.OrgID][]model.UserID
	OrgsTickets map[model.OrgID][]model.TicketID

	// for reverse lookup
	UsersOrgs   map[model.UserID]model.OrgID
	TicketsOrgs map[model.TicketID]model.OrgID
}

// New creates an instance of Storage and preprocess the data to store it in its
// corresponding data structures.
// The initialization process for every entity will be done on startup. Each entity is
// loaded in its own goroutine, using a sync.WaitGroup to wait for all of them to finish.
func New(organizations model.Organizations, users model.Users, tickets model.Tickets) *Storage {
	orgsMap := map[model.OrgID]model.Organization{}
	usersMap := map[model.UserID]model.User{}
	ticketsMap := map[model.TicketID]model.Ticket{}

	orgsUsers := map[model.OrgID][]model.UserID{}
	orgsTickets := map[model.OrgID][]model.TicketID{}

	usersOrgs := map[model.UserID]model.OrgID{}
	ticketsOrgs := map[model.TicketID]model.OrgID{}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		for _, org := range organizations {
			// assumes there are no duplicate IDs, otherwise the data would be overridden
			// unsafe to do type assertions without checking if it succeeded, but assuming it's correct for simplicity
			orgID := model.OrgID(org["_id"].(float64))
			orgsMap[orgID] = org
		}
	}()

	go func() {
		defer wg.Done()
		for _, user := range users {
			// assumes there are no duplicate IDs, otherwise data would be overridden
			userID := model.UserID(user["_id"].(float64))
			usersMap[userID] = user

			orgID, ok := user["organization_id"].(float64)
			if ok {
				orgID := model.OrgID(orgID)
				orgsUsers[orgID] = append(orgsUsers[orgID], userID)
				usersOrgs[userID] = orgID
			}
		}
	}()

	go func() {
		defer wg.Done()

		for _, ticket := range tickets {
			// assumes there are no duplicate IDs, otherwise data would be overridden
			ticketID := model.TicketID(ticket["_id"].(string))
			ticketsMap[ticketID] = ticket

			orgID, ok := ticket["organization_id"].(float64)
			if ok {
				orgID := model.OrgID(orgID)
				orgsTickets[orgID] = append(orgsTickets[orgID], ticketID)
				ticketsOrgs[ticketID] = orgID
			}
		}
	}()

	wg.Wait()

	return &Storage{
		UsersMap:         usersMap,
		TicketsMap:       ticketsMap,
		OrganizationsMap: orgsMap,
		OrgsUsers:        orgsUsers,
		OrgsTickets:      orgsTickets,
		UsersOrgs:        usersOrgs,
		TicketsOrgs:      ticketsOrgs,
	}
}

// Users ...
func (s *Storage) Users(term, value string) []model.UserResult {
	fmt.Printf("searching users for: %q:%q\n", term, value)
	return []model.UserResult{}
}

// Tickets ...
func (s *Storage) Tickets(term, value string) []model.TicketResult {
	fmt.Printf("searching tickets for: %q:%q\n", term, value)
	return []model.TicketResult{}
}
