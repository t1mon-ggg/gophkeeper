package tls_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/t1mon-ggg/gophkeeper/pkg/helpers"
	_ "github.com/t1mon-ggg/gophkeeper/pkg/server/tls"
)

func TestTLSinit(t *testing.T) {
	require.True(t, helpers.FileExists("./ssl"))
	require.True(t, helpers.FileExists("./ssl/server.pem"))
	require.True(t, helpers.FileExists("./ssl/server.crt"))
	err := os.Remove("./ssl/server.pem")
	require.NoError(t, err)
	err = os.Remove("./ssl/server.crt")
	require.NoError(t, err)
	err = os.RemoveAll("./ssl")
	require.NoError(t, err)
}
