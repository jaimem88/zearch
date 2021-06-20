package model

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
	OrganizationName string
}

type UserResult struct {
	User
	OrganizationName string
	TicketSubjects   []string
}

type Organization map[string]interface{}
type User map[string]interface{}
type Ticket map[string]interface{}
type Organizations []Organization
type Users []User
type Tickets []Ticket
