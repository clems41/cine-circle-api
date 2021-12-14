package testSampler

import (
	"cine-circle-api/internal/repository/model"
	"github.com/stretchr/testify/require"
)

// GetExemple retourne un exemple créé en DB
func (sampler *Sampler) GetExemple() (exemple *model.Exemple) {
	exemple = &model.Exemple{
		// TODO add your custom fields here
	}
	err := sampler.DB.
		Create(exemple).
		Error
	require.NoError(sampler.t, err)
	return
}
