package app

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/jaimem88/zearch/internal/model"
	"github.com/jaimem88/zearch/internal/store"
)

var (
	selectTemplate = &promptui.SelectTemplates{
		Active:   `👉 {{ . | cyan | bold }}`,
		Inactive: `  {{ . }}`,
		Selected: `✅ {{ . }}`,
	}
)

// Storage defines the methods that the App store requires in order to get
// the Organizations, Users and Tickets from the underlying storage.
type Storage interface {
	Organizations(term, value string) ([]model.OrganizationResult, error)
	Users(term, value string) ([]model.UserResult, error)
	Tickets(term, value string) ([]model.TicketResult, error)
}

// App handles the CLI interaction with the user and does the
// information presentation to stdout
type App struct {
	store Storage
}

// New creates an App with the defined Storage
func New(store Storage) *App {
	return &App{
		store: store,
	}
}

// Run the App and handle user input vua promptui
func (a *App) Run() error {
	welcomePrompt := promptui.Prompt{
		Label: "Hi Zendesk! Press return to continue",
		Templates: &promptui.PromptTemplates{
			Prompt: "👋 {{ . | blue | bold }}",
		},
	}

	_, err := welcomePrompt.Run()
	if err != nil {
		return err
	}

	actionPrompt := promptui.Select{
		Label:     "What would you like to do?",
		Items:     []string{"Zearch Zendesk", "View searchable fields", "Quit"},
		Templates: selectTemplate,
	}

	stop := false
	for !stop {
		n, _, err := actionPrompt.Run()
		if err != nil {
			return err
		}

		switch n {
		case 0:
			if err := a.handleSearch(); err != nil {
				return fmt.Errorf("search failed: %w", err)
			}
		case 1:
			//a.PrintSearchableFields(orgs[0], users[0], tickets[0])
		case 2:
			confirmQuit := promptui.Prompt{
				Label:     "Are you sure you want to quit??",
				IsConfirm: true,
			}

			quit, err := confirmQuit.Run()
			if err != nil && !errors.Is(err, promptui.ErrAbort) {
				// promptui.Prompt returns an empty error when IsConfirm: true https://github.com/manifoldco/promptui/issues/81
				// we need to check the error type to properly handle the error instead of ignoring it
				return err
			}

			if strings.ToLower(quit) == "y" {
				fmt.Printf("See ya!")
				stop = true
			}
		default:
			return fmt.Errorf("unknown option: %d\n", n)
		}
	}

	return nil
}

func (a *App) handleSearch() error {
	selectEntity := promptui.Select{
		Label:     "Select a search option:",
		Items:     []string{"Users", "Tickets", "Organizations"},
		Templates: selectTemplate,
	}

	_, entity, err := selectEntity.Run()
	if err != nil {
		return err
	}

	promptTerm := promptui.Prompt{
		Label: "Type search term:",
	}

	term, err := promptTerm.Run()
	if err != nil {
		return err
	}

	promptValue := promptui.Prompt{
		Label: "Type search value:",
	}

	value, err := promptValue.Run()
	if err != nil {
		return err
	}

	return a.Search(entity, term, value)
}

func (a *App) handleQuit() (bool, error) {
	confirmQuit := promptui.Prompt{
		Label:     "Are you sure you want to quit??",
		IsConfirm: true,
	}

	quit, err := confirmQuit.Run()
	if err != nil {
		return false, err
	}

	if strings.ToLower(quit) == "y" {
		fmt.Printf("See ya!")
		return true, nil
	}

	return false, nil
}

func (a *App) Search(entity, term, value string) error {
	printDashes(80)

	switch strings.ToLower(entity) {
	case "organizations":
		orgResults, err := a.store.Organizations(term, value)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				fmt.Println("No results found")
				return nil
			}

			return err
		}

		for _, orgResult := range orgResults {
			err := model.OrgResultTemplate.Execute(os.Stdout, orgResult)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Total organizations found: %d\n", len(orgResults))
	case "users":
		userResults, err := a.store.Users(term, value)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				fmt.Println("No results found")
				return nil
			}

			return err
		}

		for _, userResult := range userResults {
			err := model.UserResultTemplate.Execute(os.Stdout, userResult)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Total users found: %d\n", len(userResults))
	case "tickets":
		ticketResults, err := a.store.Tickets(term, value)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				fmt.Println("No results found")
				return nil
			}

			return err
		}

		for _, ticketResult := range ticketResults {
			err := model.TicketResultTemplate.Execute(os.Stdout, ticketResult)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Total tickets found: %d\n", len(ticketResults))
	default:
		return fmt.Errorf("unkown entity: %s", entity)
	}

	return nil
}

func (a *App) printSearchableFields(organization model.Organization, user model.User, ticket model.Ticket) {
	printDashes(80)
	printFields("Organizations", organization.String())

	printDashes(80)
	printFields("Users", user.String())

	printDashes(80)
	printFields("Tickets", ticket.String())
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
