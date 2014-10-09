package revelswagger

type Specification struct {
	Swagger  int64
	Info     info
	Host     string
	Schemes  []string
	BasePath string
	Produces []string
	Paths    map[string]route
}

type info struct {
	Title       string
	Description string
	Version     string
}

type route struct {
	Get *method
}

type method struct {
	Summary     string
	Description string
	Parameters  []parameter
	Tags        []string
	Responses   map[string]response
}

type parameter struct {
	Name        string
	In          string
	Description string
	Required    bool
	Type        string
	Format      string
	Minimum     *int
	Maximum     *int
	Pattern     string
	Enum        []string
	Default     string
}

type response struct {
	Description string
	Schema      schema
}

type schema struct {
	Type  string
	Items map[string]string
}
