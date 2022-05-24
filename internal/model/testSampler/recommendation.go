package testSampler

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"time"
)

func (sampler *Sampler) GetRecommendation(customFields map[string]interface{}) (recommendation *model.Recommendation) {
	var media model.Media
	mediaCustomFields, ok := customFields["Media"]
	if ok {
		media, ok = mediaCustomFields.(model.Media)
		require.Truef(sampler.t, ok, "incorrect type for custom fields")
	} else {
		mediaSample := sampler.GetMedia(nil)
		require.NotNil(sampler.t, mediaSample)
		media = *mediaSample
	}

	var sender model.User
	senderCustomFields, ok := customFields["Sender"]
	if ok {
		sender, ok = senderCustomFields.(model.User)
		require.Truef(sampler.t, ok, "incorrect type for custom fields")
	} else {
		senderSample := sampler.GetUser(nil)
		require.NotNil(sampler.t, senderSample)
		sender = *senderSample
	}

	var circles []model.Circle
	circlesCustomFields, ok := customFields["Circles"]
	if ok {
		circles, ok = circlesCustomFields.([]model.Circle)
		require.Truef(sampler.t, ok, "incorrect type for custom fields")
	} else {
		for range fakeData.FakeRange(1, 3) {
			circleSample := sampler.GetCircle(nil)
			require.NotNil(sampler.t, circleSample)
			circles = append(circles, *circleSample)
		}
	}

	recommendation = &model.Recommendation{
		Sender:  sender,
		Circles: circles,
		Media:   media,
		Text:    fake.Sentences(),
		Date:    time.Now(),
	}

	err := fakeData.FillStructWithFields(recommendation, customFields)
	require.NoError(sampler.t, err)
	return
}
