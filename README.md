# AutoStruct

`AutoStruct` is a Go library designed to streamline struct initialization by automatically populating fields based on struct tags.

## Installation

Install `auto-struct` using `go get`:

```bash
go get -u github.com/arsmn/auto-struct@latest
```

## Usage

### Defining a Struct

```go
type Person struct {
	Name string `auto:"person"`
	Age  int    `auto:"20"`
}
```

### Initializing a Struct

Using generics:
```go
person := autostruct.New[Person]()
```

Using a pre-defined variable (panics on error):
```go
var p Person
autostruct.MustSet(&p)
```

Using a pre-defined variable (returns error):
```go
var p Person
err := autostruct.Set(&p)
if err != nil {
	// handle error
}
```

## Supported Types

| Primitive Types      | Composite Types         |
|----------------------|-------------------------|
| `bool`               | `struct`                |
| `string`             | `array`                 |
| `int`                | `slice`                 |
| `int8`               | `map`                   |
| `int16`              | `pointer`               |
| `int32`              | `interface`             |
| `int64`              | `channel`               |
| `uint`               | `time.Time`             |
| `uint8`              | `time.Duration`         |
| `uint16`             | `json.RawMessage`       |
| `uint32`             |                         |
| `uint64`             |                         |
| `float32`            |                         |
| `float64`            |                         |
| `complex64`          |                         |
| `complex128`         |                         |
| `rune`               |                         |
| `byte`               |                         |

## Example

```go
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
```

## Options

### WithTag
Override the default tag value (`auto`) used for initialization.

```go
type Test struct {
	String string `default:"abc"`
}

test := autostruct.New[Test](autostruct.WithTag("default"))
```

### WithCache
Enable caching to improve performance by reusing generated values.

```go
type Test struct {
	String string `auto:"abc"`
}

cache := autostruct.NewCache()
test1 := autostruct.New[Test](autostruct.WithCache(cache))
test2 := autostruct.New[Test](autostruct.WithCache(cache))
```

### WithDeepCopy
When caching is enabled, reference types will point to the same values. Use DeepCopy to ensure each instance has its own copy.

```go
type Test struct {
	Slice []string `auto:"len(5),repeat(abc)"`
}

cache := autostruct.NewCache()
test1 := autostruct.New[Test](autostruct.WithCache(cache), autostruct.WithDeepCopy())
test2 := autostruct.New[Test](autostruct.WithCache(cache), autostruct.WithDeepCopy())
```
`test1.Slice` and `test2.Slice` will point to different underlying arrays.

## Benchmark

The following benchmarks were run on a Linux system (amd64) with an Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz:

| Benchmark          | Iterations | Time per Operation | Memory per Operation | Allocations per Operation |
|--------------------|------------|--------------------|----------------------|---------------------------|
| Cached             | 95,714     | 12,332 ns/op       | 4,346 B/op           | 177 allocs/op             |
| Deep Copy          | 73,220     | 16,838 ns/op       | 6,011 B/op           | 225 allocs/op             |
| Not Cached         | 10,000     | 105,668 ns/op      | 48,564 B/op          | 751 allocs/op             |
