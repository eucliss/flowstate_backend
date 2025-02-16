package sources

type Source interface {
	Query(query Query) []interface{}
}

type SourceInitializer interface {
	Initialize() (Source, error)
}

type Query struct {
	Query      string
	SourceType string
	Limit      int
	Start      int
	End        int
}
