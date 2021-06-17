package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jaimem88/zearch/internal/parser"

	"github.com/jaimem88/zearch/internal/model"
)

func TestStorage_Organizations(t *testing.T) {
	orgData := readOrgs(t)

	tests := []struct {
		name           string
		userData       []*model.User
		ticketData     []*model.Ticket
		orgData        []*model.Organization
		term           string
		value          string
		expectedResult []*model.OrganizationResult
		expectedError  error
	}{
		{
			name:    "by_id_with_no_users_no_tickets",
			orgData: orgData,
			term:    "_id",
			value:   "101",
			expectedResult: []*model.OrganizationResult{
				{
					Organization: orgData[0],
				},
			},
		},
		{
			name:           "by_id_does_not_exist",
			orgData:        orgData,
			term:           "_id",
			value:          "0",
			expectedResult: nil,
			expectedError:  ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.userData, tt.ticketData, tt.orgData)
			got, err := s.Organizations(tt.term, tt.value)
			if tt.expectedError != nil {
				require.EqualError(t, err, tt.expectedError.Error())
				return
			}

			require.NoError(t, err)

			require.Len(t, got, len(tt.expectedResult))
			assert.Len(t, got[0].UserNames, len(tt.expectedResult[0].UserNames))
			assert.Len(t, got[0].TicketSubjects, len(tt.expectedResult[0].TicketSubjects))
		})
	}
}

func readOrgs(t *testing.T) []*model.Organization {
	t.Helper()

	var orgs []*model.Organization
	err := parser.ReadJSONFile("testdata/organizations.json", &orgs)
	require.NoError(t, err)

	return orgs
}

func readUsers(t *testing.T) []*model.User {
	t.Helper()

	var users []*model.User
	err := parser.ReadJSONFile("testdata/users.json", &users)
	require.NoError(t, err)

	return users
}

func readTickets(t *testing.T) []*model.Ticket {
	t.Helper()

	var tickets []*model.Ticket
	err := parser.ReadJSONFile("testdata/tickets.json", &tickets)
	require.NoError(t, err)

	return tickets
}
