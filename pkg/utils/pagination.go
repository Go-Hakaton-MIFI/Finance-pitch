package utils

type PaginatedEntities[T any] struct {
	Items            []T
	Total            int
	PageNumber       int
	ObjectsCount     int
	ObjectsCounTotal int
	PageCount        int
}

type RestfullPaginatedEntities[T any] struct {
	Items    []T
	Next     *string
	Previous *string
}
