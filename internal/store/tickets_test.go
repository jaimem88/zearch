package store

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jaimem88/zearch/internal/model"
)

func TestStorage_Tickets(t *testing.T) {
	orgData := readOrgs(t)
	ticketData := readTickets(t)

	tests := []struct {
		name string
		//userData       model.Users
		ticketData     model.Tickets
		orgData        model.Organizations
		term           string
		value          string
		expectedResult []*model.TicketResult
		expectedError  error
	}{
		{
			name:       "search_by_id_with_no_organizations",
			ticketData: ticketData,
			term:       "_id",
			value:      "27c447d9-cfda-4415-9a72-d5aa12942cf1",
			expectedResult: []*model.TicketResult{
				{
					Ticket: ticketData[0],
				},
			},
		},
		{
			name:           "search_by_id_does_not_exist",
			ticketData:     ticketData,
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
			name:       "search_by_id_with_org",
			orgData:    orgData,
			ticketData: ticketData,
			term:       "_id",
			value:      "27c447d9-cfda-4415-9a72-d5aa12942cf1",
			expectedResult: []*model.TicketResult{
				{
					Ticket:           ticketData[0],
					OrganizationName: "Enthaze",
				},
			},
		},
		{
			name:       "search_by_has_incidents_boolean",
			ticketData: ticketData,
			term:       "has_incidents",
			value:      "false",
			expectedResult: []*model.TicketResult{
				{
					Ticket: ticketData[1],
				},
			},
		},
		{
			name:       "search_by_tag_string_slice",
			ticketData: ticketData,
			term:       "tags",
			value:      "Massachusetts",
			expectedResult: []*model.TicketResult{
				{
					Ticket: ticketData[0],
				},
				{
					Ticket: ticketData[1],
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.orgData, nil, tt.ticketData)
			got, err := s.Tickets(tt.term, tt.value)
			if tt.expectedError != nil {
				require.EqualError(t, err, tt.expectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, len(tt.expectedResult))

			// sort by user ID to ensure comparison below will be deterministic
			sort.SliceStable(got, func(i int, j int) bool {
				return got[i].Ticket["_id"].(string) <= got[j].Ticket["_id"].(string)
			})

			for k, expectedResult := range tt.expectedResult {
				assert.Equal(t, expectedResult.Ticket["_id"], got[k].Ticket["_id"])
				assert.Equal(t, expectedResult.OrganizationName, got[k].OrganizationName)
			}
		})
	}
}
