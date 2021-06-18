package model

import "strings"

//
// These types serve as aliases to help read the code
type (
	UserID   int
	OrgID    int
	TicketID string
)

//
//type Organization struct {
//	ID            OrgID    `json:"_id"`
//	URL           string   `json:"url"`
//	ExternalID    string   `json:"external_id"`
//	Name          string   `json:"name"`
//	DomainNames   []string `json:"domain_names"`
//	CreatedAt     string   `json:"created_at"`
//	Details       string   `json:"details"`
//	SharedTickets bool     `json:"shared_tickets"`
//	Tags          []string `json:"tags"`
//}
//
//// GetJSONTagsFromStruct users refecltion to iterate over t and extract the JSON tag for each field.
//// Assumes the JSON tag is one word, no commas or other fields like `omitempty`.
//// Does not support nested structs
//func GetJSONTagsFromStruct(t interface{}) string {
//	sb := strings.Builder{}
//
//	val := reflect.ValueOf(t)
//	for i := 0; i < val.Type().NumField(); i++ {
//
//		sb.WriteString(val.Type().Field(i).Tag.Get("json"))
//		sb.WriteRune('\n')
//	}
//
//	return sb.String()
//}
//
type OrganizationResult struct {
	Organization
	UserNames      []string
	TicketSubjects []string
}

//
//type Ticket struct {
//	ID             TicketID `json:"_id"`
//	URL            string   `json:"url"`
//	ExternalID     string   `json:"external_id"`
//	CreatedAt      string   `json:"created_at"`
//	Type           string   `json:"type"`
//	Subject        string   `json:"subject"`
//	Description    string   `json:"description"`
//	Priority       string   `json:"priority"`
//	Status         string   `json:"status"`
//	SubmitterID    int      `json:"submitter_id"`
//	AssigneeID     int      `json:"assignee_id"`
//	OrganizationID OrgID    `json:"organization_id"`
//	Tags           []string `json:"tags"`
//	HasIncidents   bool     `json:"has_incidents"`
//	DueAt          string   `json:"due_at"`
//	Via            string   `json:"via"`
//}
//
type TicketResult struct {
	Ticket
	OrganizationName  string
	TicketDescription []string
}

//
//type User struct {
//	ID             UserID   `json:"_id"`
//	URL            string   `json:"url"`
//	ExternalID     string   `json:"external_id"`
//	Name           string   `json:"name"`
//	Alias          string   `json:"alias"`
//	CreatedAt      string   `json:"created_at"`
//	Active         bool     `json:"active"`
//	Verified       bool     `json:"verified"`
//	Shared         bool     `json:"shared"`
//	Locale         string   `json:"locale"`
//	Timezone       string   `json:"timezone"`
//	LastLoginAt    string   `json:"last_login_at"`
//	Email          string   `json:"email"`
//	Phone          string   `json:"phone"`
//	Signature      string   `json:"signature"`
//	OrganizationID OrgID    `json:"organization_id"`
//	Tags           []string `json:"tags"`
//	Suspended      bool     `json:"suspended"`
//	Role           string   `json:"role"`
//}
//
type UserResult struct {
	User
	OrganizationName string
	Tickets          []string
}

type Organization map[string]interface{}
type User map[string]interface{}
type Ticket map[string]interface{}
type Organizations []Organization
type Users []User
type Tickets []Ticket

func (o Organization) String() string {
	sb := strings.Builder{}
	for field, _ := range o {
		sb.WriteString(field)
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (u User) String() string {
	sb := strings.Builder{}
	for field, _ := range u {
		sb.WriteString(field)
		sb.WriteRune('\n')
	}

	return sb.String()
}
func (t Ticket) String() string {
	sb := strings.Builder{}
	for field, _ := range t {
		sb.WriteString(field)
		sb.WriteRune('\n')
	}

	return sb.String()
}
