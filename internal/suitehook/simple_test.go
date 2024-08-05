package suitehook_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus/internal/suitehook"
)

type SimpleSuite struct {
	suitehook.Suite

	v []string
}

func TestSimpleSuite(t *testing.T) {
	s := SimpleSuite{
		v: []string{},
	}

	s.BeforeTest(func() func() {
		s.v = append(s.v, "A0")
		return func() {
			s.v = append(s.v, "A9")
		}
	})
	s.BeforeSub(func() func() {
		s.v = append(s.v, "b0")
		return func() {
			s.v = append(s.v, "b9")
		}
	})

	s.AfterSub(func() {
		s.v = append(s.v, "y9")
	})
	s.AfterTest(func() {
		s.v = append(s.v, "Z9")
	})

	suite.Run(t, &s)

	require.Equal(t, []string{
		"A0", // BeforeTest at init  ═════════════╗
		"b0", // BeforeSub at init ┄┄┄┄┄┄┄┄┄┄┄┄┄┐ ║
		//                                      ┊ ║
		// By 1st sub-test                      ┊ ║
		"d0", // BeforeSub at Foo ┄┄┄┄┄┄┄┄┄┄┄┄┐ ┊ ║
		"d9", // Return of BeforeSub at Foo ┄┄┘ ┊ ║
		"b9", // Return of BeforeSub at init ┄┄┄┘ ║
		"w9", // AfterSub at Foo                  ║
		"y9", // AfterSub at init                 ║
		//                                        ║
		"X9", // AfterTest at Foo                 ║
		"A9", // Return of BeforeTest at init ════╝
		"Z9", // AfterTest at init
	}, s.v)
}

func (t *SimpleSuite) TestFoo() {
	t.Panics(func() {
		t.BeforeTest(func() func() {
			t.v = append(t.v, "C0")
			return nil
		})
	}, "BeforeTest should not be used after the test is started")

	t.BeforeSub(func() func() {
		t.v = append(t.v, "d0")
		return func() {
			t.v = append(t.v, "d9")
		}
	})
	t.AfterSub(func() {
		t.v = append(t.v, "w9")
	})
	t.AfterTest(func() {
		t.v = append(t.v, "X9")
	})

	t.Run("1st sub-test", func() {
		t.Equal([]string{
			"A0",
			"b0",
			"d0",
		}, t.v)
	})
}
