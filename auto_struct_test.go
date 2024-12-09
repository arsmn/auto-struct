package autostruct

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type Name struct {
	FN **string `default:"FN"`
	LN *string  `default:"LN"`
}

type Age struct {
	Day   *uint8    `default:"10"`
	Month **int64   `default:"20"`
	Year  ***uint16 `default:"3030"`
}

type Example struct {
	Foo    **bool         `default:"true"`
	Bar    string         `default:"bar"`
	Qux    int8           `default:"123"`
	Name   *Name          `default:"value(struct)"`
	Age    **Age          `default:"value(struct)"`
	Arr1   [5]string      `default:"json([\"1\", \"2\", \"3\", \"4\"])"`
	Arr2   [2][]string    `default:"json([[\"1\", \"2\", \"3\", \"4\"], [\"5\", \"6\", \"7\", \"8\"]])"`
	Arr3   [3]*Name       `default:"repeat(struct)"`
	Arr4   [4]int         `default:"repeat(1)"`
	Slice1 []string       `default:"len(5),cap(10),repeat(1)"`
	Slice2 []string       `default:"json([\"1\", \"2\", \"3\", \"4\", \"5\"])"`
	Slice3 []*Name        `default:"len(1),repeat(struct)"`
	Slice4 [][2]string    `default:"json([[\"1\", \"2\"], [\"3\", \"4\"]])"`
	Map1   map[string]int `default:"json({\"1\": 1})"`
	Dur1   time.Duration  `default:"3s"`
	Dur2   *time.Duration `default:"5h30m15s"`
	Time1  time.Time      `default:"2024-12-09T02:20:35Z"`
	Time2  *time.Time     `default:"value(2024-12-09 02:20:35),layout(DateTime)"`
}

func Test_New(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		act := New[**Example]()

		foo := true
		foo1 := &foo

		fn := "FN"
		fn1 := &fn

		ln := "LN"
		ln1 := &ln

		day := uint8(10)
		month := int64(20)
		month1 := &month
		year := uint16(3030)
		year1 := &year
		year2 := &year1

		age := &Age{
			Day:   &day,
			Month: &month1,
			Year:  &year2,
		}

		dur := time.Duration((5 * time.Hour) + (30 * time.Minute) + (15 * time.Second))
		time_ := time.Date(2024, 12, 9, 2, 20, 35, 0, time.UTC)

		exp := &Example{
			Foo: &foo1,
			Bar: "bar",
			Qux: 123,
			Name: &Name{
				FN: &fn1,
				LN: ln1,
			},
			Age:  &age,
			Arr1: [5]string{"1", "2", "3", "4", ""},
			Arr2: [2][]string{
				{"1", "2", "3", "4"},
				{"5", "6", "7", "8"},
			},
			Arr3: [3]*Name{
				{
					FN: &fn1,
					LN: ln1,
				},
				{
					FN: &fn1,
					LN: ln1,
				},
				{
					FN: &fn1,
					LN: ln1,
				},
			},
			Arr4:   [4]int{1, 1, 1, 1},
			Slice1: []string{"1", "1", "1", "1", "1"},
			Slice2: []string{"1", "2", "3", "4", "5"},
			Slice3: []*Name{
				{
					FN: &fn1,
					LN: ln1,
				},
			},
			Slice4: [][2]string{
				{"1", "2"},
				{"3", "4"},
			},
			Map1: map[string]int{
				"1": 1,
			},
			Dur1:  time.Second * 3,
			Dur2:  &dur,
			Time1: time.Date(2024, 12, 9, 2, 20, 35, 0, time.UTC),
			Time2: &time_,
		}

		if !cmp.Equal(&exp, act) {
			t.Error(cmp.Diff(&exp, act))
		}
	})
}
