package recommendationDom

type Creation struct {
	SenderID  uint   `json:"-"`
	MovieID   uint   `json:"movieId"`
	Comment   string `json:"comment"`
	CircleIDs []uint `json:"circleIds"`
	UserIDs   []uint `json:"userIds"`
}

func (c Creation) Valid() error {
	if c.SenderID == 0 {
		return errSenderIDNull
	}
	if c.MovieID == 0 {
		return errMovieIDNull
	}
	if c.Comment == "" {
		return errCommentEmpty
	}
	if len(c.CircleIDs) == 0 && len(c.UserIDs) == 0 {
		return errMissingRecipient
	}
	return nil
}
