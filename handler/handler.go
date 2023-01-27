package handler

type Handler struct {
	db *Store.DataBase
}

func NewHandler(db *Store.DataBase) *Handler {
	return &Handler{db}
}
