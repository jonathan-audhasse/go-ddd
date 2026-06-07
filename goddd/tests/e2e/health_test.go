package e2e_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {

	res, err := http.Get(fmt.Sprintf(e2eBaseUrl + "/health"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// assertion
	resData, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, []byte("{}\n"), resData)
}
