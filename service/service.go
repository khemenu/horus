package service

import "khepri.dev/horus"

type services struct {
	auth horus.AuthService
}

func (s *services) Auth() horus.AuthService {
	return s.auth
}
