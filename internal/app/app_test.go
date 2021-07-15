package app

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jaimem88/zearch/internal/model"
)

func TestSearch_ByOrganization(t *testing.T) {
	buf := &bytes.Buffer{}
	app := New(&mockStore{
		orgResults: []model.OrganizationResult{
			{
				Organization: model.Organization{
					"_id":         125,
					"url":         "http://initech.zendesk.com/api/v2/organizations/125.json",
					"external_id": "42a1a845-70cf-40ed-a762-acb27fd606cc",
					"name":        "Strezzö",
					"domain_names": []string{
						"techtrix.com",
						"teraprene.com",
						"corpulse.com",
						"flotonic.com",
					},
					"created_at":     "2016-02-21T06:11:51 -11:00",
					"details":        "MegaCorp",
					"shared_tickets": false,
					"tags": []string{
						"Vance",
						"Ray",
						"Jacobs",
						"Frank",
					},
				},
			},
			{
				Organization: model.Organization{
					"_id":         124,
					"url":         "http://initech.zendesk.com/api/v2/organizations/124.json",
					"external_id": "15c21605-cbc6-440f-8da2-6e1601aed5fa",
					"name":        "Bitrex",
					"domain_names": []string{
						"unisure.com",
						"boink.com",
						"quinex.com",
						"poochies.com",
					},
					"created_at":     "2016-05-11T12:16:15 -10:00",
					"details":        "Non profit",
					"shared_tickets": true,
					"tags": []string{
						"Lott",
						"Hunter",
						"Beasley",
						"Glass",
					},
				},
			},
		},
	},
		buf)

	err := app.Search("organizations", "name or name", "Bitrex or Strezzö")
	if err != nil {
		t.Errorf("%+v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `Bitrex`) {
		t.Error("output does not contain name")
	}

	if !strings.Contains(out, `Strezzö`) {
		t.Error("output does not contain name")
	}

	err = app.Search("organizations", "_id or name", "Strezzö")
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if !(err.Error() == "2 terms do not match 1 values") {
		t.Fatalf("expected error message does not match: %q", err)
	}
}

type mockStore struct {
	orgResults []model.OrganizationResult
	err        error
}

func (ms *mockStore) Organizations(term, value string) ([]model.OrganizationResult, error) {
	return ms.orgResults, ms.err
}

func (ms *mockStore) Users(term, value string) ([]model.UserResult, error) {

	return nil, nil
}
func (ms *mockStore) Tickets(term, value string) ([]model.TicketResult, error) {

	return nil, nil
}
func (ms *mockStore) GetSearchableFields() map[string][]string {

	return nil
}
