package espresso

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDeclaratorNormal(t *testing.T) {
	ctx := &declareContext[struct{}]{}
	fakeHandler := func(Context[struct{}]) error { return nil }
	cmpType := func(a, b reflect.Type) bool {
		if a == nil && b == nil {
			return true
		}
		if a == nil || b == nil {
			return false
		}
		return a.String() == b.String()
	}
	cmpBindFunc := func(a, b bindFunc) bool {
		if a == nil && b == nil {
			return true
		}
		if a == nil || b == nil {
			return false
		}
		return fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b)
	}
	cmpHandlers := func(a, b []Handler[struct{}]) bool {
		if a == nil && b == nil {
			return true
		}
		if a == nil || b == nil {
			return false
		}
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	}

	declarator := ctx.Endpoint("method", "/path/with/:param1/and/:param2", fakeHandler)
	if diff := cmp.Diff(ctx.brew.handlers, []Handler[struct{}]{fakeHandler}, cmp.Comparer(cmpHandlers)); diff != "" {
		t.Errorf("declare handlers diff: (-got, +want)\n%s", diff)
	}

	var param1 string
	var param2 int
	declarator.BindPath("param1", &param1)
	declarator.BindPath("param2", &param2)
	var form1 string
	var form2 int
	declarator.BindForm("form1", &form1)
	declarator.BindForm("form2", &form2)

	declarator.Response("response/type")

	func() {
		defer func() {
			r := recover()
			if c, ok := r.(declareChcecker); ok && c.DeclareDone() {
				return
			}
			t.Fatalf("call declarator.End() panics with %v, should be endpointDeclareFinished{}", r)
		}()

		if err := declarator.End(); err != nil {
			t.Errorf("declarator.End() returns error: %v", err)
		}
	}()

	want := &endpoint{
		Method: "method",
		Path:   "/path/with/:param1/and/:param2",
		PathParams: []*binding{
			{"param1", bindStr, reflect.TypeOf("")},
			{"param2", bindInt, reflect.TypeOf(int(0))},
		},
		FormParams: []*binding{
			{"form1", bindStr, reflect.TypeOf("")},
			{"form2", bindInt, reflect.TypeOf(int(0))},
		},
		ResponseMime: "response/type",
	}

	if diff := cmp.Diff(ctx.endpoint, want, cmp.Comparer(cmpType), cmp.Comparer(cmpBindFunc)); diff != "" {
		t.Errorf("endpoint diff: (-got +want)\n%s", diff)
	}
}

func TestEndpointBindPath(t *testing.T) {
	tests := []struct {
		path         string
		binding      func(d Declarator)
		wantPanicStr string
		wantBinding  []*binding
	}{
		{"/path/no/params", func(d Declarator) {}, "", nil},
		{"/path/int/:param", func(d Declarator) {
			var i int
			d.BindPath("param", &i)
		}, "", []*binding{
			{"param", bindInt, reflect.TypeOf(int(0))},
		}},
		{"/path/one/more/:param/:more", func(d Declarator) {
			var i int
			d.BindPath("param", &i)
		}, "didn't bind any variables with path params [more]", []*binding{
			{"param", bindInt, reflect.TypeOf(int(0))},
		}},
		{"/path/not/exist/param", func(d Declarator) {
			var i int
			d.BindPath("param", &i)
		}, "can't find variables with name param in path /path/not/exist/param", []*binding{
			{"param", bindInt, reflect.TypeOf(int(0))},
		}},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			ctx := &declareContext[struct{}]{}
			declarator := ctx.Endpoint("method", test.path)

			panicStr := ""
			func() {
				defer func() {
					r := recover()
					if c, ok := r.(declareChcecker); ok && c.DeclareDone() {
						return
					}
					panicStr = fmt.Sprintf("%v", r)
				}()
				test.binding(declarator)

				if err := declarator.End(); err != nil {
					t.Errorf("declarator.End() returns error: %v", err)
				}
			}()

			if got, want := panicStr, test.wantPanicStr; got != want {
				t.Errorf("panic with %q, want %q", got, want)
			}
		})
	}
}
