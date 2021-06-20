package store

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jaimem88/zearch/internal/model"
)

// Organizations implements the searcher method for the app. It searches by term and value.
// Handles a special case for _id which can be looked up in the Storage easily from the
// OrganizationsMap.
func (s *Storage) Organizations(term, value string) ([]model.OrganizationResult, error) {
	fmt.Printf("Searching organizations by: %q with value: %q\n", term, value)

	if term == "_id" {
		return s.searchOrgByID(value)
	}

	return s.searchOrgByTerm(term, value)
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
		if org[term] == nil {
			continue
		}

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

	if len(result) < 1 {
		return nil, ErrNotFound
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
