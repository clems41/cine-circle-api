package fakeData

import (
	"github.com/oleiade/reflections"
)

// FillStructWithFields will fill value fields using map. customFields map using field name as key and field value as value.
func FillStructWithFields(value interface{}, customFields map[string]interface{}) (err error) {
	for fieldName, fieldValue := range customFields {
		err = reflections.SetField(value, fieldName, fieldValue)
		if err != nil {
			return
		}
	}
	return
}
