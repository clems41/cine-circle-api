package testSampler

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
)

func (sampler *Sampler) GetCircle(customFields map[string]interface{}) (circle *model.Circle) {
	var users []model.User
	usersCustomFields, ok := customFields["Users"]
	if ok {
		users, ok = usersCustomFields.([]model.User)
		require.Truef(sampler.t, ok, "incorrect type for custom fields")
	} else {
		for range fakeData.FakeRange(4, 10) {
			userSample := sampler.GetUser(nil)
			require.NotNil(sampler.t, userSample)
			users = append(users, *userSample)
		}
	}

	circle = &model.Circle{
		Users:       users,
		Name:        fake.Title(),
		Description: fake.Sentences(),
	}
	err := fakeData.FillStructWithFields(circle, customFields)
	require.NoError(sampler.t, err)
	return
}
