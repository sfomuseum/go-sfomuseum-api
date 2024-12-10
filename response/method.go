package response

import (
	"fmt"
	"log/slog"
	"strconv"
)

type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    int    `json:"required"` // sudo make me a boolean
	Example     any    `json:"example"`  // sudo make me a string
}

func (p *Parameter) String() (string, error) {

	switch p.Type {
	case "int64":

		i, ok := p.Example.(int64)

		if !ok {
			return "", fmt.Errorf("Not a valid int64 (%T)", p.Example)
		}

		return strconv.FormatInt(i, 10), nil

	case "float64":

		f, ok := p.Example.(float64)

		if !ok {
			return "", fmt.Errorf("Not a float64 (%T)", p.Example)
		}

		return strconv.FormatFloat(f, 'g', -1, 64), nil

	case "string":

		s, ok := p.Example.(string)

		if !ok {
			return "", fmt.Errorf("Not a string (%T)", p.Example)
		}

		return s, nil

	default:
		slog.Warn("Unhandled API parameter type", "type", p.Type)
		return fmt.Sprintf("%v", p.Example), nil
	}
}

type Error struct {
	Message    string `json:"message"`
	Documented int    `json:"documented"`
}

type Method struct {
	Name                   string      `json:"name"`
	RequestMethod          string      `json:"request_method"`
	Description            string      `json:"description"`
	RequiresAuthentication int         `json:"requires_auth"` // sudo make me a boolean
	Parameters             []Parameter `json:"parameters"`
	// Errors                 map[string]Error `json:"errors"`
	Notes     []string `json:"notes"`
	Extras    int      `json:"extras"`    // sudo make me a boolean
	Paginated int      `json:"paginated"` // sudo make me a boolean
}
