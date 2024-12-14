package auth

import (
	"testing"
)

func TestAuth(t *testing.T){
	cases := []struct {
		input string
	}{
		{
			input: "Prakhar@02",
		},
	}

	for _, c := range cases {
		hashed_password, err := HashPassword(c.input)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		err = CheckHashPassword(c.input, hashed_password)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}

}