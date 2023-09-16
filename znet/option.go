package znet

type Option func(s *Server)

func WithName(name string) Option {
	return func(s *Server) {
		s.Name = name
	}
}

func WithIP(ip string) Option {
	return func(s *Server) {
		s.IP = ip
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}
