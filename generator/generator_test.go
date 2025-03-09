package generator

import (
	"fmt"
	"testing"
)

func TestGenerateRepositoryImplements(t *testing.T) {
	t.Run("UserRepository", func(t *testing.T) {
		specs := ParseRepository([]string{"UserRepository"}, []string{"../testdata"}, []string{})
		src, err := GenerateRepositoryImplements(&specs[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		fmt.Println(string(src))
	})

	t.Run("UserUuidRepository", func(t *testing.T) {
		specs := ParseRepository([]string{"UserUuidRepository"}, []string{"../testdata"}, []string{})
		src, err := GenerateRepositoryImplements(&specs[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		fmt.Println(string(src))
	})
}
