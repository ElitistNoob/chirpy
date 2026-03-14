package app

import (
	"sync/atomic"

	"github.com/ElitistNoob/chirpy/internal/database"
)

type App struct {
	Queries        *database.Queries
	FileserverHits atomic.Int32
	Platform       string
}

func NewApp() *App {
	return &App{}
}
