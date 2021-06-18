package store

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jaimem88/zearch/internal/model"
)

var ErrNotFound = errors.New("not found")

type Storage struct {
	UsersMap         map[model.UserID]model.User
	TicketsMap       map[model.TicketID]model.Ticket
	OrganizationsMap map[model.OrgID]model.Organization

	OrgsUsers   map[model.OrgID][]model.UserID
	OrgsTickets map[model.OrgID][]model.TicketID

	// for reverse lookup
	UsersOrgs   map[model.UserID]model.OrgID
	TicketsOrgs map[model.TicketID]model.OrgID
}

// New creates an instance of Storage and preprocess the data to store it in its
// corresponding data structures.
func New(organizations model.Organizations, users model.Users, tickets model.Tickets) *Storage {
	orgsMap := map[model.OrgID]model.Organization{}
	usersMap := map[model.UserID]model.User{}
	ticketsMap := map[model.TicketID]model.Ticket{}

	orgsUsers := map[model.OrgID][]model.UserID{}
	orgsTickets := map[model.OrgID][]model.TicketID{}

	usersOrgs := map[model.UserID]model.OrgID{}
	ticketsOrgs := map[model.TicketID]model.OrgID{}

	//wg := sync.WaitGroup{}
	//wg.Add(3)

	//go func() {
	//	defer wg.Done()

	for _, org := range organizations {
		// assumes there are no duplicate IDs, otherwise the data would be overridden
		// unsafe to do type assertions without checking if it succeeded, but assuming it's correct for simplicity
		orgID := model.OrgID(org["_id"].(float64))
		orgsMap[orgID] = org
	}
	//}()
	//
	//go func() {
	//	defer wg.Done()
	fmt.Printf("Print the first user? k: %+v", users[0])
	for k, user := range users {
		fmt.Printf("Print a user? k: %+v v:%+v\n", k, user)
		// assumes there are no duplicate IDs, otherwise data would be overridden
		userID := model.UserID(user["_id"].(float64))
		usersMap[userID] = user
		//orgID := model.OrgID(user["organization_id"].(float64))
		//orgsUsers[orgID] = append(orgsUsers[orgID], userID)
		//usersOrgs[userID] = orgID
	}
	//}()

	//go func() {
	//	defer wg.Done()

	for _, ticket := range tickets {
		fmt.Printf("a ticket looks like: %+v\n", ticket)
		// assumes there are no duplicate IDs, otherwise data would be overridden
		ticketID := model.TicketID(ticket["_id"].(string))
		ticketsMap[ticketID] = ticket

		//orgID := model.OrgID(ticket["organization_id"].(float64))
		//orgsTickets[orgID] = append(orgsTickets[orgID], ticketID)
		//ticketsOrgs[ticketID] = orgID
	}
	//}()

	//wg.Wait()

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

func (s *Storage) Organizations(term, value string) ([]model.OrganizationResult, error) {
	fmt.Printf("Searching organizations by: %q with value: %q\n", term, value)

	if term == "_id" {
		return s.searchOrgByID(value)

	}
	return s.searchOrgByTerm(term, value)

}

func (s *Storage) Tickets(term, value string) []model.TicketResult {
	fmt.Printf("searching tickets for: %q:%q\n", term, value)
	return []model.TicketResult{}
}

func (s *Storage) Users(term, value string) []model.UserResult {
	fmt.Printf("searching users for: %q:%q\n", term, value)
	return []model.UserResult{}
}

func (s *Storage) searchOrgByID(value string) ([]model.OrganizationResult, error) {
	var results []model.OrganizationResult

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

		userNames = append(userNames, user["name"].(string))
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

		ticketSubjects = append(ticketSubjects, ticket["subject"].(string))
	}

	orgResult := model.OrganizationResult{
		Organization:   org,
		UserNames:      userNames,
		TicketSubjects: ticketSubjects,
	}

	results = append(results, orgResult)

	return results, nil
}

func (s *Storage) searchOrgByTerm(term, value string) ([]model.OrganizationResult, error) {
	var result []model.OrganizationResult

	return result, nil
}
