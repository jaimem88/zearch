package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/jaimem88/zearch/internal/cli"
	"github.com/jaimem88/zearch/internal/model"
	"github.com/jaimem88/zearch/internal/parser"
	"github.com/jaimem88/zearch/internal/store"
)

var (
	usersFilename   = flag.String("users", "data/users.json", "Filename to load users from e.g. --users data/users.json")
	ticketsFilename = flag.String("tickets", "data/tickets.json", "Filename to load users from e.g. --users data/users.json")
	orgsFilename    = flag.String("organizations", "data/organizations.json", "Filename to load users from e.g. --users data/users.json")
)

func main() {
	flag.Parse()

	////var i []map[string]interface{}
	//var i model.Users
	//err := parser.ReadJSONFile(*usersFilename, &i)
	//if err != nil {
	//	log.Fatalf("failed to load users.json: %+v", err)
	//}
	//for key, val := range i {
	//	fmt.Printf("doc: %d json: %+v id %+v\n", key, val, val["_id"])
	//	//for k, v := range val {
	//	//	fmt.Printf("map[%s]%+v\n", k, v)
	//	//
	//	//}
	//}
	//fmt.Printf("CAN I ACCESS A FIELD: %+v\n", i[0]["_id"])
	//
	////fmt.Printf("What happened?: %+v\n\ntypeOf: %q", i, reflect.TypeOf(i.([]interface{})[0]))
	//return
	var users model.Users
	var tickets model.Tickets
	var orgs model.Organizations
	err := parser.ReadJSONFile(*usersFilename, &users)
	if err != nil {
		log.Fatalf("failed to load users.json: %+v", err)
	}

	//fmt.Printf("CAN I ACCESS A FIELD: %+v\n", users[0]["_id"])
	//return
	err = parser.ReadJSONFile(*ticketsFilename, &tickets)
	if err != nil {
		log.Fatalf("failed to load tickets.json: %+v", err)
	}
	err = parser.ReadJSONFile(*orgsFilename, &orgs)
	if err != nil {
		log.Fatalf("failed to load organizations.json: %+v", err)
	}

	c := cli.New(store.New(orgs, users, tickets))
	p := promptui.Prompt{
		Label: "Hi Zendesk! Press return to continue",
		Templates: &promptui.PromptTemplates{
			Prompt: "👋 {{ . | blue | bold }}",
		},
	}
	_, err = p.Run()
	if err != nil {
		log.Fatalf(err.Error())
	}

	selectTemplate := &promptui.SelectTemplates{
		Active:   `👉 {{ . | cyan | bold }}`,
		Inactive: `  {{ . }}`,
		Selected: `✅ {{ . }}`,
	}

	prompt := promptui.Select{
		Label:     "What would you like to do?",
		Items:     []string{"Zearch Zendesk", "View searchable fields", "Quit"},
		Templates: selectTemplate,
	}

	stop := false
	for !stop {
		n, _, err := prompt.Run()
		if err != nil {
			log.Fatalf(err.Error())
		}

		switch n {
		case 0:
			selectEntity := promptui.Select{
				Label:     "Select a search option:",
				Items:     []string{"Users", "Tickets", "Organizations"},
				Templates: selectTemplate,
			}
			_, entity, err := selectEntity.Run()
			if err != nil {
				log.Fatalf(err.Error())
			}

			promptTerm := promptui.Prompt{
				Label: "Type search term:",
			}
			term, err := promptTerm.Run()
			if err != nil {
				log.Fatalf(err.Error())
			}
			promptValue := promptui.Prompt{
				Label: "Type search value:",
			}
			value, err := promptValue.Run()
			if err != nil {
				log.Fatalf(err.Error())
			}
			c.Search(entity, term, value)
		case 1:
			c.PrintSearchableFields(orgs[0], users[0], tickets[0])
		case 2:
			p := promptui.Prompt{
				Label:     "Are you sure you want to quit??",
				IsConfirm: true,
			}

			quit, _ := p.Run()
			if strings.ToLower(quit) == "y" {
				fmt.Printf("See ya!")
				stop = true
			}
		}
	}
}
