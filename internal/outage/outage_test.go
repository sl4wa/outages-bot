package outage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	ot0 = time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	ot1 = time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC)
)

func makeTestOutage(streetID int, streetName string, buildings []string, start, end time.Time, comment string) *Outage {
	period, _ := NewPeriod(start, end)
	addr, _ := NewAddress(streetID, streetName, buildings, "Львів")
	return &Outage{
		Period:      period,
		Address:     addr,
		Description: NewDescription(comment),
	}
}

func TestOutage_Create(t *testing.T) {
	period, _ := NewPeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	addr, _ := NewAddress(1, "Стрийська", []string{"10", "12"}, "Львів")
	desc := NewDescription("test")
	o := &Outage{ID: 42, Period: period, Address: addr, Description: desc}
	assert.Equal(t, 42, o.ID)
	assert.Equal(t, period, o.Period)
	assert.Equal(t, addr, o.Address)
	assert.Equal(t, desc, o.Description)
}

func TestOutagesEqual_SameOutages_True(t *testing.T) {
	a := []*Outage{makeTestOutage(1, "Стрийська", []string{"10"}, ot0, ot1, "c")}
	b := []*Outage{makeTestOutage(1, "Стрийська", []string{"10"}, ot0, ot1, "c")}
	assert.True(t, OutagesEqual(a, b))
}

func TestOutagesEqual_BothNil_True(t *testing.T) {
	assert.True(t, OutagesEqual(nil, nil))
}

func TestOutagesEqual_EmptySlices_True(t *testing.T) {
	assert.True(t, OutagesEqual([]*Outage{}, []*Outage{}))
}

func TestOutagesEqual_NilVsEmpty_True(t *testing.T) {
	// len(nil) == len([]) == 0, so positional comparison considers them equal
	assert.True(t, OutagesEqual(nil, []*Outage{}))
}

func TestOutagesEqual_DifferentComment_False(t *testing.T) {
	a := []*Outage{makeTestOutage(1, "Стрийська", []string{"10"}, ot0, ot1, "old")}
	b := []*Outage{makeTestOutage(1, "Стрийська", []string{"10"}, ot0, ot1, "new")}
	assert.False(t, OutagesEqual(a, b))
}

func TestOutagesEqual_DifferentPeriod_False(t *testing.T) {
	t2 := time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC)
	a := []*Outage{makeTestOutage(1, "Стрийська", []string{"10"}, ot0, ot1, "c")}
	b := []*Outage{makeTestOutage(1, "Стрийська", []string{"10"}, ot0, t2, "c")}
	assert.False(t, OutagesEqual(a, b))
}

func TestOutagesEqual_DifferentCity_False(t *testing.T) {
	period, _ := NewPeriod(ot0, ot1)
	addrA, _ := NewAddress(1, "Стрийська", []string{"10"}, "Львів")
	addrB, _ := NewAddress(1, "Стрийська", []string{"10"}, "Київ")
	a := []*Outage{{Period: period, Address: addrA, Description: NewDescription("c")}}
	b := []*Outage{{Period: period, Address: addrB, Description: NewDescription("c")}}
	assert.False(t, OutagesEqual(a, b))
}

func TestOutagesEqual_DifferentStreet_False(t *testing.T) {
	a := []*Outage{makeTestOutage(1, "Стрийська", []string{"10"}, ot0, ot1, "c")}
	b := []*Outage{makeTestOutage(2, "Наукова", []string{"10"}, ot0, ot1, "c")}
	assert.False(t, OutagesEqual(a, b))
}

func TestOutagesEqual_DifferentBuildings_False(t *testing.T) {
	a := []*Outage{makeTestOutage(1, "Стрийська", []string{"10"}, ot0, ot1, "c")}
	b := []*Outage{makeTestOutage(1, "Стрийська", []string{"10", "12"}, ot0, ot1, "c")}
	assert.False(t, OutagesEqual(a, b))
}

func TestOutagesEqual_ReorderedOutages_False(t *testing.T) {
	o1 := makeTestOutage(1, "Стрийська", []string{"10"}, ot0, ot1, "alpha")
	o2 := makeTestOutage(2, "Наукова", []string{"5"}, ot0, ot1, "beta")
	assert.False(t, OutagesEqual([]*Outage{o1, o2}, []*Outage{o2, o1}))
}

func TestOutagesEqual_SameInstantDifferentTimezone_True(t *testing.T) {
	kyiv := time.FixedZone("UTC+3", 3*60*60)
	a := []*Outage{makeTestOutage(1, "Стрийська", []string{"10"}, ot0, ot1, "c")}
	b := []*Outage{makeTestOutage(1, "Стрийська", []string{"10"}, ot0.In(kyiv), ot1.In(kyiv), "c")}
	assert.True(t, OutagesEqual(a, b))
}

func TestOutagesEqual_ReorderedBuildings_False(t *testing.T) {
	a := []*Outage{makeTestOutage(1, "Стрийська", []string{"12", "10", "14"}, ot0, ot1, "c")}
	b := []*Outage{makeTestOutage(1, "Стрийська", []string{"10", "14", "12"}, ot0, ot1, "c")}
	assert.False(t, OutagesEqual(a, b))
}
