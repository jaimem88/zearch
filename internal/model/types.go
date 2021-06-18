package model

import (
	"sort"
	"strings"
)

// These types serve as aliases to help read the code
type (
	UserID   float64
	OrgID    float64
	TicketID string
)

// OrganizationResult contains the result of a search
type OrganizationResult struct {
	Organization
	UserNames      []string
	TicketSubjects []string
}

type TicketResult struct {
	Ticket
	OrganizationName  string
	TicketDescription []string
}

type UserResult struct {
	User
	OrganizationName string
	TicketsForOrg    []string
}

type Organization map[string]interface{}
type User map[string]interface{}
type Ticket map[string]interface{}
type Organizations []Organization
type Users []User
type Tickets []Ticket

// String implements the stringer interface.
// Will print the key names from the Organization.
func (o Organization) String() string {
	items := make([]string, 0, len(o))
	for field := range o {
		items = append(items, field)
	}

	sort.Strings(items)

	return strings.Join(items, "\n")
}

// String implements the stringer interface.
// Will print the key names from the User.
func (u User) String() string {
	items := make([]string, 0, len(u))
	for field := range u {
		items = append(items, field)
	}

	sort.Strings(items)

	return strings.Join(items, "\n")
}

// String implements the stringer interface.
// Will print the key names from the Ticket.
func (t Ticket) String() string {
	items := make([]string, 0, len(t))
	for field := range t {
		items = append(items, field)
	}

	sort.Strings(items)

	return strings.Join(items, "\n")
}
