package rf

import (
	"fmt"
	"testing"
)

type Test struct{ name string }

func (self Test) Repeat(n int) (count int, result string) {
	count = n
	for i := 0; i < n; i++ {
		result += self.name
	}
	return
}

func (self Test) local() {}

func TestCallMethod(t *testing.T) {
	var to = Test{":Test:"}
	var res, ok = CallMethod(to, "Repeat", 3)
	t.Log(res, ok)
	res, ok = CallMethod(to, "Unknown", 4)
	t.Log(res, ok)
	defer func() { t.Log(recover()) }()
	res, ok = CallMethod(to, "local")
}

func ExampleCallMethod() {
	var to = Test{":Test:"}
	var res, ok = CallMethod(to, "Repeat", 3)
	fmt.Println(res, ok)

	res, ok = CallMethod(to, "Unknown", 4)
	fmt.Println(res, ok)

	defer func() { fmt.Println(recover()) }()
	res, ok = CallMethod(to, "local")

	// Output:
	// [3 :Test::Test::Test:] true
	// [] false
	// CallMethod rf.Test.local() failed reflect: Call of unexported method
}

func ExampleCallMethod_nil() {
	var to Test
	var res, ok = CallMethod(to, "Repeat", 3)
	fmt.Println(res, ok)
	// Output:
	// [3 ] true
}

func ExampleCallMethod_insufficientarguments() {
	var to = Test{":Test:"}
	defer func() { fmt.Println(recover()) }()
	CallMethod(to, "Repeat")
	// Output:
	// CallMethod rf.Test.Repeat() failed reflect: Call with too few input arguments
}

func ExampleCallMethod_badarguments() {
	var to = Test{":Test:"}
	defer func() { fmt.Println(recover()) }()
	CallMethod(to, "Repeat", to)
	// Output:
	// CallMethod rf.Test.Repeat() failed reflect: Call using rf.Test as type int
}

func ExampleCallMethod_extraarguments() {
	var to = Test{":Test:"}
	defer func() { fmt.Println(recover()) }()
	CallMethod(to, "Repeat", 2, "extra")
	// Output:
	// CallMethod rf.Test.Repeat() failed reflect: Call with too many input arguments
}

func ExampleCallMethod_dontpanic() {
	var dontPanic = func(to Test, m string, a interface{}) (res []interface{}, err error) {
		defer func() { err = recover().(error) }()
		res, _ = CallMethod(to, m, a)
		return
	}

	res, err := dontPanic(Test{":Test:"}, "Repeat", "string")
	fmt.Println(res, err)
	// Output:
	// [] CallMethod rf.Test.Repeat() failed reflect: Call using string as type int
}

func BenchmarkNoCallMethod(b *testing.B) {
	//b.ReportAllocs()
	var to = Test{":Test:"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		to.Repeat(3)
	}
}
func BenchmarkCallMethod(b *testing.B) {
	//b.ReportAllocs()
	var to = Test{":Test:"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CallMethod(to, "Repeat", 3)
	}
}
