package rf

import ref "reflect"
import "fmt"

/* CallMethod calls object method by name
 * object.method(args...)
 */
func CallMethod(object interface{}, method string, args ...interface{}) ([]interface{}, bool) {
	if m := ref.ValueOf(object).MethodByName(method); !m.IsValid() {
		return nil, false
	} else {
		var a = make([]ref.Value, len(args))
		for i, arg := range args {
			a[i] = ref.ValueOf(arg)
		}
		defer func() {
			if r := recover(); r != nil {
				panic(CallMethodPanic{object, method, r})
			}
		}()
		a = m.Call(a)
		var res = make([]interface{}, len(a))
		for i, r := range a {
			res[i] = r.Interface()
		}
		return res, true
	}
}

type CallMethodPanic struct {
	Object interface{}
	Method string
	Err    interface{}
}

func (self CallMethodPanic) Error() string {
	return fmt.Sprintf("CallMethod %T.%s() failed %s",
		self.Object, self.Method, self.Err)
}
