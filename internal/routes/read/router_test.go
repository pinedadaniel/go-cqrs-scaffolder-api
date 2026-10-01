package read

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	mocks := initMocks()

	mocks.Router.On("Group", "", mock.Anything).Return(expectedRouteGroup())

	handlers := Handlers{
		CheckoutData: mocks.CheckoutData,
	}

	Register(mocks.Router, handlers)

	mocks.Router.AssertExpectations(t)
	mocks.CheckoutData.AssertExpectations(t)
}
