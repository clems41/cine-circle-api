package testSampler

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"strings"
)

func (sampler *Sampler) GetUser(customFields map[string]interface{}) (user *model.User) {
	user = &model.User{
		Username:  strings.ToLower(fakeData.UniqueUsername()),
		LastName:  fake.LastName(),
		FirstName: fake.FirstName(),
		Email:     fakeData.UniqueEmail(),
		Active:    true,
	}
	err := fakeData.FillStructWithFields(user, customFields)
	require.NoError(sampler.t, err)
	return
}
