package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jaimem88/zearch/internal/model"
	"github.com/jaimem88/zearch/internal/store"
)

type Searcher interface {
	Organizations(term, value string) ([]model.OrganizationResult, error)
	Tickets(term, value string) []model.TicketResult
	Users(term, value string) []model.UserResult
}

type CLI struct {
	searcher Searcher
}

func New(searcher Searcher) *CLI {
	return &CLI{
		searcher: searcher,
	}
}

func (c *CLI) Search(entity, term, value string) error {
	printDashes(80)
	switch strings.ToLower(entity) {
	case "organizations":
		r, err := c.searcher.Organizations(term, value)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				fmt.Println("No results found")
				return nil
			}

			return err
		}

		if len(r) == 0 {
			fmt.Println("No results found")
			return nil
		}

		for _, rr := range r {
			err := model.OrgResultTemplate.Execute(os.Stdout, rr)
			if err != nil {
				return err
			}
		}

	case "tickets":
		c.searcher.Tickets(term, value)
	case "users":
		c.searcher.Users(term, value)
	default:
		return nil
	}

	return nil
}

func (c *CLI) PrintSearchableFields(organization model.Organization, user model.User, ticket model.Ticket) {
	printDashes(80)
	printFields("Users", user.String())

	printDashes(80)
	printFields("Tickets", ticket.String())

	printDashes(80)
	printFields("Organizations", organization.String())
}

func printFields(param, fields string) {
	fmt.Printf("Search %s by:\n%s", param, fields)
}

func printDashes(n int) {
	fmt.Println()
	for n > 0 {
		n--
		fmt.Print("-")
	}
	fmt.Println()
}
