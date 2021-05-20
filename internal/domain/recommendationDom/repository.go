package recommendationDom

import (
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ repository = (*Repository)(nil)

type repository interface {
	Create(recommendation *repositoryModel.Recommendation) (err error)
	GetUserIDsFromCircle(circleID uint) (userIDs []uint, err error)
	GetUserIDsCloseToUser(userID uint) (userIDs []uint, err error)
	CheckIfMovieExists(movieID uint) (exists bool, err error)
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Migrate() {

	err := r.DB.AutoMigrate(&repositoryModel.Recommendation{})
	if err != nil {
		logger.Sugar.Fatalf("Error occurs when migrating movieRepository : %s", err.Error())
	}

	err = r.DB.
		Exec("CREATE INDEX IF NOT EXISTS idx_recommendation_circle_recommendation_id ON recommendation_circle (recommendation_id)").
		Error
	if err != nil {
		logger.Sugar.Fatalf("Error while creating index : %s", err.Error())
	}

	err = r.DB.
		Exec("CREATE INDEX IF NOT EXISTS idx_recommendation_circle_circle_id ON recommendation_circle (circle_id)").
		Error
	if err != nil {
		logger.Sugar.Fatalf("Error while creating index : %s", err.Error())
	}

	err = r.DB.
		Exec("CREATE INDEX IF NOT EXISTS idx_recommendation_user_recommendation_id ON recommendation_user (recommendation_id)").
		Error
	if err != nil {
		logger.Sugar.Fatalf("Error while creating index : %s", err.Error())
	}

	err = r.DB.
		Exec("CREATE INDEX IF NOT EXISTS idx_recommendation_user_user_id ON recommendation_user (user_id)").
		Error
	if err != nil {
		logger.Sugar.Fatalf("Error while creating index : %s", err.Error())
	}

}

func (r *Repository) Create(recommendation *repositoryModel.Recommendation) (err error) {

	return r.DB.Transaction(func(tx *gorm.DB) error {

		circles := recommendation.Circles
		users := recommendation.Users
		recommendation.Circles = nil
		recommendation.Users = nil

		err = tx.Create(&recommendation).Error
		if err != nil {
			return errors.WithStack(err)
		}

		for _, user := range users {
			err = tx.
				Exec("INSERT INTO recommendation_user (recommendation_id,user_id) VALUES (?,?) ON CONFLICT DO NOTHING", recommendation.GetID(), user.GetID()).
				Error
			if err != nil {
				return errors.WithStack(err)
			}
		}

		for _, circle := range circles {
			err = tx.
				Exec("INSERT INTO recommendation_circle (recommendation_id,circle_id) VALUES (?,?) ON CONFLICT DO NOTHING", recommendation.GetID(), circle.GetID()).
				Error
			if err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})
}

func (r *Repository) GetUserIDsFromCircle(circleID uint) (userIDs []uint, err error) {
	var users []repositoryModel.User
	err = r.DB.
		Table("users").
		Joins("INNER JOIN circle_user ON circle_user.user_id = users.id and circle_user.circle_id = ?", circleID).
		Select("users.id").
		Scan(&users).
		Error
	if err != nil {
		return userIDs, errors.WithStack(err)
	}

	for _, user := range users {
		userIDs = append(userIDs, user.GetID())
	}

	return
}

// GetUserIDsCloseToUser list id of all users that are at least in one circle with specific user
func (r *Repository) GetUserIDsCloseToUser(userID uint) (userIDs []uint, err error) {
	var users []repositoryModel.User
	err = r.DB.
		Table("users").
		Joins(`INNER JOIN circle_user ON circle_user.user_id = users.id 
			and circle_user.circle_id IN (select circle_id from circle_user where user_id = ?)`, userID).
		Select("users.id").
		Scan(&users).
		Error
	if err != nil {
		return userIDs, errors.WithStack(err)
	}

	for _, user := range users {
		userIDs = append(userIDs, user.GetID())
	}

	return
}

func (r *Repository) CheckIfMovieExists(movieID uint) (exists bool, err error) {
	var movie repositoryModel.Movie
	err = r.DB.
		Take(&movie, "id = ?", movieID).
		Error
	if err != nil {
		return false, errors.WithStack(err)
	}

	return movie.ID == movieID, nil
}
