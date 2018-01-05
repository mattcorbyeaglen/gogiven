package gogiven

import (
	"testing"
	"runtime"
)

type some struct {
	context              *TestingT
	InterestingGivens    *InterestingGivens
	CapturedIO           *CapturedIO
	runtimeCaller        *runtime.Func
	frameProgramCounters []uintptr
}

func newSome(testContext *TestingT,
	runtimeCaller *runtime.Func,
	frameProgramCounters [] uintptr,
	givenFunc ...func(givens *InterestingGivens)) *some {

	some := new(some)
	some.runtimeCaller = runtimeCaller
	some.frameProgramCounters = frameProgramCounters
	some.context = testContext
	some.CapturedIO = newCapturedIO()
	givens := newInterestingGivens()

	if len(givenFunc) > 0 {
		for _, someGivenFunc := range givenFunc {
			someGivenFunc(givens)
		}
	}
	some.InterestingGivens = givens
	return some
}
func newInterestingGivens() *InterestingGivens {
	givens := new(InterestingGivens)
	givens.Givens = map[string]interface{}{}
	return givens
}
func newCapturedIO() *CapturedIO {
	capturedIO := new(CapturedIO)
	capturedIO.CapturedIO = map[string]interface{}{}
	return capturedIO
}

func (some *some) When(action ...func(actual *CapturedIO, givens *InterestingGivens)) *some {
	action[0](some.CapturedIO, some.InterestingGivens) // TODO: there could be multiple actions..
	return some
}

func (some *some) Then(assertions func(testingT *TestingT, actual *CapturedIO, givens *InterestingGivens)) *some {
	assertions(some.context, some.CapturedIO, some.InterestingGivens)
	generateTestOutput(some)
	return some
}

func newTestMetaData(t *testing.T) *TestingT {
	testContext := new(TestingT)
	testContext.t = t
	return testContext
}