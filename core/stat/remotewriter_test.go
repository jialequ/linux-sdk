package stat

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestRemoteWriter(t *testing.T) {
	defer gock.Off()

	gock.New(literal_7360).Reply(200).BodyString("foo")
	writer := NewRemoteWriter(literal_7360)
	err := writer.Write(&StatReport{
		Name: "bar",
	})
	assert.Nil(t, err)
}

func TestRemoteWriterFail(t *testing.T) {
	defer gock.Off()

	gock.New(literal_7360).Reply(503).BodyString("foo")
	writer := NewRemoteWriter(literal_7360)
	err := writer.Write(&StatReport{
		Name: "bar",
	})
	assert.NotNil(t, err)
}

const literal_7360 = "http://foo.com"
