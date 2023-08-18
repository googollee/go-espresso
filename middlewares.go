package espresso

func WithoutDefault() ServerOption {
	return func(s *Server) error {
		s.middlewares = nil
		return nil
	}
}
