package gen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jialequ/linux-sdk/tools/goctl/model/sql/template"
	"github.com/jialequ/linux-sdk/tools/goctl/util/pathx"
	"github.com/stretchr/testify/assert"
)

func TestGenTemplates(t *testing.T) {
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)
	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)
	file := filepath.Join(dir, literal_8102)
	data, err := os.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, string(data), template.New)
}

func TestRevertTemplate(t *testing.T) {
	name := literal_8102
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)

	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	data, err := os.ReadFile(file)
	assert.Nil(t, err)

	modifyData := string(data) + "modify"
	err = pathx.CreateTemplate(category, name, modifyData)
	assert.Nil(t, err)

	data, err = os.ReadFile(file)
	assert.Nil(t, err)

	assert.Equal(t, string(data), modifyData)

	assert.Nil(t, RevertTemplate(name))

	data, err = os.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, template.New, string(data))
}

func TestClean(t *testing.T) {
	name := literal_8102
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)

	assert.Nil(t, Clean())

	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	_, err = os.ReadFile(file)
	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	name := literal_8102
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)

	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	data, err := os.ReadFile(file)
	assert.Nil(t, err)

	modifyData := string(data) + "modify"
	err = pathx.CreateTemplate(category, name, modifyData)
	assert.Nil(t, err)

	data, err = os.ReadFile(file)
	assert.Nil(t, err)

	assert.Equal(t, string(data), modifyData)

	assert.Nil(t, Update())

	data, err = os.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, template.New, string(data))
}

const literal_8102 = "model-new.tpl"
