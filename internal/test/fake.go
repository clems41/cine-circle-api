package test

import (
	"fmt"
	"github.com/icrowley/fake"
	"math/rand"
	"reflect"
	"time"
)

func FakeTime() time.Time {
	rand.Seed(time.Now().UnixNano())
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func FakeTimePtr() *time.Time {
	fakeTime := FakeTime()
	return &fakeTime
}

func FakeTimeBefore(max time.Time) time.Time {
	rand.Seed(time.Now().UnixNano())
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max.Unix() - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func FakeTimePtrBefore(max time.Time) *time.Time {
	fakeTime := FakeTimeBefore(max)
	return &fakeTime
}

func FakeTimeAfter(min time.Time) time.Time {
	rand.Seed(time.Now().UnixNano())
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min.Unix()

	sec := rand.Int63n(delta) + min.Unix()
	return time.Unix(sec, 0)
}

func FakeTimePtrAfter(min time.Time) *time.Time {
	fakeTime := FakeTimeAfter(min)
	return &fakeTime
}

func FakeLocalPhone() string {
	return fmt.Sprintf("%d.%d.%d", rand.Int31n(89) + 10, rand.Int31n(89) + 10, rand.Int31n(89) + 10)
}

func FakeIntBetween(min, max int64) int {
	rand.Seed(time.Now().UnixNano())
	return int(rand.Int63n(max - min) + min)
}

func FakeRange(min, max int64) []int {
	rangeValue := FakeIntBetween(min, max)
	return make([]int, rangeValue)
}

func FakeBool() bool {
	return rand.Int63n(1) == 1
}

func FakePassword() string {
	return fake.Password(PasswordMinCharacter, PasswordMaxCharacter, PasswordAllowUpper,
		PasswordAllowNumber, PasswordAllowSpecial)
}

func RandomElement(slice interface{}) interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil
	}
	idx := FakeIntBetween(0, int64(s.Len() - 1))
	return s.Index(idx).Interface()
}
