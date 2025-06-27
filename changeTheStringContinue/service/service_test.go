package service

import (
	"errors"
	"github.com/MaksimHank/WorkingWithFile/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UnitTestSuite struct {
	suite.Suite
	mockProducer  *mocks.Producer
	mockPresenter *mocks.Presenter
	service       *Service
}

func (uts *UnitTestSuite) SetupTest() {
	uts.mockProducer = new(mocks.Producer)
	uts.mockPresenter = new(mocks.Presenter)
	uts.service = NewService(uts.mockProducer, uts.mockPresenter)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (uts *UnitTestSuite) TestNewServiceSuit() {
	assert.NotNil(uts.T(), uts.service)
	assert.Equal(uts.T(), uts.mockProducer, uts.service.prod)
	assert.Equal(uts.T(), uts.mockPresenter, uts.service.pres)
}

func (uts *UnitTestSuite) TestRunSuccess() {
	masked := []string{"Hello, it's my website http://***********;", "And this too http://********"}

	uts.mockProducer.On("Produce").Return(masked, nil)
	uts.mockPresenter.On("Present", masked).Return(nil)

	err := uts.service.Run()

	assert.NoError(uts.T(), err)
	uts.mockProducer.AssertExpectations(uts.T())
	uts.mockPresenter.AssertExpectations(uts.T())
}

func (uts *UnitTestSuite) TestRunProducerError() {
	uts.mockProducer.On("Produce").Return([]string{}, errors.New("fail on producer"))
	err := uts.service.Run()
	assert.Error(uts.T(), err)
	uts.mockPresenter.AssertNotCalled(uts.T(), "Present", mock.Anything)
}

func (uts *UnitTestSuite) TestRunPresenterError() {
	masked := []string{"Hello, it's my website http://***********;", "And this too http://********"}

	uts.mockProducer.On("Produce").Return(masked, nil)
	uts.mockPresenter.On("Present", masked).Return(errors.New("fail on presenter"))

	err := uts.service.Run()
	assert.Error(uts.T(), err)
}

func (uts *UnitTestSuite) TestChangeTheStringToAsterisks() {
	input := "Hello, it's my website http://example.com; And this too http://test.com"
	expected := "Hello, it's my website http://***********; And this too http://********"
	result := uts.service.changeTheStringToAsterisks(input)
	assert.Equal(uts.T(), expected, result)

	input2 := "No url here"
	expected2 := "No url here"
	result2 := uts.service.changeTheStringToAsterisks(input2)
	assert.Equal(uts.T(), expected2, result2)
}
