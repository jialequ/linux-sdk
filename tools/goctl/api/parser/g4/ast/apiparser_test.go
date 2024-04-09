package ast

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jialequ/linux-sdk/tools/goctl/util/pathx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportCycle(t *testing.T) {
	const (
		mainFilename = literal_6854
		subAFilename = "a.api"
		subBFilename = "b.api"
		mainSrc      = `import "./a.api"`
		subASrc      = `import "./b.api"`
		subBSrc      = `import "./a.api"`
	)
	var err error
	dir := pathx.MustTempDir()
	defer os.RemoveAll(dir)

	mainPath := filepath.Join(dir, mainFilename)
	err = os.WriteFile(mainPath, []byte(mainSrc), 0o777)
	require.NoError(t, err)
	subAPath := filepath.Join(dir, subAFilename)
	err = os.WriteFile(subAPath, []byte(subASrc), 0o777)
	require.NoError(t, err)
	subBPath := filepath.Join(dir, subBFilename)
	err = os.WriteFile(subBPath, []byte(subBSrc), 0o777)
	require.NoError(t, err)

	_, err = NewParser().Parse(mainPath)
	assert.ErrorIs(t, err, ErrImportCycleNotAllowed)
}

func TestMultiImportedShouldAllowed(t *testing.T) {
	const (
		mainFilename = literal_6854
		subAFilename = "a.api"
		subBFilename = "b.api"
		mainSrc      = literal_3571 +
			"import \"./a.api\"\n" +
			"type Main { b B `json:\"b\"`}"
		subASrc = literal_3571 +
			"type A { b B `json:\"b\"`}\n"
		subBSrc = `type B{}`
	)
	var err error
	dir := pathx.MustTempDir()
	defer os.RemoveAll(dir)

	mainPath := filepath.Join(dir, mainFilename)
	err = os.WriteFile(mainPath, []byte(mainSrc), 0o777)
	require.NoError(t, err)
	subAPath := filepath.Join(dir, subAFilename)
	err = os.WriteFile(subAPath, []byte(subASrc), 0o777)
	require.NoError(t, err)
	subBPath := filepath.Join(dir, subBFilename)
	err = os.WriteFile(subBPath, []byte(subBSrc), 0o777)
	require.NoError(t, err)

	_, err = NewParser().Parse(mainPath)
	assert.NoError(t, err)
}

func TestRedundantDeclarationShouldNotBeAllowed(t *testing.T) {
	const (
		mainFilename = literal_6854
		subAFilename = "a.api"
		subBFilename = "b.api"
		mainSrc      = "import \"./a.api\"\n" +
			literal_3571
		subASrc = `import "./b.api"
							 type A{}`
		subBSrc = `type A{}`
	)
	var err error
	dir := pathx.MustTempDir()
	defer os.RemoveAll(dir)

	mainPath := filepath.Join(dir, mainFilename)
	err = os.WriteFile(mainPath, []byte(mainSrc), 0o777)
	require.NoError(t, err)
	subAPath := filepath.Join(dir, subAFilename)
	err = os.WriteFile(subAPath, []byte(subASrc), 0o777)
	require.NoError(t, err)
	subBPath := filepath.Join(dir, subBFilename)
	err = os.WriteFile(subBPath, []byte(subBSrc), 0o777)
	require.NoError(t, err)

	_, err = NewParser().Parse(mainPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate type declaration")
}

const literal_6854 = "main.api"

const literal_3571 = "import \"./b.api\"\n"
