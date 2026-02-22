package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOutageInfo_Create(t *testing.T) {
	p, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	d := NewOutageDescription("test")
	info := NewOutageInfo(p, d)
	assert.Equal(t, p, info.Period)
	assert.Equal(t, d, info.Description)
}

func TestOutageInfo_EqualsIdentical(t *testing.T) {
	p, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	d := NewOutageDescription("test")
	i1 := NewOutageInfo(p, d)
	i2 := NewOutageInfo(p, d)
	assert.True(t, i1.Equals(i2))
}

func TestOutageInfo_DifferentPeriod(t *testing.T) {
	p1, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	p2, _ := NewOutagePeriod(
		time.Date(2024, 1, 2, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 16, 0, 0, 0, time.UTC),
	)
	d := NewOutageDescription("test")
	i1 := NewOutageInfo(p1, d)
	i2 := NewOutageInfo(p2, d)
	assert.False(t, i1.Equals(i2))
}

func TestOutageInfo_DifferentDescription(t *testing.T) {
	p, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	i1 := NewOutageInfo(p, NewOutageDescription("test1"))
	i2 := NewOutageInfo(p, NewOutageDescription("test2"))
	assert.False(t, i1.Equals(i2))
}

func TestOutageInfo_BothDifferent(t *testing.T) {
	p1, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	p2, _ := NewOutagePeriod(
		time.Date(2024, 1, 2, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 16, 0, 0, 0, time.UTC),
	)
	i1 := NewOutageInfo(p1, NewOutageDescription("test1"))
	i2 := NewOutageInfo(p2, NewOutageDescription("test2"))
	assert.False(t, i1.Equals(i2))
}
