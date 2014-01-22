// Package regtest contains a bunch of helper functions for testing regression algorithms

package regtest

import (
	"github.com/gonum/floats"
	"math/rand"
	"testing"

	//"fmt"
)

const throwPanic = true

func panics(f func()) (b bool) {
	defer func() {
		err := recover()
		if err != nil {
			b = true
		}
	}()
	f()
	return
}

func maybe(f func()) (b bool) {
	defer func() {
		err := recover()
		if err != nil {
			b = true
			if throwPanic {
				panic(err)
			}
		}
	}()
	f()
	return
}

type ParameterGetterSetter interface {
	NumParameters() int
	Parameters([]float64) []float64
	SetParameters([]float64)
}

func TestGetAndSetParameters(t *testing.T, p ParameterGetterSetter, name string) {

	// Test that we can get parameters from nil
	// TODO: Add panic guard
	var nilParam []float64
	f := func() {
		nilParam = p.Parameters(nil)
	}

	if maybe(f) {
		t.Errorf("%v: Parameters panicked with nil input", name)
		return
	}

	if len(nilParam) != p.NumParameters() {
		t.Errorf("%v: On nil input, incorrect length returned from Parameters()", name)
	}
	nilParamCopy := make([]float64, p.NumParameters())
	copy(nilParamCopy, nilParam)
	nonNilParam := make([]float64, p.NumParameters())
	p.Parameters(nonNilParam)
	if !floats.Equal(nilParam, nonNilParam) {
		t.Errorf("%v: Return from Parameters() with nil argument and non nil argument are different", name)
	}
	for i := range nonNilParam {
		nonNilParam[i] = rand.NormFloat64()
	}
	if !floats.Equal(nilParam, nilParamCopy) {
		t.Errorf("%v: Modifying the return from Parameters modified the underlying parameters", name)
	}
	setParam := make([]float64, p.NumParameters())
	copy(setParam, nonNilParam)
	p.SetParameters(setParam)
	if !floats.Equal(setParam, nonNilParam) {
		t.Errorf("%v: Input slice modified during call to SetParameters", name)
	}

	afterParam := p.Parameters(nil)
	if !floats.Equal(afterParam, setParam) {
		t.Errorf("%v: Set parameters followed by Parameters don't return the same argument", name)
	}

	// Test that there are panics on bad length arguments
	badLength := make([]float64, p.NumParameters()+3)

	f = func() {
		p.Parameters(badLength)
	}
	if !panics(f) {
		t.Errorf("%v: Parameters did not panic given a slice too long", name)
	}
	f = func() {
		p.SetParameters(badLength)
	}
	if !panics(f) {
		t.Errorf("%v: SetParameters did not panic given a slice too long", name)
	}
	if p.NumParameters() == 0 {
		return
	}
	badLength = badLength[:p.NumParameters()-1]
	f = func() {
		p.Parameters(badLength)
	}
	if !panics(f) {
		t.Errorf("%v: Parameters did not panic given a slice too short", name)
	}
	f = func() {
		p.SetParameters(badLength)
	}
	if !panics(f) {
		t.Errorf("%v: SetParameters did not panic given a slice too short", name)
	}
}

type InputOutputer interface {
	InputDim() int
	OutputDim() int
}

func TestInputOutputDim(t *testing.T, io InputOutputer, trueInputDim, trueOutputDim int, name string) {
	inputDim := io.InputDim()
	outputDim := io.OutputDim()
	if inputDim != trueInputDim {
		t.Errorf("%v: Mismatch in input dimension. expected %v, found %v", name, trueInputDim, inputDim)
	}
	if outputDim != trueOutputDim {
		t.Errorf("%v: Mismatch in input dimension. expected %v, found %v", name, trueOutputDim, inputDim)
	}
}
