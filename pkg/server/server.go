package server

import "github.com/Malyue/quotaguard/pkg/quota"

type Server struct {
	QuotaManager *quota.QuotaManager
}

func NewServer() *Server {
	qm := quota.NewQuotaManager()
	return &Server{
		QuotaManager: qm,
	}
}
