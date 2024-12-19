package autostruct

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type Basic struct {
	Bool1   bool   `auto:"true"`
	Bool2   bool   `auto:"false"`
	Bool3   bool   `auto:"t"`
	Bool4   bool   `auto:"0"`
	String1 string `auto:"abc"`
	String2 string `auto:"123"`
}

type Signed struct {
	Int   *int       `auto:"0"`
	Int8  **int8     `auto:"8"`
	Int16 ***int16   `auto:"16"`
	Int32 ****int32  `auto:"32"`
	Int64 *****int64 `auto:"64"`
}

type Unsigned struct {
	Uint   *uint       `auto:"0"`
	Uint8  **uint8     `auto:"8"`
	Uint16 ***uint16   `auto:"16"`
	Uint32 ****uint32  `auto:"32"`
	Uint64 *****uint64 `auto:"64"`
}

type Float struct {
	Float1 float32 `auto:"1.2345"`
	Float2 float64 `auto:"1.23456789"`
}

type Complex struct {
	Complex1 complex64  `auto:"1+2i"`
	Complex2 complex128 `auto:"3+4i"`
}

type Test struct {
	Struct1         *Basic          `auto:"struct"`
	Struct2         **Signed        `auto:"value(struct)"`
	Struct3         ***Unsigned     `auto:"value(struct)"`
	Struct4         ****Float       `auto:"value(struct)"`
	Struct5         *****Complex    `auto:"value(struct)"`
	Arr1            [5]string       `auto:"json([\"1\", \"2\", \"3\", \"4\"])"`
	Arr2            [2][]string     `auto:"json([[\"1\", \"2\", \"3\", \"4\"], [\"5\", \"6\", \"7\", \"8\"]])"`
	Arr3            [3]*Basic       `auto:"repeat(struct)"`
	Arr4            [4]int          `auto:"repeat(1)"`
	Slice1          []string        `auto:"len(5),cap(10),repeat(1)"`
	Slice2          []string        `auto:"json([\"1\", \"2\", \"3\", \"4\", \"5\"])"`
	Slice3          []*Basic        `auto:"len(1),repeat(struct)"`
	Slice4          [][2]string     `auto:"json([[\"1\", \"2\"], [\"3\", \"4\"]])"`
	Map1            map[string]int  `auto:"json({\"1\": 1})"`
	Map2            map[string]int  `auto:"len(5),value(key1:1,key2:2,key3:3)"`
	Duration1       time.Duration   `auto:"3s"`
	Duration2       *time.Duration  `auto:"value(5h30m15s)"`
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
	Chan2           chan<- int      `auto:"chan(5)"`
	Chan3           <-chan int      `auto:"chan(10)"`
	Interface1      any             `auto:"123"`
	Interface2      any             `auto:"true"`
	Interface3      any             `auto:"\"abc\""`
	Interface4      any             `auto:"{\"key\": \"value\"}"`
	Interface5      any             `auto:"[\"1\", \"2\", \"3\"]"`
}

func Test_New(t *testing.T) {
	toPointer := func(value any, depth int) any {
		if depth <= 0 {
			return value
		}

		current := reflect.ValueOf(value)
		for i := 0; i < depth; i++ {
			ptr := reflect.New(current.Type())
			ptr.Elem().Set(current)
			current = ptr
		}

		return current.Interface()
	}

	cmpOptions := func() []cmp.Option {
		return []cmp.Option{
			cmp.FilterPath(
				func(p cmp.Path) bool { return strings.HasPrefix(p.String(), "Chan") },
				cmp.Comparer(func(a, b any) bool {
					return reflect.ValueOf(a).Cap() == reflect.ValueOf(b).Cap() &&
						reflect.TypeOf(a).ChanDir() == reflect.TypeOf(b).ChanDir()
				}),
			),
		}
	}

	t.Run("success", func(t *testing.T) {
		act := New[**Test]()

		exp := &Test{
			Struct1: toPointer(Basic{
				Bool1:   true,
				Bool2:   false,
				Bool3:   true,
				Bool4:   false,
				String1: "abc",
				String2: "123",
			}, 1).(*Basic),
			Struct2: toPointer(Signed{
				Int:   toPointer(int(0), 1).(*int),
				Int8:  toPointer(int8(8), 2).(**int8),
				Int16: toPointer(int16(16), 3).(***int16),
				Int32: toPointer(int32(32), 4).(****int32),
				Int64: toPointer(int64(64), 5).(*****int64),
			}, 2).(**Signed),
			Struct3: toPointer(Unsigned{
				Uint:   toPointer(uint(0), 1).(*uint),
				Uint8:  toPointer(uint8(8), 2).(**uint8),
				Uint16: toPointer(uint16(16), 3).(***uint16),
				Uint32: toPointer(uint32(32), 4).(****uint32),
				Uint64: toPointer(uint64(64), 5).(*****uint64),
			}, 3).(***Unsigned),
			Struct4: toPointer(Float{
				Float1: float32(1.2345),
				Float2: float64(1.23456789),
			}, 4).(****Float),
			Struct5: toPointer(Complex{
				Complex1: complex64(1 + 2i),
				Complex2: complex128(3 + 4i),
			}, 5).(*****Complex),
			Arr1: [5]string{"1", "2", "3", "4", ""},
			Arr2: [2][]string{{"1", "2", "3", "4"}, {"5", "6", "7", "8"}},
			Arr3: [3]*Basic{
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
			},
			Arr4:   [4]int{1, 1, 1, 1},
			Slice1: []string{"1", "1", "1", "1", "1"},
			Slice2: []string{"1", "2", "3", "4", "5"},
			Slice3: []*Basic{
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
			},
			Slice4:          [][2]string{{"1", "2"}, {"3", "4"}},
			Map1:            map[string]int{"1": 1},
			Map2:            map[string]int{"key1": 1, "key2": 2, "key3": 3},
			Duration1:       time.Second * 3,
			Duration2:       toPointer(((5 * time.Hour) + (30 * time.Minute) + (15 * time.Second)), 1).(*time.Duration),
			Time1:           time.Date(2024, 12, 9, 2, 20, 35, 0, time.UTC),
			Time2:           toPointer(time.Date(2024, 12, 9, 2, 20, 35, 0, time.UTC), 1).(*time.Time),
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
			Chan3:           make(<-chan int, 10),
			Interface1:      float64(123),
			Interface2:      true,
			Interface3:      "abc",
			Interface4:      map[string]any{"key": "value"},
			Interface5:      []any{"1", "2", "3"},
		}

		if !cmp.Equal(&exp, act, cmpOptions()...) {
			t.Error(cmp.Diff(&exp, act, cmpOptions()...))
		}
	})

	t.Run("success-with-cache", func(t *testing.T) {
		cached := NewCache()
		act1 := New[Test](WithCache(cached))
		act2 := New[Test](WithCache(cached))

		exp := Test{
			Struct1: toPointer(Basic{
				Bool1:   true,
				Bool2:   false,
				Bool3:   true,
				Bool4:   false,
				String1: "abc",
				String2: "123",
			}, 1).(*Basic),
			Struct2: toPointer(Signed{
				Int:   toPointer(int(0), 1).(*int),
				Int8:  toPointer(int8(8), 2).(**int8),
				Int16: toPointer(int16(16), 3).(***int16),
				Int32: toPointer(int32(32), 4).(****int32),
				Int64: toPointer(int64(64), 5).(*****int64),
			}, 2).(**Signed),
			Struct3: toPointer(Unsigned{
				Uint:   toPointer(uint(0), 1).(*uint),
				Uint8:  toPointer(uint8(8), 2).(**uint8),
				Uint16: toPointer(uint16(16), 3).(***uint16),
				Uint32: toPointer(uint32(32), 4).(****uint32),
				Uint64: toPointer(uint64(64), 5).(*****uint64),
			}, 3).(***Unsigned),
			Struct4: toPointer(Float{
				Float1: float32(1.2345),
				Float2: float64(1.23456789),
			}, 4).(****Float),
			Struct5: toPointer(Complex{
				Complex1: complex64(1 + 2i),
				Complex2: complex128(3 + 4i),
			}, 5).(*****Complex),
			Arr1: [5]string{"1", "2", "3", "4", ""},
			Arr2: [2][]string{{"1", "2", "3", "4"}, {"5", "6", "7", "8"}},
			Arr3: [3]*Basic{
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
			},
			Arr4:   [4]int{1, 1, 1, 1},
			Slice1: []string{"1", "1", "1", "1", "1"},
			Slice2: []string{"1", "2", "3", "4", "5"},
			Slice3: []*Basic{
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
			},
			Slice4:          [][2]string{{"1", "2"}, {"3", "4"}},
			Map1:            map[string]int{"1": 1},
			Map2:            map[string]int{"key1": 1, "key2": 2, "key3": 3},
			Duration1:       time.Second * 3,
			Duration2:       toPointer(((5 * time.Hour) + (30 * time.Minute) + (15 * time.Second)), 1).(*time.Duration),
			Time1:           time.Date(2024, 12, 9, 2, 20, 35, 0, time.UTC),
			Time2:           toPointer(time.Date(2024, 12, 9, 2, 20, 35, 0, time.UTC), 1).(*time.Time),
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
			Chan3:           make(<-chan int, 10),
			Interface1:      float64(123),
			Interface2:      true,
			Interface3:      "abc",
			Interface4:      map[string]any{"key": "value"},
			Interface5:      []any{"1", "2", "3"},
		}

		if !cmp.Equal(exp, act1, cmpOptions()...) {
			t.Error(cmp.Diff(exp, act1, cmpOptions()...))
		}

		if !cmp.Equal(exp, act2, cmpOptions()...) {
			t.Error(cmp.Diff(exp, act2, cmpOptions()...))
		}
	})

	t.Run("success-with-deep-copy", func(t *testing.T) {
		cached := NewCache()
		act1 := New[Test](WithCache(cached), WithDeepCopy())
		act2 := New[Test](WithCache(cached), WithDeepCopy())

		exp := Test{
			Struct1: toPointer(Basic{
				Bool1:   true,
				Bool2:   false,
				Bool3:   true,
				Bool4:   false,
				String1: "abc",
				String2: "123",
			}, 1).(*Basic),
			Struct2: toPointer(Signed{
				Int:   toPointer(int(0), 1).(*int),
				Int8:  toPointer(int8(8), 2).(**int8),
				Int16: toPointer(int16(16), 3).(***int16),
				Int32: toPointer(int32(32), 4).(****int32),
				Int64: toPointer(int64(64), 5).(*****int64),
			}, 2).(**Signed),
			Struct3: toPointer(Unsigned{
				Uint:   toPointer(uint(0), 1).(*uint),
				Uint8:  toPointer(uint8(8), 2).(**uint8),
				Uint16: toPointer(uint16(16), 3).(***uint16),
				Uint32: toPointer(uint32(32), 4).(****uint32),
				Uint64: toPointer(uint64(64), 5).(*****uint64),
			}, 3).(***Unsigned),
			Struct4: toPointer(Float{
				Float1: float32(1.2345),
				Float2: float64(1.23456789),
			}, 4).(****Float),
			Struct5: toPointer(Complex{
				Complex1: complex64(1 + 2i),
				Complex2: complex128(3 + 4i),
			}, 5).(*****Complex),
			Arr1: [5]string{"1", "2", "3", "4", ""},
			Arr2: [2][]string{{"1", "2", "3", "4"}, {"5", "6", "7", "8"}},
			Arr3: [3]*Basic{
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
			},
			Arr4:   [4]int{1, 1, 1, 1},
			Slice1: []string{"1", "1", "1", "1", "1"},
			Slice2: []string{"1", "2", "3", "4", "5"},
			Slice3: []*Basic{
				{
					Bool1:   true,
					Bool2:   false,
					Bool3:   true,
					Bool4:   false,
					String1: "abc",
					String2: "123",
				},
			},
			Slice4:          [][2]string{{"1", "2"}, {"3", "4"}},
			Map1:            map[string]int{"1": 1},
			Map2:            map[string]int{"key1": 1, "key2": 2, "key3": 3},
			Duration1:       time.Second * 3,
			Duration2:       toPointer(((5 * time.Hour) + (30 * time.Minute) + (15 * time.Second)), 1).(*time.Duration),
			Time1:           time.Date(2024, 12, 9, 2, 20, 35, 0, time.UTC),
			Time2:           toPointer(time.Date(2024, 12, 9, 2, 20, 35, 0, time.UTC), 1).(*time.Time),
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
			Chan3:           make(<-chan int, 10),
			Interface1:      float64(123),
			Interface2:      true,
			Interface3:      "abc",
			Interface4:      map[string]any{"key": "value"},
			Interface5:      []any{"1", "2", "3"},
		}

		if !cmp.Equal(exp, act1, cmpOptions()...) {
			t.Error(cmp.Diff(exp, act1, cmpOptions()...))
		}

		if !cmp.Equal(exp, act2, cmpOptions()...) {
			t.Error(cmp.Diff(exp, act2, cmpOptions()...))
		}
	})
}

func Benchmark_NotCached(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New[Test]()
	}
}

func Benchmark_Cached(b *testing.B) {
	cached := NewCache()

	for i := 0; i < b.N; i++ {
		_ = New[Test](WithCache(cached))
	}
}

func Benchmark_DeepCopy(b *testing.B) {
	cached := NewCache()

	for i := 0; i < b.N; i++ {
		_ = New[Test](WithCache(cached), WithDeepCopy())
	}
}
