package chirps

import (
	"net/http"

	"github.com/google/uuid"
)

func getAuthorIdFromRequest(r *http.Request) (uuid.UUID, error) {
	urlParam := r.URL.Query().Get("author_id")
	if urlParam == "" {
		return uuid.Nil, nil
	}

	authorID, err := uuid.Parse(urlParam)
	if err != nil {
		return uuid.Nil, err
	}

	return authorID, nil
}

func getSortOrderFromRequest(r *http.Request) bool {
	sortParam := r.URL.Query().Get("sort")
	if sortParam == "" || sortParam == "asc" {
		return false
	}
	return true
}
