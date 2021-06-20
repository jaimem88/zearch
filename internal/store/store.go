package store

import (
	"errors"
	"sort"
	"sync"

	"github.com/jaimem88/zearch/internal/model"
)

// ErrNotFound returned when any search cannot find any match
var ErrNotFound = errors.New("not found")

// Storage holds an in-memory set of maps that will be used to store and lookup
// values per key
type Storage struct {
	// store the users, tickets and organizations by ID so that accessing them
	// is done in constant time
	usersMap         map[model.UserID]model.User
	ticketsMap       map[model.TicketID]model.Ticket
	organizationsMap map[model.OrgID]model.Organization

	// Keep a list of users and tickets per orgID
	orgsUsers   map[model.OrgID][]model.UserID
	orgsTickets map[model.OrgID][]model.TicketID

	searchableFields map[string][]string
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

	searchableFields := map[string][]string{}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		for k, org := range organizations {
			// assumes there are no duplicate IDs, otherwise the data would be overridden
			// unsafe to do type assertions without checking if it succeeded, but assuming it's correct for simplicity
			orgID := model.OrgID(org["_id"].(float64))
			orgsMap[orgID] = org

			// Get the searchable fields from the first element programmatically. The caveat to this approach is that
			// if other objects have more fields they won't be printed as searchable.
			if k == 0 {
				searchableFields["organizations"] = getOrgFields(org)
			}
		}
	}()

	go func() {
		defer wg.Done()

		for k, user := range users {
			userID := model.UserID(user["_id"].(float64))
			usersMap[userID] = user

			orgID, ok := user["organization_id"].(float64)
			if ok {
				orgID := model.OrgID(orgID)
				orgsUsers[orgID] = append(orgsUsers[orgID], userID)
			}

			if k == 0 {
				searchableFields["users"] = getUserFields(user)
			}
		}
	}()

	go func() {
		defer wg.Done()

		for k, ticket := range tickets {
			ticketID := model.TicketID(ticket["_id"].(string))
			ticketsMap[ticketID] = ticket

			orgID, ok := ticket["organization_id"].(float64)
			if ok {
				orgID := model.OrgID(orgID)
				orgsTickets[orgID] = append(orgsTickets[orgID], ticketID)
			}

			if k == 0 {
				searchableFields["tickets"] = getTicketFields(ticket)
			}
		}
	}()

	wg.Wait()

	return &Storage{
		usersMap:         usersMap,
		ticketsMap:       ticketsMap,
		organizationsMap: orgsMap,
		orgsUsers:        orgsUsers,
		orgsTickets:      orgsTickets,
		searchableFields: searchableFields,
	}
}

func getOrgFields(org model.Organization) []string {
	fields := make([]string, 0, len(org))
	for k := range org {
		fields = append(fields, k)
	}

	sort.Strings(fields)
	return fields
}

func getUserFields(user model.User) []string {
	fields := make([]string, 0, len(user))
	for k := range user {
		fields = append(fields, k)
	}

	sort.Strings(fields)
	return fields
}

func getTicketFields(ticket model.Ticket) []string {
	fields := make([]string, 0, len(ticket))
	for k := range ticket {
		fields = append(fields, k)
	}

	sort.Strings(fields)
	return fields
}

// GetSearchableFields returns the list of fields per entity contained in the store
func (s *Storage) GetSearchableFields() map[string][]string {
	return s.searchableFields
}
