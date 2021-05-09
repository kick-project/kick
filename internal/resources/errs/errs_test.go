package errs_test

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"testing"

	errs "github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/stretchr/testify/assert"
)

func setup() (*errs.Handler, *bytes.Buffer) {
	str := ``
	buf := bytes.NewBufferString(str)
	e := errs.Handler{
		Ex: &exit.Handler{
			Mode: exit.MPanic,
		},
		Logger: log.New(buf, "", log.LstdFlags),
	}

	return &e, buf
}

func expectPanic(t *testing.T, r interface{}, msg string) {
	assert.NotNil(t, r)
	switch v := r.(type) {
	case error:
		assert.Equal(t, r.(error).Error(), msg)
	default:
		t.Errorf("Unexpected recovery type \"%T\"", v)
	}
}

func expectNil(t *testing.T, r interface{}) {
	assert.Nil(t, r)
}

// TestErrors_Panic_Paniced tests to see if a panic ensused and succeeds if it does
func TestErrors_Panic_Paniced(t *testing.T) {
	e, _ := setup()
	msg := "my error"

	defer func() {
		r := recover()
		expectPanic(t, r, msg)
	}()
	err := errors.New(msg)
	e.Panic(err)
}

// TestErrors_Panic_None tests to see if a panic ensused and fails if it does
func TestErrors_Panic_None(t *testing.T) {
	e, _ := setup()

	defer func() {
		r := recover()
		expectNil(t, r)
	}()
	e.Panic(nil)
}

// TestErrors_PanicF_Paniced tests to see if a panic ensused and succeeds if it does
func TestErrors_PanicF_Paniced(t *testing.T) {
	e, _ := setup()
	format := `Expecting a panic for error: %v`
	msg := "My error"

	defer func() {
		r := recover()
		expectPanic(t, r, fmt.Sprintf(format, msg))
	}()
	err := errors.New(msg)
	e.PanicF(format, err)
}

// TestErrors_PanicF_Paniced tests to see if a panic ensused and succeeds if it does
func TestErrors_PanicF_None(t *testing.T) {
	e, _ := setup()

	defer func() {
		r := recover()
		expectNil(t, r)
	}()
	e.PanicF(`Expecting no panic for nil: %v`, nil)
}

// TestErrors_LogF_True tests that that a log message is created when an error has occurred
func TestErrors_LogF_True(t *testing.T) {
	e, _ := setup()
	format := `Logging an error: %v`
	msg := "My error"

	err := errors.New(msg)
	assert.True(t, e.LogF(format, err))
}

// TestErrors_LogF_True tests that that a log message is created when an error has occurred
func TestErrors_LogF_False(t *testing.T) {
	e, _ := setup()
	format := `Logging an error: %v`

	assert.False(t, e.LogF(format, nil))
}

// TestErrors_Fatal tests that the a the code will exit and log a message when an error is encountered
func TestErrors_Fatal(t *testing.T) {
	e, _ := setup()
	msg := "My error"

	err := errors.New(msg)
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		switch v := r.(type) {
		case string:
			assert.Equal(t, r.(string), "Exit 255\n")
		default:
			t.Errorf("Unexpected recovery type \"%T\"", v)
		}
	}()
	e.Fatal(err)
}

// TestErrors_Fatalf tests that the a the code will exit and log a message when an error is encountered
func TestErrors_Fatalf(t *testing.T) {
	e, _ := setup()
	format := `Expecting an error: %v`
	msg := "My error"
	err := errors.New(msg)

	defer func() {
		r := recover()
		assert.NotNil(t, r)
		switch v := r.(type) {
		case string:
			assert.Equal(t, r.(string), "Exit 255\n")
		default:
			t.Errorf("Unexpected recovery type \"%T\"", v)
		}

	}()
	e.FatalF(format, err)
}

// TestErrors_Panic_Paniced tests to see if a panic ensused and succeeds if it does
func TestPanic_Paniced(t *testing.T) {
	exit.Mode(exit.MPanic)
	msg := "my error"
	defer func() {
		r := recover()
		expectPanic(t, r, msg)
	}()
	err := errors.New(msg)
	errs.Panic(err)
}

// TestErrors_Panic_None tests to see if a panic ensused and fails if it does
func TestPanic_None(t *testing.T) {
	exit.Mode(exit.MPanic)
	defer func() {
		r := recover()
		expectNil(t, r)
	}()
	errs.Panic(nil)
}

// TestErrors_PanicF_Paniced tests to see if a panic ensused and succeeds if it does
func TestPanicF_Paniced(t *testing.T) {
	exit.Mode(exit.MPanic)
	format := `Expecting a panic for error: %v`
	msg := "My error"

	defer func() {
		r := recover()
		expectPanic(t, r, fmt.Sprintf(format, msg))
	}()
	err := errors.New(msg)
	errs.PanicF(format, err)
}

// TestErrors_PanicF_Paniced tests to see if a panic ensused and succeeds if it does
func TestPanicF_None(t *testing.T) {
	exit.Mode(exit.MPanic)
	defer func() {
		r := recover()
		expectNil(t, r)
	}()
	errs.PanicF(`Expecting no panic for nil: %v`, nil)
}

// TestErrors_LogF_True tests that that a log message is created when an error has occurred
func TestLogF_True(t *testing.T) {
	exit.Mode(exit.MPanic)
	format := `Logging an error: %v`
	msg := "My error"

	err := errors.New(msg)
	assert.True(t, errs.LogF(format, err))
}

// TestErrors_LogF_True tests that that a log message is created when an error has occurred
func TestLogF_False(t *testing.T) {
	exit.Mode(exit.MPanic)
	format := `Logging an error: %v`

	assert.False(t, errs.LogF(format, nil))
}

// TestErrors_Fatal tests that the a the code will exit and log a message when an error is encountered
func TestFatal(t *testing.T) {
	exit.Mode(exit.MPanic)
	msg := "My error"

	err := errors.New(msg)
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		switch v := r.(type) {
		case string:
			assert.Equal(t, r.(string), "Exit 255\n")
		default:
			t.Errorf("Unexpected recovery type \"%T\"", v)
		}
	}()
	errs.Fatal(err)
}

// TestErrors_Fatalf tests that the a the code will exit and log a message when an error is encountered
func TestFatalf(t *testing.T) {
	exit.Mode(exit.MPanic)
	format := `Expecting an error: %v`
	msg := "My error"
	err := errors.New(msg)

	defer func() {
		r := recover()
		assert.NotNil(t, r)
		switch v := r.(type) {
		case string:
			assert.Equal(t, r.(string), "Exit 255\n")
		default:
			t.Errorf("Unexpected recovery type \"%T\"", v)
		}

	}()
	errs.FatalF(format, err)
}
