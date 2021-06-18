package store

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jaimem88/zearch/internal/model"
)

// Users implements the searcher method for the app. It searches by term and value.
// Handles a special case for _id which can be looked up in the Storage easily from the
// UsersMap.
func (s *Storage) Users(term, value string) ([]model.UserResult, error) {
	fmt.Printf("Searching users by: %q with value: %q\n", term, value)

	if term == "_id" {
		return s.searchUserByID(value)
	}

	return s.searchUserByTerm(term, value)
}

func (s *Storage) searchUserByID(value string) ([]model.UserResult, error) {
	var results []model.UserResult

	id, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}

	userID := model.UserID(id)

	user, ok := s.UsersMap[userID]
	if !ok {
		return nil, ErrNotFound
	}

	orgID := getUserOrgID(user)
	userResult := model.UserResult{
		User:             user,
		OrganizationName: s.getOrgName(orgID),
		TicketsForOrg:    s.getTicketsForOrg(orgID),
	}

	results = append(results, userResult)

	return results, nil
}

func getUserOrgID(user model.User) model.OrgID {
	orgID, ok := user["organization_id"].(float64)
	if !ok {
		orgID = 0
	}

	return model.OrgID(orgID)
}

// searchUserByTerm will iterate over each element of the UsersMap and accessing
// the term directly. Once found, the user and its ID will be saved in a slice to later fetch
// the related tickets and users.
func (s *Storage) searchUserByTerm(term, value string) ([]model.UserResult, error) {
	type userAndID struct {
		id   model.UserID
		user model.User
	}
	var result []model.UserResult
	var foundUsers []*userAndID

	// search all users for a match in a specific field
	for userID, user := range s.UsersMap {
		if user[term] == nil {
			continue
		}

		if findUserMatch(user, term, value) {
			foundUsers = append(foundUsers, &userAndID{
				id:   userID,
				user: user,
			})
		}
	}

	for _, userAndID := range foundUsers {
		orgID := getUserOrgID(userAndID.user)
		userResult := model.UserResult{
			User:             userAndID.user,
			OrganizationName: s.getOrgName(orgID),
			TicketsForOrg:    s.getTicketsForOrg(orgID),
		}

		result = append(result, userResult)
	}

	return result, nil
}

func findUserMatch(user model.User, term, value string) bool {
	foundMatch := false

	switch v := user[term].(type) {
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

func (s *Storage) getOrgName(orgID model.OrgID) string {
	org, ok := s.OrganizationsMap[orgID]
	if !ok {
		return ""
	}

	// assume name is always a string
	return org["name"].(string)
}
