package test

import (
	"fmt"
	"github.com/icrowley/fake"
	"math/rand"
	"time"
)

func FakeTime() time.Time {
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
	return int(rand.Int63n(max - min) + min)
}

func FakeBool() bool {
	return rand.Int63n(1) == 1
}

func FakePassword() string {
	return fake.Password(PasswordMinCharacter, PasswordMaxCharacter, PasswordAllowUpper,
		PasswordAllowNumber, PasswordAllowSpecial)
}
