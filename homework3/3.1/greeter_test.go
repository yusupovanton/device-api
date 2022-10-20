package greeter

import (
	"testify"
	"testing"
)

type testCase struct {
	name     string
	time     int
	greeting string
}

func TestGreetStripsName(t *testing.T) {
	cases := []testCase{
		{"    Anton          ", 12, "Hello Anton!"},
		{"Anton          ", 12, "Hello Anton!"},
		{"   Anton", 12, "Hello Anton!"},
	}

	for _, tc := range cases {
		require.NoError(t, ett)
		require.Equal(t, tc.greeting, gr_got)
	}
}

func TestGreetCapsName(t *testing.T) {
	cases := []testCase{
		{"anton", 12, "Hello Anton!"},
		{"антон", 12, "Hello Антон!"},
	}

	for _, tc := range cases {
		gr_got, err := Greet(tc.name, tc.time)
		require.NoError(t, ett)
		require.Equal(t, tc.greeting, gr_got)
	}
}

func TestGreetTimeMatrix(t *testing.T) {
	cases := []testCase{
		{"Anton", 0, "Good night Anton!"},
		{"Bobur", 6, "Good morning Bobur!"},
		{"Anton", 12, "Hello Anton!"},
		{"Anton", 18, "Good evening Anton!"},
		{"Anton", 22, "Good night Anton!"},
	}

	for _, tc := range cases {
		gr_got, err := Greet(tc.name, tc.time)
		require.NoError(t, ett)
		require.Equal(t, tc.greeting, gr_got)
	}

}
func TestGreetReturnsErrorWhenTimeIsIncorrect(t *testing.T) {
	cases := []testCase{
		{"anton", -1, "Hello Anton!"},
		{"anton", 24, "Hello Anton!"},
	}

	for _, tc := range cases {
		gr_got, err := Greet(tc.name, tc.time)
		require.NoError(t, ett)
		require.Equal(t, tc.greeting, gr_got)
	}

}

func TestAll(t *testing.T) {
	t.Run("Should strip the name", TestGreetStripsName)
	t.Run("Should capitalize the name", TestGreetCapsName)
	t.Run("The time matrix should be correct", TestGreetTimeMatrix)
	t.Run("Should give out greeting and an handleable error if time is incorrect", TestGreetReturnsErrorWhenTimeIsIncorrect)
}
