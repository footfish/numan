package app

import (
	"github.com/footfish/numan"
	"github.com/footfish/numan/internal/storage"
)

// numanService implements the numanAPI
type numanService struct {
	db numan.API
}

// NewNumanService instantiates a new NumberService.
func NewNumanService(dsn string) numan.API {
	return &numanService{
		db: storage.NewStore(dsn),
	}
}

//Close closes db connection
func (s *numanService) Close() {
	s.db.Close()
}
