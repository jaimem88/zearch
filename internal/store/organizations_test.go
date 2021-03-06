package store

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jaimem88/zearch/internal/model"
)

func TestStorage_Organizations(t *testing.T) {
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
		expectedResult []*model.OrganizationResult
		expectedError  error
	}{
		{
			name:    "search_by_id_with_no_users_no_tickets",
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
			name:           "search_by_id_does_not_exist",
			orgData:        orgData,
			term:           "_id",
			value:          "0",
			expectedResult: nil,
			expectedError:  ErrNotFound,
		},
		{
			name:           "search_by_term_does_not_exist",
			ticketData:     ticketData,
			term:           "unknown",
			expectedResult: nil,
			expectedError:  ErrNotFound,
		},
		{
			name:       "search_by_id_with_users_and_tickets",
			orgData:    orgData,
			userData:   userData,
			ticketData: ticketData,
			term:       "_id",
			value:      "101",
			expectedResult: []*model.OrganizationResult{
				{
					Organization:   orgData[0],
					UserNames:      []string{"Francis Bailey"},
					TicketSubjects: []string{"A Problem in Guyana"},
				},
			},
		},
		{
			name:    "search_by_shared_tickets_boolean",
			orgData: orgData,
			term:    "shared_tickets",
			value:   "true",
			expectedResult: []*model.OrganizationResult{
				{
					Organization: orgData[1],
				},
			},
		},
		{
			name:    "search_by_tag_string_slice",
			orgData: orgData,
			term:    "tags",
			value:   "Farley",
			expectedResult: []*model.OrganizationResult{
				{
					Organization: orgData[0],
				},
				{
					Organization: orgData[1],
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.orgData, tt.userData, tt.ticketData)
			got, err := s.Organizations(tt.term, tt.value)
			if tt.expectedError != nil {
				require.EqualError(t, err, tt.expectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, len(tt.expectedResult))

			// sort by orgID to ensure comparison below will be deterministic
			sort.SliceStable(got, func(i int, j int) bool {
				return got[i].Organization["_id"].(float64) <= got[j].Organization["_id"].(float64)
			})

			for k, expectedResult := range tt.expectedResult {
				assert.Equal(t, expectedResult.Organization["_id"], got[k].Organization["_id"])
				assert.ElementsMatch(t, expectedResult.UserNames, got[k].UserNames)
				assert.ElementsMatch(t, expectedResult.TicketSubjects, got[k].TicketSubjects)
			}
		})
	}
}
