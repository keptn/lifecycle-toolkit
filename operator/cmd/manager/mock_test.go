package manager

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockManager(t *testing.T) {
	mgr := MockManager{}

	mgr.On("Start", mock.Anything).Return(errors.New(""))
	mgr.On("AddHealthzCheck", mock.Anything, mock.Anything).Return(errors.New(""))
	mgr.On("AddReadyzCheck", mock.Anything, mock.Anything).Return(errors.New(""))

	assert.EqualError(t, mgr.Start(context.TODO()), "")
	assert.EqualError(t, mgr.AddHealthzCheck("", nil), "")
	assert.EqualError(t, mgr.AddReadyzCheck("", nil), "")

	mgr.AssertCalled(t, "Start", mock.Anything)
	mgr.AssertCalled(t, "AddHealthzCheck", mock.Anything, mock.Anything)
	mgr.AssertCalled(t, "AddReadyzCheck", mock.Anything, mock.Anything)
}
