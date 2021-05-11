package circleDom

/*var _ Repository = (*circleRepository)(nil)

type circleRepository struct {
	DB *gorm.DB
}

func NewCircleRepository(DB *gorm.DB) *circleRepository {
	return &circleRepository{DB: DB}
}

func (r circleRepository) Migrate() {

	err := r.DB.AutoMigrate(&Circle{})
	if err != nil {
		logger.Sugar.Fatalf("Error occurs when migrating circleRepository : %s", err.Error())
	}

	err = r.DB.
		Exec("CREATE INDEX IF NOT EXISTS idx_circle_user_circle ON circle_user (circle_id)").
		Error
	if err != nil {
		logger.Sugar.Fatalf("Error while creating index : %s", err.Error())
	}

	err = r.DB.
		Exec("CREATE INDEX IF NOT EXISTS idx_circle_user_user ON circle_user (user_id)").
		Error
	if err != nil {
		logger.Sugar.Fatalf("Error while creating index : %s", err.Error())
	}

}

func (r circleRepository) Create(creation Creation) (result Result, err error) {
	circle := Circle{
		Name:        creation.Name,
		Description: creation.Description,
	}

	var users []userDom.User
	err = r.DB.
		Find(&users, "id IN (?)", creation.UsersID).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedError(err)
	}
	if len(users) != len(creation.UsersID) {
		return result, typedErrors.NewRepositoryResourceNotFoundErrorf("Not all users has been found, got only %d", len(users))
	}

	circle.Users = users
	err = r.DB.
		Create(&circle).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedError(err)
	}

	result = r.toResult(circle)
	return
}

func (r circleRepository) Update(update Update) (result Result, err error){
	var circle Circle
	err = r.DB.
		Take(&circle, "id = ?", update.CircleID).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedError(err)
	}

	var users []userDom.User
	err = r.DB.
		Find(&users, "id IN (?)", update.UsersID).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedError(err)
	}
	if len(users) != len(update.UsersID) {
		return result, typedErrors.NewRepositoryResourceNotFoundErrorf("Not all users has been found, got only %d", len(users))
	}

	// Remove previous associations
	err = r.DB.
		Model(&circle).
		Association("Users").
		Clear()
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedError(err)
	}

	err = r.DB.
		Model(&circle).
		Updates(Circle{Name: update.Name, Description: update.Description, Users: users}).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedError(err)
	}

	result = r.toResult(circle)
	return

}

func (r circleRepository) Delete(delete Delete) (err error){
	var circle Circle
	err = r.DB.
		Take(&circle, "id = ?", delete.CircleID).
		Error
	if err != nil {
		return typedErrors.NewRepositoryQueryFailedError(err)
	}
	err = r.DB.
		Delete(&circle).
		Error
	return
}

func (r circleRepository) Get(get Get) (result Result,err error){
	var circle Circle
	err = r.DB.
		Preload("Users").
		Take(&circle, "id = ?", get.CircleID).
		Error
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedError(err)
	}
	result = r.toResult(circle)
	return
}

func (r circleRepository) toResult(circle Circle) (result Result){
	result = Result{
		CircleID:    circle.GetID(),
		Name:        circle.Name,
		Description: circle.Description,
		Users:       nil,
	}

	for _, user := range circle.Users {
		result.Users = append(result.Users, userDom.Result{
			UserID:      user.GetID(),
			Username:    user.Username,
			DisplayName: user.DisplayName,
		})
	}
	return
}*/
