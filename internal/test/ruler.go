package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
	"time"
)

type Ruler struct {
	t *testing.T
}

func NewRuler(t *testing.T) *Ruler {
	return &Ruler{t: t}
}

func formatFieldName(callStrings []string, suffixes ...string) string {
	out := strings.Join(callStrings, ".")
	suffixesString := strings.Join(suffixes, ".")
	if len(out) > 0 && len(suffixesString) > 0 {
		out += "."
	}
	out += suffixesString

	return out
}

type IgnoredField struct {
}

type NotEmptyField struct {
}

type EmptyField struct {
}

type ZeroValueField struct {
}

type SameValueField struct {
	Expected interface{}
}

func (ruler *Ruler) CheckStruct(actual interface{}, expected map[string]interface{}, callStrings ...string) {
	t := ruler.t
	actualType := reflect.TypeOf(actual)
	actualValue := reflect.ValueOf(actual)

	// This utility is for testing structs !
	if actualType.Kind() != reflect.Struct {
		assert.Fail(ruler.t, "Non struct", "checkStruct was used with a non-struct object, please fix your test for %s", formatFieldName(callStrings))
		return
	}

	// We browse all expected fields and test if they correspond to the actual struct or not
	for fieldKey, expectedFieldInterface := range expected {
		_, structFieldFound := actualType.FieldByName(fieldKey)
		if !structFieldFound {
			assert.Fail(t, "Field not found", "Expected field not found : %s", formatFieldName(callStrings, fieldKey))
			continue
		}

		actualFieldValue := actualValue.FieldByName(fieldKey)
		switch castedExpectedField := expectedFieldInterface.(type) {
		case []map[string]interface{}:
			// A list of maps should correspond to a list of structs
			if actualFieldValue.Kind() != reflect.Slice {
				assert.Fail(ruler.t, "No slice found", "A slice was expected for field %s. We found : %s", formatFieldName(callStrings, fieldKey), actualFieldValue.Interface())
				continue
			}

			// The two slices lengths must match
			if actualFieldValue.Len() != len(castedExpectedField) {
				assert.Fail(ruler.t, "Wrong slice size", "For field %s, The actual slice has a len of %d, while a size of %d was expected", formatFieldName(callStrings, fieldKey), actualFieldValue.Len(), len(castedExpectedField))
				continue
			}

			// All elements should match : recursion !
			for itemIndex, expectedItem := range castedExpectedField {
				itemValue := actualFieldValue.Index(itemIndex)
				ruler.CheckStruct(itemValue.Interface(), expectedItem, append(callStrings, fieldKey, fmt.Sprintf("[%d]", itemIndex))...)
			}

		case map[string]interface{}:
			ruler.CheckStruct(actualFieldValue.Interface(), castedExpectedField, append(callStrings, fieldKey)...)
		case time.Time:
			checkStructTime(actualFieldValue, ruler, castedExpectedField, callStrings, fieldKey)
		case *time.Time:
			checkStructTime(actualFieldValue.Elem(), ruler, *castedExpectedField, callStrings, fieldKey)
		case IgnoredField, *IgnoredField:
			continue
		case NotEmptyField, *NotEmptyField:
			assert.NotEmptyf(t, actualFieldValue.Interface(), "Check that field %s is not empty", formatFieldName(callStrings, fieldKey))
		case EmptyField, *EmptyField:
			assert.Emptyf(t, actualFieldValue.Interface(), "Check that field %s is empty", formatFieldName(callStrings, fieldKey))
		case SameValueField:
			assert.EqualValues(t, expectedFieldInterface.(SameValueField).Expected, actualFieldValue.Interface(), "Check of Same Value for field %s", formatFieldName(callStrings, fieldKey))
		case ZeroValueField:
			assert.Zerof(t, actualFieldValue.Interface(), "Field '%s' should have zero value but has value %v", formatFieldName(callStrings, fieldKey), actualFieldValue.Interface())
		default:
			assert.Equalf(t, expectedFieldInterface, actualFieldValue.Interface(), "Check of Equality for field %s", formatFieldName(callStrings, fieldKey))
		}
		continue
	}

	// If some fields were not reviewed, we fail and tell what fields were missing
	allFieldsReviewed := len(expected) == actualType.NumField()
	if !allFieldsReviewed {
		missingFields := make([]string, 0)
		for i := 0; i < actualType.NumField(); i++ {
			name := actualType.Field(i).Name
			_, found := expected[name]
			if !found {
				missingFields = append(missingFields, name)
			}
		}

		fieldName := formatFieldName(callStrings)
		precision := ""
		if fieldName != "" {
			precision = fmt.Sprintf("of '%s'", fieldName)
		}

		assert.Truef(t, allFieldsReviewed, "Some fields %sare missing : %s", precision, missingFields)

	}
}

// areSameTimes checks if two times are similar.
func (ruler *Ruler) areSameTimes(time1 time.Time, time2 time.Time) bool {
	return time1.Local().Truncate(time.Minute) == time2.Local().Truncate(time.Minute)
}

func checkStructTime(actualFieldValue reflect.Value, ruler *Ruler, expectedTime time.Time, callStrings []string, fieldKey string) {
	trueTime, ok := actualFieldValue.Interface().(time.Time)
	if ok {
		assert.True(ruler.t, ruler.areSameTimes(expectedTime, trueTime), "Check of Equality for field %s, expected %v, got %v", formatFieldName(callStrings, fieldKey), expectedTime, trueTime)
	} else {
		trueTimePtr, ok2 := actualFieldValue.Interface().(*time.Time)
		if ok2 {
			assert.True(ruler.t, ruler.areSameTimes(expectedTime, *trueTimePtr), "Check of Equality for field %s, expected %v, got %v", formatFieldName(callStrings, fieldKey), expectedTime, trueTime)
		} else {
			assert.Fail(ruler.t, "Expected Time", "Expected time value for field %s, got type %T with value : %s", formatFieldName(callStrings, fieldKey), actualFieldValue.Interface(), actualFieldValue.Interface())
		}
	}
}
