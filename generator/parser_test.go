package generator

import (
	"fmt"
	"testing"
)

func TestParseRepository(t *testing.T) {
	specs := ParseRepository([]string{"UserRepository"}, []string{"../testdata"}, []string{})
	if len(specs) != 1 {
		t.Fatalf("expect 1, got %d", len(specs))
	}

	for _, v := range specs[0].Methods {
		fmt.Println(v.FullName())
	}
}

func TestParseModel(t *testing.T) {
	specs := ParseModel([]string{"User"}, []string{"../testdata"}, []string{})
	if len(specs) != 1 {
		t.Fatalf("expect 1, got %d", len(specs))
	}

	for i := 0; i < specs[0].Struct.NumFields(); i++ {
		fmt.Println(specs[0].Struct.Field(i).Name(), specs[0].Struct.Field(i).Type(), specs[0].Struct.Tag(i))
	}
}
