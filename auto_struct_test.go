package autostruct

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type Name struct {
	FN **string `auto:"FN"`
	LN *string  `auto:"LN"`
}

type Age struct {
	Day   *uint8    `auto:"10"`
	Month **int64   `auto:"20"`
	Year  ***uint16 `auto:"3030"`
}

type Example struct {
	Foo             **bool          `auto:"true"`
	Bar             string          `auto:"bar"`
	Qux             int8            `auto:"123"`
	Name            *Name           `auto:"value(struct)"`
	Age             **Age           `auto:"value(struct)"`
	Arr1            [5]string       `auto:"json([\"1\", \"2\", \"3\", \"4\"])"`
	Arr2            [2][]string     `auto:"json([[\"1\", \"2\", \"3\", \"4\"], [\"5\", \"6\", \"7\", \"8\"]])"`
	Arr3            [3]*Name        `auto:"repeat(struct)"`
	Arr4            [4]int          `auto:"repeat(1)"`
	Slice1          []string        `auto:"len(5),cap(10),repeat(1)"`
	Slice2          []string        `auto:"json([\"1\", \"2\", \"3\", \"4\", \"5\"])"`
	Slice3          []*Name         `auto:"len(1),repeat(struct)"`
	Slice4          [][2]string     `auto:"json([[\"1\", \"2\"], [\"3\", \"4\"]])"`
	Map1            map[string]int  `auto:"json({\"1\": 1})"`
	Map2            map[string]int  `auto:"len(5),value(key1:1,key2:2,key3:3)"`
	Dur1            time.Duration   `auto:"3s"`
	Dur2            *time.Duration  `auto:"5h30m15s"`
	Time1           time.Time       `auto:"2024-12-09T02:20:35Z"`
	Time2           *time.Time      `auto:"value(2024-12-09 02:20:35),layout(DateTime)"`
	Rune1           rune            `auto:"rune(1)"`
	Rune2           rune            `auto:"rune(a)"`
	Runes1          []rune          `auto:"rune(abc)"`
	Byte1           byte            `auto:"byte(1)"`
	Byte2           byte            `auto:"byte(a)"`
	Bytes1          []byte          `auto:"byte(abc)"`
	JSONRawMessage1 json.RawMessage `auto:"json({\"key\": \"value\"})"`
	JSONRawMessage2 json.RawMessage `auto:"{\"key\": \"value\"}"`
	Chan1           chan int        `auto:"chan"`
	Chan2           chan<- int      `auto:"chan,cap(5)"`
	Chan3           <-chan int      `auto:"chan,cap(5)"`
	Interface1      any             `auto:"123"`
	Interface2      any             `auto:"true"`
	Interface3      any             `auto:"\"abc\""`
	Interface4      any             `auto:"{\"key\": \"value\"}"`
	Interface5      any             `auto:"[\"1\", \"2\", \"3\"]"`
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
			Map2: map[string]int{
				"key1": 1,
				"key2": 2,
				"key3": 3,
			},
			Dur1:            time.Second * 3,
			Dur2:            &dur,
			Time1:           time.Date(2024, 12, 9, 2, 20, 35, 0, time.UTC),
			Time2:           &time_,
			Rune1:           '1',
			Rune2:           'a',
			Runes1:          []rune{'a', 'b', 'c'},
			Byte1:           byte('1'),
			Byte2:           byte('a'),
			Bytes1:          []byte("abc"),
			JSONRawMessage1: json.RawMessage(`{"key": "value"}`),
			JSONRawMessage2: json.RawMessage(`{"key": "value"}`),
			Chan1:           make(chan int),
			Chan2:           make(chan<- int, 5),
			Chan3:           make(<-chan int, 5),
			Interface1:      float64(123),
			Interface2:      true,
			Interface3:      "abc",
			Interface4:      map[string]any{"key": "value"},
			Interface5:      []any{"1", "2", "3"},
		}

		if !cmp.Equal(&exp, act, cmp.FilterPath(
			func(p cmp.Path) bool { return strings.HasPrefix(p.String(), "Chan") },
			cmp.Comparer(func(a, b any) bool {
				return reflect.ValueOf(a).Cap() == reflect.ValueOf(b).Cap() &&
					reflect.TypeOf(a).ChanDir() == reflect.TypeOf(b).ChanDir()
			}),
		)) {
			t.Error(cmp.Diff(&exp, act))
		}
	})
}
