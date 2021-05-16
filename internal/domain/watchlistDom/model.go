package watchlistDom

type Creation struct {
	MovieID string `json:"movieId"`
	UserID  uint   `json:"userId"`
}

type Delete struct {
	MovieID string `json:"movieId"`
	UserID  uint   `json:"userId"`
}
