# StructBot

StructBot is a library for defining function based factories. It excels for use cases involving test cases where you would often define the same struct multiple times (specifically in table driven tests) or where you already have a struct initialized but need to modify a small handful of fields.

## Installation

`go get github.com/tjpnz/structbot`

## Usage

### Defining Factories

Factories are defined as functions returning a struct and error. They're registered to StructBot using `RegisterFactory`:

```go
type MyStruct struct {
    Foo string
    Bar string
}

sb := structbot.New()
sb.RegisterFactory("MyFactory", func() (interface{}, error) {
    return &MyStruct{
        Foo: "Foo",
        Bar: "Bar",
    }, nil
})
```

`RegisterFactory` returns an instance of StructBot which allows multiple factories to be registered at once using chaining.

### Using Factories

`Factory` returns a previously defined factory which can be invoked using `Create`, `Patch`, `MustCreate` or `MustPatch` - the latter two take a `testing.T` where errors are passed to `Fatalf` which is especially useful for testing.

```go
const baseCustomer = "BaseCustomer"

var sb *structbot.StructBot

func init() {
    sb = structbot.New()
    sb.RegisterFactory(baseCustomer, func() (interface{}, error) {
        return &Customer{
            FamilyName: "Cleese",
            FirstName: "John",
            Birthday: time.Date(1939, 10, 27, 0, 0, 0, 0, time.UTC),
            Address: &Address{}
        }, nil
    })
}

func Test_Create(t *testing.T) {
    for name, tc := range map[string]struct{
        in  *Customer
        out error
    }{
        "Valid": {
            in: sb.Factory(baseCustomer).MustCreate(t),
        }
        "ValidationError_Birthday": {
            in: sb.Factory(baseCustomer).MustPatch(t, &Customer{
                Birthday: time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
            }),
            out: ErrDateInFuture,
        },
        "ValidationError_City": {
            in: sb.Factory(baseCustomer).MustPatch(t, &Customer{
                Address: &Address{
                    City: "Palmerston North",
                }
            }),
            out: ErrInvalidCity,
        },
    } {
        t.Run(name, func(t *testing.T) {
            // Test code
        })
    }
}
```

### Patch Semantics

StructBot uses Mergo under the hood for patching with overrides enabled. Patching works with nested structs of arbitrary depth. Structs passed to `Patch` or `MustPatch` won't overwrite the value on the destination for fields where the source has a zero value. See the [Mergo docs](https://github.com/imdario/mergo#usage) for more.
