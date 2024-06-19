package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("test")
	assert.NotNil(t, err)
	errID := GetID(err)
	assert.Equal(t, NoID, errID)
	assert.Equal(t, "test", err.Error())
}

func TestWrap(t *testing.T) {
	err := Wrap(errors.New("test"), "Wrapped")
	assert.NotNil(t, err)
	errID := GetID(err)
	assert.Equal(t, NoID, errID)
	assert.Equal(t, "Wrapped: test", FullError(err))
	assert.Equal(t, "Wrapped", err.Error())
	assert.Equal(t, "test", Unwrap(err).Error())
}

func TestNewErrorID(t *testing.T) {
	testID := ErrorID(4)
	err := testID.New("test")
	assert.NotNil(t, err)
	errID := GetID(err)
	assert.Equal(t, testID, errID)
	assert.NotEqual(t, NoID, errID)
	assert.Equal(t, "(4) test", FullError(err))
	assert.Equal(t, "test", err.Error())
}

func TestWrapErrorID(t *testing.T) {
	testID := ErrorID(1)
	err := testID.Wrap(errors.New("test"), "Wrapped")
	assert.NotNil(t, err)
	errID := GetID(err)
	assert.Equal(t, testID, errID)
	assert.NotEqual(t, NoID, errID)
	assert.Equal(t, "(1) Wrapped: test", FullError(err))
	assert.Equal(t, "Wrapped", err.Error())
	assert.Equal(t, "test", Unwrap(err).Error())

	testID2 := ErrorID(2)
	err2 := testID2.Wrap(err, "Wrapped2")
	err2 = Wrap(err2, "Double")
	assert.NotNil(t, err)
	errID2 := GetID(err2)
	assert.Equal(t, testID2, errID2)
	assert.Equal(t, "(2) Double: Wrapped2: (1) Wrapped: test", FullError(err2))
	assert.Equal(t, "Double", err2.Error())
	assert.Equal(t, "(2) Wrapped2: (1) Wrapped: test", FullError(Unwrap(err2)))
	assert.Equal(t, "Wrapped2", Unwrap(err2).Error())
	assert.Equal(t, "(1) Wrapped: test", FullError(Unwrap(Unwrap(err2))))
	assert.Equal(t, "Wrapped", Unwrap(Unwrap(err2)).Error())
	assert.Equal(t, "test", FullError(Unwrap(Unwrap(Unwrap(err2)))))
	assert.Equal(t, "test", Unwrap(Unwrap(Unwrap(err2))).Error())
	assert.Equal(t, "", FullError(Unwrap(Unwrap(Unwrap(Unwrap(err2))))))
	assert.Equal(t, nil, Unwrap(Unwrap(Unwrap(Unwrap(err2)))))
	assert.Equal(t, true, errors.Is(err2, testID2.AsError()))
}
