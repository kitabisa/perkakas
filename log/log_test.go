package log

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	gock "gopkg.in/h2non/gock.v1"
)

type LogTestSuite struct {
	suite.Suite
	Logger       *Logger
	host         string
	endpoint     string
	url          string
	mockResponse map[string]interface{}
	responseTime string
	logID        string
}

func (suite *LogTestSuite) SetupTest() {
	suite.Logger = NewLogger("test")
	suite.logID = "3891967c-8589-42e0-a493-4fb6a0287992" // should be dynamic uuid. This static uuid was for testing purpose
	suite.host = "https://programming-quotes-api.herokuapp.com"
	suite.endpoint = "quotes"
	suite.url = fmt.Sprintf("%s/%s", suite.host, suite.endpoint)
	suite.responseTime = "2019-09-30T10:25:14+07:00"
	suite.mockResponse = map[string]interface{}{
		"id":     "a1234-abcd",
		"quotes": "Everything the light touches, is our kingdom",
	}
}

func (suite *LogTestSuite) TestEmptyLog() {
	suite.Logger.Print()
}

func (suite *LogTestSuite) TestMultipleCallLog() {
	suite.Logger.Print()
	suite.Logger.Print("message1", "message2")
	suite.Logger.AddMessage(ErrorLevel, "This is error")
	suite.Logger.AddMessage(ErrorLevel)
	suite.Logger.Print()
	suite.Logger.AddMessage(InfoLevel, "This is info")
	suite.Logger.Print("message3", "message4")
}

func (suite *LogTestSuite) TestLog() {
	defer gock.Off()

	req := gock.New(suite.host).
		Get(suite.endpoint).
		Reply(200).
		AddHeader("X-Test-Response", "Hello").
		AddHeader("Time", suite.responseTime).
		JSON(suite.mockResponse)

	resp, err := http.Get(suite.url)
	assert.Nil(suite.T(), err, "Nil expected")

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	suite.Logger.SetRequest(req)
	suite.Logger.SetResponse(resp, b)
	suite.Logger.AddMessage(TraceLevel, "This is trace")
	suite.Logger.AddMessage(DebugLevel, "This is debug")
	suite.Logger.AddMessage(InfoLevel, "This is info")
	suite.Logger.AddMessage(WarnLevel, "This is warning")
	suite.Logger.AddMessage(ErrorLevel, errors.New("This is error"))
	suite.Logger.AddMessage(FatalLevel, "This is fatal")
	suite.Logger.AddMessage(PanicLevel, "This is panic")
	suite.Logger.Print()
	assert.Equal(suite.T(), true, gock.IsDone(), "must be equal")
}

func (suite *LogTestSuite) TestMessageStackType() {
	message := map[string]interface{}{
		"id": "12345",
	}
	assert.Panics(suite.T(), func() { ensureStackType(message) }, "should panic")
}

func (suite *LogTestSuite) TestFindMaxLevel() {
	msgs := []message{
		message{Level: TraceLevel},
		message{Level: DebugLevel},
		message{Level: WarnLevel},
		message{Level: ErrorLevel},
	}

	maxLevel := suite.Logger.findMaxLevel(msgs)
	assert.Equal(suite.T(), ErrorLevel, maxLevel)
}

func TestLogTestSuite(t *testing.T) {
	suite.Run(t, new(LogTestSuite))
}
