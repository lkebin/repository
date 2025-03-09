package repository

import (
	"testing"
)

func TestID(t *testing.T) {
	t.Run("test invalid id", func(t *testing.T) {
		id := ID("e3654c3b-39c3-49e5-8cac-0dde20b05b2")
		if id.IsValid() {
			t.Errorf("expect false, got true")
		}
	})

	t.Run("test valid id", func(t *testing.T) {
		id := ID("154c48eb-a116-48ae-825c-deee2a664d33")
		if !id.IsValid() {
			t.Errorf("expect true, got false")
		}
	})

	t.Run("test short id", func(t *testing.T) {
		id := ID("154c48eb-a116-48ae-825c-deee2a664d33")
		if id.Short() != "154c48eba11648ae825cdeee2a664d33" {
			t.Errorf("expect 154c48eba11648ae825cdeee2a664d33, got %s", id.Short())
		}
	})

	t.Run("test parse id from short", func(t *testing.T) {
		id, err := ParseIDFromShort("154c48eba11648ae825cdeee2a664d33")
		if err != nil {
			t.Fatal(err)
		}
		if id.String() != "154c48eb-a116-48ae-825c-deee2a664d33" {
			t.Errorf("expect 154c48eb-a116-48ae-825c-deee2a664d33, got %s", id.String())
		}
	})
}
