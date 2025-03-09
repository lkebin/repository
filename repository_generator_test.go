package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRepositoryImplements(t *testing.T) {
	specs := ParseRepository([]string{"UserRepository"}, []string{"repository/testdata"}, []string{})
	tpl, err := GenerateRepositoryImplements(&specs[0])

	assert.Nil(t, err)
	fmt.Println(tpl)
}
