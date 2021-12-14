package swagger

import "github.com/go-openapi/spec"

// EnrichSwaggerObject will create understandable method for swagger package based on Info fields to enrich swagger.
func EnrichSwaggerObject(info Info) func(swo *spec.Swagger) {
	return func(swo *spec.Swagger) {
		swo.Info = &spec.Info{
			InfoProps: spec.InfoProps{
				Title:       info.Title,
				Description: info.Description,
				Contact: &spec.ContactInfo{
					ContactInfoProps: spec.ContactInfoProps{
						Name:  info.ContactName,
						Email: info.ContactEmail,
					},
				},
				License: &spec.License{
					LicenseProps: spec.LicenseProps{
						Name: licenseName,
						URL:  licenseUrl,
					},
				},
				Version: info.Version,
			},
		}
		var tags []spec.Tag
		for tagName, tagDescription := range info.Tags {
			tags = append(tags, spec.Tag{
				TagProps: spec.TagProps{
					Name:        string(tagName),
					Description: string(tagDescription),
				},
			})
		}
		swo.Tags = tags
	}
}
