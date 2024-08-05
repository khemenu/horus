package suitehook

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PreHook func() func()
type PostHook func()

type frame struct {
	beforeTest []PreHook
	beforeSub  []PreHook
	afterSub   []PostHook
	afterTest  []PostHook
}

type Suite struct {
	suite.Suite
	*require.Assertions

	outer frame
	inner frame
	curr  frame

	has_been_init bool
}

func (s *Suite) BeforeTest(f PreHook) {
	if s.has_been_init {
		panic("test already been setup!")
	}
	s.curr.beforeTest = append(s.curr.beforeTest, f)
}

func (s *Suite) BeforeSub(f PreHook) {
	s.curr.beforeSub = append(s.curr.beforeSub, f)
}

func (s *Suite) AfterSub(f PostHook) {
	s.curr.afterSub = append(s.curr.afterSub, f)
}

func (s *Suite) AfterTest(f PostHook) {
	s.curr.afterTest = append(s.curr.afterTest, f)
}

func (s *Suite) SetupTest() {
	s.Assertions = s.Require()
	if !s.has_been_init {
		s.has_been_init = true
		s.outer = s.curr
	} else {
		s.curr = s.outer
	}

	for _, f := range s.curr.beforeTest {
		g := f()
		if g != nil {
			s.curr.afterTest = append(s.curr.afterTest, g)
		}
	}
}

func (s *Suite) SetupSubTest() {
	s.inner = s.curr
	for _, f := range s.curr.beforeSub {
		g := f()
		if g != nil {
			s.curr.afterSub = append(s.curr.afterSub, g)
		}
	}
}

func (s *Suite) TearDownSubTest() {
	for i := range s.curr.afterSub {
		s.curr.afterSub[len(s.curr.afterSub)-i-1]()
	}
	s.curr = s.inner
}

func (s *Suite) TearDownTest() {
	for i := range s.curr.afterTest {
		s.curr.afterTest[len(s.curr.afterTest)-i-1]()
	}
}

func (s *Suite) Run(name string, test func()) bool {
	ok := s.Suite.Run(name, func() {
		s.Assertions = s.Require()
		test()
	})

	s.Assertions = s.Require()
	return ok
}
