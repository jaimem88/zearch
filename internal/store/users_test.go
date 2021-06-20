package store

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jaimem88/zearch/internal/model"
)

func TestStorage_Users(t *testing.T) {
	orgData := readOrgs(t)
	userData := readUsers(t)
	ticketData := readTickets(t)

	tests := []struct {
		name           string
		userData       model.Users
		ticketData     model.Tickets
		orgData        model.Organizations
		term           string
		value          string
		expectedResult []*model.UserResult
		expectedError  error
	}{
		{
			name:     "search_by_id_with_no_organizations_no_tickets",
			userData: userData,
			term:     "_id",
			value:    "1",
			expectedResult: []*model.UserResult{
				{
					User: userData[0],
				},
			},
		},
		{
			name:           "search_by_id_does_not_exist",
			userData:       userData,
			term:           "_id",
			value:          "0",
			expectedResult: nil,
			expectedError:  ErrNotFound,
		},
		{
			name:           "search_by_term_does_not_exist",
			userData:       userData,
			term:           "unknown",
			expectedResult: nil,
			expectedError:  ErrNotFound,
		},
		{
			name:       "search_by_id_with_org_and_ticket",
			orgData:    orgData,
			userData:   userData,
			ticketData: ticketData,
			term:       "_id",
			value:      "1",
			expectedResult: []*model.UserResult{
				{
					User:             userData[0],
					OrganizationName: "Enthaze",
					TicketSubjects: []string{
						"A Problem in Guyana",
					},
				},
			},
		},
		{
			name:     "search_by_verified_boolean",
			userData: userData,
			term:     "verified",
			value:    "true",
			expectedResult: []*model.UserResult{
				{
					User: userData[1],
				},
			},
		},
		{
			name:     "search_by_tag_string_slice",
			userData: userData,
			term:     "tags",
			value:    "Leola",
			expectedResult: []*model.UserResult{
				{
					User: userData[0],
				},
				{
					User: userData[1],
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.orgData, tt.userData, tt.ticketData)
			got, err := s.Users(tt.term, tt.value)
			if tt.expectedError != nil {
				require.EqualError(t, err, tt.expectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, len(tt.expectedResult))

			// sort by user ID to ensure comparison below will be deterministic
			sort.SliceStable(got, func(i int, j int) bool {
				return got[i].User["_id"].(float64) <= got[j].User["_id"].(float64)
			})

			for k, expectedResult := range tt.expectedResult {
				assert.Equal(t, expectedResult.User["_id"], got[k].User["_id"])
				assert.Equal(t, expectedResult.OrganizationName, got[k].OrganizationName)
				assert.ElementsMatch(t, expectedResult.TicketSubjects, got[k].TicketSubjects)
			}
		})
	}
}
