package autostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	Foo  **bool      `default:"true"`
	Bar  string      `default:"33"`
	Qux  int8        `default:"33"`
	Name *Name       `default:"nested"`
	Age  **Age       `default:"nested"`
	Arr  [5]string   `default:"[\"1\", \"2\", \"3\", \"4\"]"`
	Arr2 [2][]string `default:"[[\"1\", \"2\", \"3\", \"4\"], [\"5\", \"6\", \"7\", \"8\"]]"`
	Arr3 [3]*Name    `default:"nested"`
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

		exp := &Example{
			Foo: &foo1,
			Bar: "33",
			Qux: 33,
			Name: &Name{
				FN: &fn1,
				LN: ln1,
			},
			Age: &age,
			Arr: [5]string{"1", "2", "3", "4", ""},
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
		}

		assert.Equal(t, &exp, act)
	})
}
