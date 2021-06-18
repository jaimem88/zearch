package store

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

// Organizations implements the searcher method for the app. It searches by term and value.
// Handles a special case for _id which can be looked up in the Storage easily from the
// map.
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

	orgResult := model.OrganizationResult{
		Organization:   org,
		UserNames:      s.getUsersForOrg(orgID),
		TicketSubjects: s.getTicketsForOrg(orgID),
	}

	results = append(results, orgResult)

	return results, nil
}

// searchOrgByTerm will iterate over each element of the OrganizationsMap and accessing
// the term directly. Once found, the organization and its ID will be saved in a slice to later fetch
// the related tickets and users.
func (s *Storage) searchOrgByTerm(term, value string) ([]model.OrganizationResult, error) {
	type orgAndID struct {
		id  model.OrgID
		org model.Organization
	}
	var result []model.OrganizationResult
	var foundOrgs []*orgAndID

	// search all organizations for a match in a specific field
	for orgID, org := range s.OrganizationsMap {
		if findOrgMatch(org, term, value) {
			foundOrgs = append(foundOrgs, &orgAndID{
				id:  orgID,
				org: org,
			})
		}
	}

	for _, orgAndID := range foundOrgs {
		orgResult := model.OrganizationResult{
			Organization:   orgAndID.org,
			UserNames:      s.getUsersForOrg(orgAndID.id),
			TicketSubjects: s.getTicketsForOrg(orgAndID.id),
		}

		result = append(result, orgResult)
	}

	return result, nil
}

func findOrgMatch(org model.Organization, term, value string) bool {
	foundMatch := false

	switch v := org[term].(type) {
	case string:
		foundMatch = v == value
	case int:
		foundMatch = strconv.Itoa(v) == value
	case float64:
		// assume there are no decimals
		foundMatch = strconv.Itoa(int(v)) == value
	case bool:
		foundMatch = strconv.FormatBool(v) == value
	case []interface{}:
		s := ""
		for _, elem := range v {
			// assume they are strings and try to format them
			// and append them to s
			s = fmt.Sprintf("%s;%s", s, elem)
		}

		foundMatch = strings.Contains(s, value)
	default:
		fmt.Printf("unhandled type for v: %T\n", v)
	}

	return foundMatch
}

func (s *Storage) getUsersForOrg(orgID model.OrgID) []string {
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

	return userNames
}

func (s *Storage) getTicketsForOrg(orgID model.OrgID) []string {
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

	return ticketSubjects
}
