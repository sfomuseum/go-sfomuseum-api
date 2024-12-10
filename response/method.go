package response

type Parameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    int    `json:"required"`
	Example     any    `json:"example"`
}

type Error struct {
	Message    string `json:"message"`
	Documented int    `json:"documented"`
}

type Method struct {
	Name                   string      `json:"name"`
	RequestMethod          string      `json:"request_method"`
	Description            string      `json:"description"`
	RequiresAuthentication int         `json:"requires_auth"`
	Parameters             []Parameter `json:"parameters"`
	// Errors                 map[string]Error `json:"errors"`
	Notes     []string `json:"notes"`
	Extras    int      `json:"extras"`
	Paginated int      `json:"paginated"`
}