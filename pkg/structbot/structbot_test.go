package structbot_test

import (
	"fmt"
	"github.com/tjpnz/structbot/pkg/structbot"
)

type MyStruct struct {
	Foo string
	Bar string
}

func (s *MyStruct) String() string {
	return fmt.Sprintf("Foo: %s, Bar: %s", s.Foo, s.Bar)
}

func Example() {
	sb := structbot.New()
	sb.RegisterFactory("MyStruct", func() (interface{}, error) {
		return &MyStruct{
			Foo: "A",
			Bar: "B",
		}, nil
	})

	res, err := sb.Factory("MyStruct").Create()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	res, err = sb.Factory("MyStruct").Patch(&MyStruct{Foo: "AA"})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	// Output:
	// Foo: A, Bar: B
	// Foo: AA, Bar: B
}
