package testSampler

import (
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
)

func (sampler *Sampler) GetCircle() (circle *model.Circle) {
	circle = &model.Circle{
		Name:        fake.Title(),
		Description: fake.Sentences(),
	}

	err := sampler.DB.Create(circle).Error
	require.NoError(sampler.t, err)
	return
}

func (sampler *Sampler) GetCircleWithName(name string) (circle *model.Circle) {
	circle = &model.Circle{
		Name:        name,
		Description: fake.Sentences(),
	}

	err := sampler.DB.Create(circle).Error
	require.NoError(sampler.t, err)
	return
}

func (sampler *Sampler) GetCircleWithUsers() (circle *model.Circle) {
	circle = &model.Circle{
		Name:        fake.Title(),
		Description: fake.Sentences(),
	}

	for range fakeData.FakeRange(4, 10) {
		user := sampler.GetUser()
		require.NotNil(sampler.t, user)
		circle.Users = append(circle.Users, *user)
	}

	err := sampler.DB.Create(circle).Error
	require.NoError(sampler.t, err)
	return
}
