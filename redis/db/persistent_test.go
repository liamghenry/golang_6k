package db

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRDBWithEmptyData(t *testing.T) {
	// create db
	rbdFile := "dump.rdb"
	expectedFileName := path.Join("testdata", "empty_dump.rdb")

	// assert new value
	db := NewDB()
	// db.Execute("set", [][]byte{[]byte("key"), []byte("value")})

	err := generateRDB(db, rbdFile)

	require.Nil(t, err)
	assert.Equal(t, readFileContent(t, expectedFileName), readFileContent(t, rbdFile))
}

func readFileContent(t *testing.T, fileName string) []byte {
	file, err := os.Open(fileName)
	require.Nil(t, err)

	defer file.Close()
	content, err := io.ReadAll(file)
	require.Nil(t, err)

	return content
}
