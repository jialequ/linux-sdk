package util

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchFiles(t *testing.T) {
	dir, err := filepath.Abs("./")
	assert.Nil(t, err)

	files, err := MatchFiles("./*.sql")
	assert.Nil(t, err)
	assert.Equal(t, []string{filepath.Join(dir, "studeat.sql"), filepath.Join(dir, literal_4713), filepath.Join(dir, literal_8375)}, files)

	files, err = MatchFiles("./??.sql")
	assert.Nil(t, err)
	assert.Equal(t, []string{filepath.Join(dir, literal_8375)}, files)

	files, err = MatchFiles("./*.sq*")
	assert.Nil(t, err)
	assert.Equal(t, []string{filepath.Join(dir, "studeat.sql"), filepath.Join(dir, literal_4713), filepath.Join(dir, literal_8375), filepath.Join(dir, "xx.sql1")}, files)

	files, err = MatchFiles("./student.sql")
	assert.Nil(t, err)
	assert.Equal(t, []string{filepath.Join(dir, literal_4713)}, files)
}

const literal_4713 = "student.sql"

const literal_8375 = "xx.sql"
