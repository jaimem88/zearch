package store

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/jaimem88/zearch/internal/model"
)

var ErrNotFound = errors.New("not found")

type Storage struct {
	UsersMap         map[model.UserID]*model.User
	TicketsMap       map[model.TicketID]*model.Ticket
	OrganizationsMap map[model.OrgID]*model.Organization

	OrgsUsers   map[model.OrgID][]model.UserID
	OrgsTickets map[model.OrgID][]model.TicketID

	// for reverse lookup
	UsersOrgs   map[model.UserID]model.OrgID
	TicketsOrgs map[model.TicketID]model.OrgID
}

// New creates an instance of Storage and preprocess the data to store it in its
// corresponding data structures.
func New(userData []*model.User, ticketData []*model.Ticket, orgData []*model.Organization) *Storage {
	users := map[model.UserID]*model.User{}
	tickets := map[model.TicketID]*model.Ticket{}
	orgs := map[model.OrgID]*model.Organization{}

	orgsUsers := map[model.OrgID][]model.UserID{}
	orgsTickets := map[model.OrgID][]model.TicketID{}

	usersOrgs := map[model.UserID]model.OrgID{}
	ticketsOrgs := map[model.TicketID]model.OrgID{}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		for _, org := range orgData {
			// assumes there are no duplicate IDs, otherwise data would be overridden
			orgs[org.ID] = org
		}
	}()

	go func() {
		defer wg.Done()

		for _, user := range userData {
			// assumes there are no duplicate IDs, otherwise data would be overridden
			users[user.ID] = user
			orgsUsers[user.OrganizationID] = append(orgsUsers[user.OrganizationID], user.ID)
			usersOrgs[user.ID] = user.OrganizationID
		}
	}()

	go func() {
		defer wg.Done()

		for _, ticket := range ticketData {
			// assumes there are no duplicate IDs, otherwise data would be overridden
			tickets[ticket.ID] = ticket
			orgsTickets[ticket.OrganizationID] = append(orgsTickets[ticket.OrganizationID], ticket.ID)
			ticketsOrgs[ticket.ID] = ticket.OrganizationID
		}
	}()

	wg.Wait()

	return &Storage{
		UsersMap:         users,
		TicketsMap:       tickets,
		OrganizationsMap: orgs,
		OrgsUsers:        orgsUsers,
		OrgsTickets:      orgsTickets,
		UsersOrgs:        usersOrgs,
		TicketsOrgs:      ticketsOrgs,
	}
}

func (s *Storage) Organizations(term, value string) ([]*model.OrganizationResult, error) {
	fmt.Printf("Searching organizations by: %q with value: %q\n", term, value)

	var results []*model.OrganizationResult
	var err error

	if term == "_id" {
		results, err = s.searchOrgByID(value)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (s *Storage) Tickets(term, value string) []*model.TicketResult {
	fmt.Printf("searching tickets for: %q:%q\n", term, value)
	return []*model.TicketResult{}
}

func (s *Storage) Users(term, value string) []*model.UserResult {
	fmt.Printf("searching users for: %q:%q\n", term, value)
	return []*model.UserResult{}
}

func (s *Storage) searchOrgByID(value string) ([]*model.OrganizationResult, error) {
	var results []*model.OrganizationResult

	id, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}

	orgID := model.OrgID(id)

	org, ok := s.OrganizationsMap[orgID]
	if !ok {
		return nil, ErrNotFound
	}

	usersForOrg := s.OrgsUsers[orgID]
	// initializing a slice with capacity allows us to use `append` preventing it
	// from allocating a new slice  when the capacity is reached see https://golang.org/pkg/builtin/#append
	userNames := make([]string, 0, len(usersForOrg))
	for _, userID := range usersForOrg {
		user, ok := s.UsersMap[userID]
		if !ok {
			// skip if we can't find the userID for some reason
			continue
		}

		userNames = append(userNames, user.Name)
	}

	ticketsForOrg := s.OrgsTickets[orgID]
	ticketSubjects := make([]string, 0, len(ticketsForOrg))
	// if we had generics, maybe this could have been implemented once
	for _, ticketID := range ticketsForOrg {
		ticket, ok := s.TicketsMap[ticketID]
		if !ok {
			// skip if we can't find the ticketID for some reason
			continue
		}

		ticketSubjects = append(ticketSubjects, ticket.Subject)
	}

	orgResult := &model.OrganizationResult{
		Organization:   org,
		UserNames:      userNames,
		TicketSubjects: ticketSubjects,
	}

	results = append(results, orgResult)

	return results, nil
}
