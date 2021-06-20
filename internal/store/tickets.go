package store

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jaimem88/zearch/internal/model"
)

// Tickets implements the searcher method for the app. It searches by term and value.
// Handles a special case for _id which can be looked up in the Storage easily from the
// TicketsMap.
func (s *Storage) Tickets(term, value string) ([]model.TicketResult, error) {
	fmt.Printf("Searching tickets by: %q with value: %q\n", term, value)

	if term == "_id" {
		return s.searchTicketByID(value)
	}

	return s.searchTicketByTerm(term, value)
}

func (s *Storage) searchTicketByID(value string) ([]model.TicketResult, error) {
	var results []model.TicketResult

	ticket, ok := s.TicketsMap[model.TicketID(value)]
	if !ok {
		return nil, ErrNotFound
	}

	orgID := getTicketOrgID(ticket)
	ticketResult := model.TicketResult{
		Ticket:           ticket,
		OrganizationName: s.getOrgName(orgID),
	}

	results = append(results, ticketResult)

	return results, nil
}

func getTicketOrgID(ticket model.Ticket) model.OrgID {
	orgID, ok := ticket["organization_id"].(float64)
	if !ok {
		orgID = 0
	}

	return model.OrgID(orgID)
}

// searchTicketByTerm will iterate over each element of the TicketMap and accessing
// the term directly.
func (s *Storage) searchTicketByTerm(term, value string) ([]model.TicketResult, error) {
	var result []model.TicketResult
	var foundTicket []model.Ticket

	// search all tickets for a match in a specific field
	for _, ticket := range s.TicketsMap {
		if ticket[term] == nil {
			continue
		}

		if findTicketMatch(ticket, term, value) {
			foundTicket = append(foundTicket, ticket)
		}
	}

	for _, ticket := range foundTicket {
		orgID := getTicketOrgID(ticket)
		ticketResult := model.TicketResult{
			Ticket:           ticket,
			OrganizationName: s.getOrgName(orgID),
		}

		result = append(result, ticketResult)
	}

	if len(result) < 1 {
		return nil, ErrNotFound
	}

	return result, nil
}

func findTicketMatch(ticket model.Ticket, term, value string) bool {
	foundMatch := false

	switch v := ticket[term].(type) {
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
