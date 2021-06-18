package store

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jaimem88/zearch/internal/model"
	"github.com/jaimem88/zearch/internal/reader"
)

func readOrgs(t *testing.T) model.Organizations {
	t.Helper()

	var orgs model.Organizations
	err := reader.ReadJSONFile("testdata/organizations.json", &orgs)
	require.NoError(t, err)

	return orgs
}

func readUsers(t *testing.T) model.Users {
	t.Helper()

	var users model.Users
	err := reader.ReadJSONFile("testdata/users.json", &users)
	require.NoError(t, err)

	return users
}

func readTickets(t *testing.T) model.Tickets {
	t.Helper()

	var tickets model.Tickets
	err := reader.ReadJSONFile("testdata/tickets.json", &tickets)
	require.NoError(t, err)

	return tickets
}
