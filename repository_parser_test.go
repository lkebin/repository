package repository

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestParseRepository(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	var specs = ParseRepository([]string{"UserRepository"}, []string{"repository/testdata"}, []string{})
	assert.Len(t, specs, 1)
	for _, v := range specs[0].Methods {
		fmt.Println(v.FullName())
	}
}

func TestParseModel(t *testing.T) {
	var specs = ParseModel([]string{"User"}, []string{"repository/testdata"}, []string{})
	assert.Len(t, specs, 1)

	for i := 0; i < specs[0].Struct.NumFields(); i++ {
		fmt.Println(specs[0].Struct.Field(i).Name(), specs[0].Struct.Field(i).Type(), specs[0].Struct.Tag(i))
	}
}
