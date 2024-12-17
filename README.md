# AutoStruct

`AutoStruct` is a Go library that simplifies struct initialization by automatically populating fields based on struct tags.

## Installation

You can install `auto-struct` using `go get`:

```bash
go get -u github.com/arsmn/auto-struct@latest
```

## Usage

### Defining

```go
type Person struct {
	Name string `auto:"person"`
	Age  int    `auto:"20"`
}
```

### Initialization

Using generics
```go
person := autostruct.New[Person]()
```

Using pre-defined variable (panics on error)
```go
var p Person
person := autostruct.MustSet(p)
```

Using pre-defined variable (returns error)
```go
var p Person
person, err := autostruct.Set(p)
if err != nil {
    // handle error
}
```

## Supported Types

### Primitive Types

- [x] `int`
- [x] `int8`
- [x] `int16`
- [x] `int32`
- [x] `int64`
- [x] `uint`
- [x] `uint8`
- [x] `uint16`
- [x] `uint32`
- [x] `uint64`
- [x] `uintptr`
- [x] `float32`
- [x] `float64`
- [x] `complex64`
- [x] `complex128`
- [x] `bool`
- [x] `string`
- [x] `duration`
- [X] `rune`
- [ ] `byte`

### Composite Types

- [x] `struct`
- [x] `array`
- [x] `slice`
- [x] `map`
- [x] `pointer`
- [x] `time`
- [ ] `interface`
- [ ] `channel`


## Example

```go
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
	Foo    **bool         `auto:"true"`
	Bar    string         `auto:"bar"`
	Qux    int8           `auto:"123"`
	Name   *Name          `auto:"value(struct)"`
	Age    **Age          `auto:"value(struct)"`
	Arr1   [5]string      `auto:"json([\"1\", \"2\", \"3\", \"4\"])"`
	Arr2   [2][]string    `auto:"json([[\"1\", \"2\", \"3\", \"4\"], [\"5\", \"6\", \"7\", \"8\"]])"`
	Arr3   [3]*Name       `auto:"repeat(struct)"`
	Arr4   [4]int         `auto:"repeat(1)"`
	Slice1 []string       `auto:"len(5),cap(10),repeat(1)"`
	Slice2 []string       `auto:"json([\"1\", \"2\", \"3\", \"4\", \"5\"])"`
	Slice3 []*Name        `auto:"len(1),repeat(struct)"`
	Slice4 [][2]string    `auto:"json([[\"1\", \"2\"], [\"3\", \"4\"]])"`
	Map1   map[string]int `auto:"json({\"1\": 1})"`
	Map2   map[string]int `auto:"len(5),value(key1:1,key2:2,key3:3)"`
	Dur1   time.Duration  `auto:"3s"`
	Dur2   *time.Duration `auto:"5h30m15s"`
	Time1  time.Time      `auto:"2024-12-09T02:20:35Z"`
	Time2  *time.Time     `auto:"value(2024-12-09 02:20:35),layout(DateTime)"`
}
```
