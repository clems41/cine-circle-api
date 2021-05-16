package userDom

type CommonFields struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type Creation struct {
	CommonFields
	Username string `json:"username"`
	Password string `json:"password"`
}

type Update struct {
	UserID uint `json:"-"`
	CommonFields
}

type Get struct {
	UserID   uint   `json:"-"`
	Username string `json:"-"`
	Email    string `json:"-"`
}

type UpdatePassword struct {
	UserID      uint   `json:"-"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type Delete struct {
	UserID uint `json:"-"`
}

type View struct {
	UserID      uint   `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type ViewMe struct {
	UserID      uint   `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type Filters struct {
	Keyword string
}

func (c CommonFields) Valid() (err error) {
	if c.Email == "" {
		return errValidEmail
	}
	if c.DisplayName == "" {
		return errValidDisplayName
	}
	return nil
}

func (c Creation) Valid() (err error) {
	if c.Password == "" {
		return errValidPassword
	}
	if c.Username == "" {
		return errValidUsername
	}
	return c.CommonFields.Valid()
}

func (u Update) Valid() (err error) {
	if u.UserID == 0 {
		return errValidID
	}
	return u.CommonFields.Valid()
}

func (d Delete) Valid() (err error) {
	if d.UserID == 0 {
		err = errValidID
	}
	return
}

func (g Get) Valid() (err error) {
	if g.UserID == 0 && g.Username == "" && g.Email == "" {
		err = errValidGet
	}
	return
}

func (f Filters) Valid() (err error) {
	if len(f.Keyword) < 3 {
		err = errValidKeyword
	}
	return
}

func (up UpdatePassword) Valid() (err error) {
	if up.UserID == 0 {
		err = errValidID
	}
	if up.OldPassword == "" {
		err = errValidOldPassword
	}
	if up.NewPassword == "" {
		err = errValidNewPassword
	}
	return
}
