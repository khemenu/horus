package suitehook_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus/internal/suitehook"
)

type OuterDoubleSuite struct {
	suitehook.Suite

	v []string

	am_i_second bool
}

func TestOuterDoubleSuite(t *testing.T) {
	s := OuterDoubleSuite{
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

	s.BeforeTest(func() func() {
		s.v = append(s.v, "A1")
		return func() {
			s.v = append(s.v, "A8")
		}
	})
	s.BeforeSub(func() func() {
		s.v = append(s.v, "b1")
		return func() {
			s.v = append(s.v, "b8")
		}
	})

	s.AfterSub(func() {
		s.v = append(s.v, "y9")
	})
	s.AfterTest(func() {
		s.v = append(s.v, "Z9")
	})

	s.AfterSub(func() {
		s.v = append(s.v, "y8")
	})
	s.AfterTest(func() {
		s.v = append(s.v, "Z8")
	})

	suite.Run(t, &s)

	require.Equal(t, []string{
		"A0", "A1", // BeforeTest at init  ═════════════════╗
		"b0", "b1", // BeforeSub at init ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┐ ║
		//                                                ┊ ║
		// By 1st sub-test                                ┊ ║
		"d0", "d1", // BeforeSub at Foo/Bar ┄┄┄┄┄┄┄┄┄┄┄┄┐ ┊ ║
		"d8", "d9", // Return of BeforeSub at Foo/Bar ┄┄┘ ┊ ║
		"b8", "b9", // Return of BeforeSub at init ┄┄┄┄┄┄┄┘ ║
		"w8", "w9", // AfterSub at Foo/Bar                  ║
		"y8", "y9", // AfterSub at init                     ║
		//                                                  ║
		"X8", "X9", // AfterTest at Foo/Bar                 ║
		"A8", "A9", // Return of BeforeTest at init ════════╝
		"Z8", "Z9", // AfterTest at init
		//
		"A0", "A1", // BeforeTest at init  ═════════════════╗
		"b0", "b1", // BeforeSub at init ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┐ ║
		//                                                ┊ ║
		// By 2nd sub-test                                ┊ ║
		"d0", "d1", // BeforeSub at Foo/Bar ┄┄┄┄┄┄┄┄┄┄┄┄┐ ┊ ║
		"d8", "d9", // Return of BeforeSub at Foo/Bar ┄┄┘ ┊ ║
		"b8", "b9", // Return of BeforeSub at init ┄┄┄┄┄┄┄┘ ║
		"w8", "w9", // AfterSub at Foo/Bar                  ║
		"y8", "y9", // AfterSub at init                     ║
		//                                                  ║
		"X8", "X9", // AfterTest at Foo/Bar                 ║
		"A8", "A9", // Return of BeforeTest at init ════════╝
		"Z8", "Z9", // AfterTest at init
	}, s.v)
}

func (t *OuterDoubleSuite) TestFoo() {
	am_i_second := t.am_i_second
	t.am_i_second = true

	t.makeHook()

	t.Run("sub-test", func() {
		if am_i_second {
			t.secondCheck()
		} else {
			t.firstCheck()
		}
	})
}

func (t *OuterDoubleSuite) TestBar() {
	am_i_second := t.am_i_second
	t.am_i_second = true

	t.makeHook()

	t.Run("sub-test", func() {
		if am_i_second {
			t.secondCheck()
		} else {
			t.firstCheck()
		}
	})
}

func (t *OuterDoubleSuite) makeHook() {
	t.BeforeSub(func() func() {
		t.v = append(t.v, "d0")
		return func() {
			t.v = append(t.v, "d9")
		}
	})
	t.BeforeSub(func() func() {
		t.v = append(t.v, "d1")
		return func() {
			t.v = append(t.v, "d8")
		}
	})

	t.AfterSub(func() {
		t.v = append(t.v, "w9")
	})
	t.AfterTest(func() {
		t.v = append(t.v, "X9")
	})

	t.AfterSub(func() {
		t.v = append(t.v, "w8")
	})
	t.AfterTest(func() {
		t.v = append(t.v, "X8")
	})
}

func (t *OuterDoubleSuite) firstCheck() {
	t.Equal([]string{
		"A0", "A1",
		"b0", "b1",

		"d0", "d1",
	}, t.v)
}

func (t *OuterDoubleSuite) secondCheck() {
	t.Equal([]string{
		"A0", "A1",
		"b0", "b1",

		"d0", "d1",
		"d8", "d9",
		"b8", "b9",
		"w8", "w9",
		"y8", "y9",

		"X8", "X9",
		"A8", "A9",
		"Z8", "Z9",

		"A0", "A1",
		"b0", "b1",

		"d0", "d1",
	}, t.v)
}
