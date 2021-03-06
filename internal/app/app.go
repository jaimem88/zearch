package app

import (
	"errors"
	"fmt"
	"io"
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
	GetSearchableFields() map[string][]string
}

// App handles the CLI interaction with the user and does the
// information presentation to stdout
type App struct {
	store Storage
	out   io.Writer
}

// New creates an App with the defined Storage
func New(store Storage, out io.Writer) *App {
	return &App{
		store: store,
		out:   out,
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
			a.printSearchableFields()
		case 2:
			stop, err = a.handleQuit()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown option: %d", n)
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
	if err != nil && !errors.Is(err, promptui.ErrAbort) {
		// promptui.Prompt returns an empty error when IsConfirm: true https://github.com/manifoldco/promptui/issues/81
		// we need to check the error type to properly handle the error instead of ignoring it
		return false, err
	}

	if strings.ToLower(quit) == "y" {
		fmt.Fprintf(a.out, "See ya!")
		return true, nil
	}

	return false, nil
}

func (a *App) Search(entity, term string, value string) error {
	a.printDashes(80)

	switch strings.ToLower(entity) {
	case "organizations":
		return a.searchOrganizations(term, value)
	case "users":
		return a.searchUsers(term, value)
	case "tickets":
		return a.searchTickets(term, value)
	default:
		return fmt.Errorf("unkoown entity: %s", entity)
	}
}

func (a *App) searchOrganizations(term, value string) error {
	var orgResults []model.OrganizationResult
	// TODO: handle `and` properly
	terms := strings.Split(term, " or ")
	values := strings.Split(value, " or ")

	if len(terms) != len(values) {
		return fmt.Errorf("%d terms do not match %d values", len(terms), len(values))
	}

	for k, term := range terms {
		orgResultsSubset, err := a.store.Organizations(term, values[k])
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				fmt.Fprintf(a.out, "No results found")
				return nil
			}

			return err
		}

		orgResults = append(orgResults, orgResultsSubset...)
	}

	for _, orgResult := range orgResults {
		err := model.OrgResultTemplate.Execute(a.out, orgResult)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(a.out, "Total organizations found: %d\n", len(orgResults))

	return nil
}
func (a *App) searchUsers(term, value string) error {
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

	fmt.Fprintf(a.out, "Total users found: %d\n", len(userResults))

	return nil
}

func (a *App) searchTickets(term, value string) error {
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

	fmt.Fprintf(a.out, "Total tickets found: %d\n", len(ticketResults))

	return nil
}

func (a *App) printSearchableFields() {
	a.printDashes(80)
	fields := a.store.GetSearchableFields()
	a.printFields("Organizations", fields["organizations"])

	a.printDashes(80)
	a.printFields("Users", fields["users"])

	a.printDashes(80)
	a.printFields("Tickets", fields["tickets"])
}

func (a *App) printFields(param string, fields []string) {
	fmt.Fprintf(a.out, "Search %s by:\n%s", param, strings.Join(fields, "\n"))
}

func (a *App) printDashes(n int) {
	fmt.Fprintln(a.out)
	for n > 0 {
		n--
		fmt.Fprint(a.out, "-")
	}

	fmt.Fprintln(a.out)
}
