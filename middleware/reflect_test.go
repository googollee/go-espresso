package middleware

import (
	"fmt"
	"reflect"
	"testing"
)

type SomeStruct struct {
	Str string
	I   int
	F   float32
}

func BenchmarkTypeOf(b *testing.B) {
	got := reflect.TypeOf(&SomeStruct{}).String()
	if got, want := got, "*espresso.SomeStruct"; got != want {
		b.Fatalf("got: %v, want: %v", got, want)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = reflect.TypeOf(&SomeStruct{}).String()
	}
}

func BenchmarkSprintfT(b *testing.B) {
	got := fmt.Sprintf("%T", &SomeStruct{})
	if got, want := got, "*espresso.SomeStruct"; got != want {
		b.Fatalf("got: %v, want: %v", got, want)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%T", &SomeStruct{})
	}
}
