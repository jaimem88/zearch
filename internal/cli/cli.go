package cli

import (
	"fmt"
	"strings"

	"github.com/jaimem88/zearch/internal/model"
)

type Searcher interface {
	Organizations(term string, value interface{}) *model.OrganizationResult
	Tickets(term string, value interface{}) *model.TicketResult
	Users(term string, value interface{}) *model.UserResult
}

type CLI struct {
	searcher Searcher
}

func New(searcher Searcher) *CLI {
	return &CLI{
		searcher: searcher,
	}
}

func (c *CLI) Search(entity, term string, value interface{}) (interface{}, error) {
	switch strings.ToLower(entity) {
	case "organizations":
		return c.searcher.Organizations(term, value), nil
	case "tickets":
		return c.searcher.Tickets(term, value), nil
	case "users":
		return c.searcher.Users(term, value), nil
	default:
		return nil, fmt.Errorf("unknown option: %s")
	}
}

func (c *CLI) SearchableFields() (interface{}, error) {
	return nil, nil
}
