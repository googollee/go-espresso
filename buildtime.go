package espresso

type buildtimeEndpoint struct {
	method      string
	path        string
	pathBinding map[string]binder
}
