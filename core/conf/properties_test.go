package conf

import (
	"os"
	"testing"

	"github.com/jialequ/linux-sdk/core/fs"
	"github.com/stretchr/testify/assert"
)

func TestProperties(t *testing.T) {
	text := `app.name = test

    app.program=app

    # this is comment
    app.threads = 5`
	tmpfile, err := fs.TempFilenameWithText(text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

	props, err := LoadProperties(tmpfile)
	assert.Nil(t, err)
	assert.Equal(t, "test", props.GetString(literal_0452))
	assert.Equal(t, "app", props.GetString(literal_7948))
	assert.Equal(t, 5, props.GetInt(literal_4038))

	val := props.ToString()
	assert.Contains(t, val, literal_0452)
	assert.Contains(t, val, literal_7948)
	assert.Contains(t, val, literal_4038)
}

func TestPropertiesEnv(t *testing.T) {
	text := `app.name = test

    app.program=app

	app.env1 = ${FOO}
	app.env2 = $none

    # this is comment
    app.threads = 5`
	tmpfile, err := fs.TempFilenameWithText(text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

	t.Setenv("FOO", "2")

	props, err := LoadProperties(tmpfile, UseEnv())
	assert.Nil(t, err)
	assert.Equal(t, "test", props.GetString(literal_0452))
	assert.Equal(t, "app", props.GetString(literal_7948))
	assert.Equal(t, 5, props.GetInt(literal_4038))
	assert.Equal(t, "2", props.GetString("app.env1"))
	assert.Equal(t, "", props.GetString("app.env2"))

	val := props.ToString()
	assert.Contains(t, val, literal_0452)
	assert.Contains(t, val, literal_7948)
	assert.Contains(t, val, literal_4038)
	assert.Contains(t, val, "app.env1")
	assert.Contains(t, val, "app.env2")
}

func TestLoadPropertiesbadContent(t *testing.T) {
	filename, err := fs.TempFilenameWithText("hello")
	assert.Nil(t, err)
	defer os.Remove(filename)
	_, err = LoadProperties(filename)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
}

func TestSetString(t *testing.T) {
	key := "a"
	value := "the value of a"
	props := NewProperties()
	props.SetString(key, value)
	assert.Equal(t, value, props.GetString(key))
}

func TestSetInt(t *testing.T) {
	key := "a"
	value := 101
	props := NewProperties()
	props.SetInt(key, value)
	assert.Equal(t, value, props.GetInt(key))
}

func TestLoadBadFile(t *testing.T) {
	_, err := LoadProperties("nosuchfile")
	assert.NotNil(t, err)
}

const literal_0452 = "app.name"

const literal_7948 = "app.program"

const literal_4038 = "app.threads"
