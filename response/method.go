package response

type Parameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    int    `json:"required"` // sudo make me a boolean
	Example     any    `json:"example"`  // sudo make me a string
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
