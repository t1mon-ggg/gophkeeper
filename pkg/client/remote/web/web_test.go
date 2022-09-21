package web

// import (
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/require"

// 	"github.com/t1mon-ggg/gophkeeper/pkg/client/config"
// 	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
// 	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
// 	"github.com/t1mon-ggg/gophkeeper/pkg/server/web"
// )

// const (
// 	pubKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
// Version: GopenPGP 2.4.10
// Comment: https://gopenpgp.org

// xjMEYyhjxhYJKwYBBAHaRw8BAQdABBGMYkg//P2HBC0JCUbjqhJy7I6X9RFm4XVd
// dKDNqAzNMWFjNWM2YjEwZDgyYjRkYmVhZGQ1OTdiZWYwMTA1OTcxIDx0ZXN0QGVt
// YWlsLmNvbT7CjAQTFggAPgUCYyhjxgmQjNtOjc05I/QWIQT+HnsCjM5BpqRmWPmM
// 206NzTkj9AIbAwIeAQIZAQMLCQcCFQgDFgACAiIBAABvtgEAk/WEMtwHzJhobY1J
// zI8BlwEWcGimU4tTZlRHny1qaEoBAKU7nA+fK6ZBscwqkUGT5OBHGAnVatOF3ioT
// CFU6cSANzjgEYyhjxhIKKwYBBAGXVQEFAQEHQFD9Bv7VAt2IOfO2DTRZDO5lZPad
// XBUdBd/6FXk0F1dUAwEKCcJ4BBgWCAAqBQJjKGPGCZCM206NzTkj9BYhBP4eewKM
// zkGmpGZY+YzbTo3NOSP0AhsMAADziAD/eKxXpURDFFGW2nH386BR0fz7Y3shDZmA
// CoxeB8h3ohgA/RtSpyonZVmeAhLWmJT2zKlgYERG8g11iYjeo48itewK
// =IyYM
// -----END PGP PUBLIC KEY BLOCK-----`

// 	pubkey2 = `-----BEGIN PGP PUBLIC KEY BLOCK-----
// Version: GopenPGP 2.4.10
// Comment: https://gopenpgp.org

// xjMEYyl+JxYJKwYBBAHaRw8BAQdAK8cD6k/5v5gVH5JQhGWTpaDjocDGPcgqgqAw
// LD5FzYLNTWFjNWM2YjEwZDgyYjRkYmVhZGQ1OTdiZWYwMTA1OTcxIDxhYzVjNmIx
// MGQ4MmI0ZGJlYWRkNTk3YmVmMDEwNTk3MUBsb2NhbGhvc3Q+wowEExYIAD4FAmMp
// ficJkORvv87KfTAYFiEE5ayzm55SeZ8yOQYg5G+/zsp9MBgCGwMCHgECGQEDCwkH
// AhUIAxYAAgIiAQAAMKUA/RM65zM/S6UYPbm4nS5W3k428MpXy0QBTvY3Ag26R334
// AP45s5BzoqZavdoubp/Q3A4KyImiTDrh++HTanFBgYWOBs44BGMpficSCisGAQQB
// l1UBBQEBB0C3xeHY581EWi77TbDB+YIONRPTL5FQE/7pl+V6fg3tOgMBCgnCeAQY
// FggAKgUCYyl+JwmQ5G+/zsp9MBgWIQTlrLObnlJ5nzI5BiDkb7/Oyn0wGAIbDAAA
// Z6QBAJq0nRRwhVkMVwJ7gjs8zPo/aI+vyuFRpy2cGlS3uBv1AP93LzeFVKBhJTfH
// EmGWnYPnJ3Aqj0izaWpoRpWSjYuVCw==
// =Tpcs
// -----END PGP PUBLIC KEY BLOCK-----`
// 	pubkey3 = `-----BEGIN PGP PUBLIC KEY BLOCK-----
// Version: GopenPGP 2.4.10
// Comment: https://gopenpgp.org

// xjMEYyl+JxYJKwYBBAHaRw8BAQdAK8cD6k/5v5gVH5JQhGWTpaDjocDGPcgqgqAw
// LD5FzYLNTWFjNWM2YjEwZDgyYjRkYmVhZGQ1OTdiZWYwMTA1OTcxIDxhYzVjNmIx
// MGQ4MmI0ZGJlYWRkNTk3YmVmMDEwfgbUBsb2NhbGhvc3Q+wowEExYIAD4FAmMp
// ficJkORvv87KfTAYFiEE5ayzm55SeZ8yOQYg5G+/zsp9MBgCGwMCHgECGQEDCwkH
// AhUIAxYAAgIiAQAAMKUA/RM65zM/S6UYPbm4nS5W3k428MpXy0QBTvY3Ag26R334
// AP45s5BzoqZavdoubp/Q3A4KyImiTDrh++HTanFBgYWOBs44BGMpficSCisGAQQB
// l1UBBQEBB0C3xeHY581EWi77TbDB+YIONRPTL5FQE/7pl+V6fg3tOgMBCgnCeAQY
// FggAKgUCYyl+JwmQ5G+/zsp9MBgWIQTlrLObnlJ5nzI5BiDkb7/Oyn0wGAIbDAAA
// Z6QBAJq0nRRwhVkMVwJ7gjs8zPo/aI+vyuFRpy2cGlS3uBv1AP93LzeFVKBhJTfH
// EmGWnYPnJ3Aqj0izaWpoRpWSjYuVCw==
// =Tpcs
// -----END PGP PUBLIC KEY BLOCK-----`

// 	content = `-----BEGIN PGP MESSAGE-----
// Version: GopenPGP 2.4.10
// Comment: https://gopenpgp.org

// wV4DAraayzAyyfcSAQdArz6Y1K/Bg53b+FprsbG4v4Ewn6+LvBNDhOX/Oxv7pWww
// 7c3fSrtLXFzf6kU2jqXLmDWGodGcZOD2h4pfUri+5ZyMLSy+qw4c/wvdsJ+i99Gz
// 0ukBV8RmiLa7fw8aFTrSeYAFLnjKg5LV0qnqZHJwiowXXPqUvNUEf62TBjMsJhWk
// Amr3Mtw3OTH2gg2Vxp/vhXqISPNzzRhFn1kDlXDXBfTWAs4PUPmEmY+8TZVS2C/y
// G5vOMSKgtrUCtyccfWwinyBG5/6iCxrhBC9QvE0kJ73UNxrTIL8PkkAAskM67clQ
// t4ZBFvdAYf5uN9HScXyOqePI73GJt0hQHiZtE04w33dRcGe11jLRereqoyrkshym
// Dpy4PgQJSqGXOwK0T95p+DUMhI99C2Kws29kXIUw7k+dmkyGkVuK6OCgtAwLsj+6
// dIn2syFI6uoOtTxPQ+YH5Y1yThMw44dXi+V/yG2GoZx+FvwTEXEPzx3Ut1z9Mkrd
// hNhuiIEaNnUnH0nd4GbkrMnr5Me/4WFqIMnf6VLaf1SrF4ObfHa8AsYEdnATq8HE
// Y9BNThrfRU+ZFQ7Rn4dbX4xhkbOKLEWy3SmQktm8gqsKtbE+VZu2EIWApW/pO2XL
// CjvKdSj+rix/xaLv8HDQ/KDpFrAXQp+jae4FM5tW8TQT21Yg90japPwrL/dA5iOb
// UthtMv7fB14ZuUuHYuAksTX6A1BU24jJmWkZlTI8elxwGgCy0Hn3r5oDUK3rVt6/
// iA6r/FpVDi+4JQmIfrXT+bqN9AQIuEWpLorJoR+2o11VeMAl3M1H7Z4/aEiYxdjP
// U6C3B67OmI0QlxDgP3tpYhf3QjY/3RbIHxYZ9fnDZTqglk6YIaRWZA69Llr1xxdS
// MCdi+JRRTSiA5kDt3nwuJVsL8/qUKwQstmC15eI1p9YFXykq1Qfwv7fsCJqUtfwB
// MwG1SRrQ7f3C2vXaFG+78TG1MuvMYlnD03p8UXBhtj8hmnZYLe+fPb7McIGs9IJB
// qo7KaZqTc/JdwvKDS28Wu/xDmGvk5g3BVPgs0KvV6FXVscb/FHZA1htMuwe5l9LC
// 3+2EjMSMkwvcXK7F+KimDkdZE2iKRENWAg==
// =iqvd
// -----END PGP MESSAGE-----`

// 	content2 = `-----BEGIN PGP MESSAGE-----
// Version: GopenPGP 2.4.10
// Comment: https://gopenpgp.org

// wV4DAraayzAyyfcSAQdArz6Y1K/Bg53b+FprsbG4v4Ewn6+LvBNDhOX/Oxv7pWww
// 7c3fSrtLXFzf6kU2jqXLmDWGodGcZOD2h4pfUri+5ZyMLSy+qw4c/wvdsJ+i99Gz
// 0ukBV8RmiLa7fw8aFTrSeYAFLnjKg5LV0qnqZHJwiowXXPqUvNUEf62TBjMsJhWk
// Amr3Mtw3OTH2gg2Vxp/vhXqISPNzzRhFn1kDlXDXBfTWAs4PUPmEmY+8TZVS2C/y
// G5vOMSKgtrUCtyccfWwinyBG5/6iCxrhBC9QvE0kJ73UNxrTIL8PkkAAskM67clQ
// t4ZBFvdAYf5uN9HScXyOqePI73GJt0hQHiZtE04w33dRcGe11jLRereqoyrkshym
// Dpy4PgQJSqGXOwK0T95p+DUMhI99C2Kws29kXIUw7k+dmkyGkVuK6OCgtAwLsj+6
// dIn2syFI6uoOtTxPQ+YH5Y1yThMw44dXi+V/yG2GoZx+FvwTEXEPzx3Ut1z9Mkrd
// hNhuiIEaNnUnH0nd4svdffbgkbOKLEWy3SmQktm8gqsKtbE+VZu2EIWApW/pO2XL
// CjvKdSj+rix/xaLv8HDQ/KDpFrAXQp+jae4FM5tW8TQT21Yg90japPwrL/dA5iOb
// UthtMv7fB14ZuUuHYuAksTX6A1BU24jJmWkZlTI8elxwGgCy0Hn3r5oDUK3rVt6/
// iA6r/FpVDi+4JQmIfrXT+bqN9AQIuEWpLorJoR+2o11VeMAl3M1H7Z4/aEiYxdjP
// U6C3B67OmI0QlxDgP3tpYhf3QjY/3RbIHxYZ9fnDZTqglk6YIaRWZA69Llr1xxdS
// MCdi+JRRTSiA5kDt3nwuJVsL8/qUKwQstmC15eI1p9YFXykq1Qfwv7fsCJqUtfwB
// MwG1SRrQ7f3C2vXaFG+78TG1MuvMYlnD03p8UXBhtj8hmnZYLe+fPb7McIGs9IJB
// qo7KaZqTc/JdwvKDS28Wu/xDmGvk5g3BVPgs0KvV6FXVscb/FHZA1htMuwe5l9LC
// 3+2EjMSMkwvcXK7F+KimDkdZE2iKRENWAg==
// =iqvd
// -----END PGP MESSAGE-----`
// )

// func startWeb() *web.Server {
// 	s := web.New()
// 	wg := sync.WaitGroup{}
// 	go s.Start(&wg)
// 	return s
// }

// func TestNew(t *testing.T) {
// 	cfg := config.New()
// 	// zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// }

// func TestSignup(t *testing.T) {
// 	cfg := config.New()
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	t.Run("create new resty client", func(t *testing.T) {
// 		err := c.Register("test", "password", pubKey)
// 		require.NoError(t, err)
// 	})
// }

// func TestLogin(t *testing.T) {
// 	cfg := config.New()
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	t.Run("Login resty client", func(t *testing.T) {
// 		err := c.Login("test", "password", pubKey)
// 		require.NoError(t, err)
// 	})
// }

// func TestPush(t *testing.T) {
// 	cfg := config.New()
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	err := c.Login("test", "password", pubKey)
// 	require.NoError(t, err)
// 	t.Run("save some content", func(t *testing.T) {
// 		err := c.Push(content, "eoijrnv9384gj")
// 		require.NoError(t, err)
// 	})
// }

// func TestLogs(t *testing.T) {
// 	cfg := config.New()
// 	log := zerolog.New().WithPrefix("web-test")
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	err := c.Login("test", "password", pubKey)
// 	require.NoError(t, err)
// 	t.Run("get some logs", func(t *testing.T) {
// 		actions, err := c.GetLogs()
// 		require.NoError(t, err)
// 		require.NotEmpty(t, actions)
// 		log.Tracef("%+v", nil, actions)
// 	})
// }

// func TestPull(t *testing.T) {
// 	cfg := config.New()
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	t.Run("get some content", func(t *testing.T) {
// 		secret, err := c.Pull("eoijrnv9384gj")
// 		require.NoError(t, err)
// 		require.Equal(t, []byte(content), secret)
// 	})
// }

// func TestVersion(t *testing.T) {
// 	cfg := config.New()
// 	log := zerolog.New()
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	t.Run("get versions", func(t *testing.T) {
// 		versions, err := c.Versions()
// 		require.NoError(t, err)
// 		require.NotEmpty(t, versions)
// 		log.Tracef("%+v", nil, versions)
// 	})
// }

// func TestAddPGP(t *testing.T) {
// 	cfg := config.New()
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	t.Run("get pgp list", func(t *testing.T) {
// 		err := c.AddPGP(pubkey2)
// 		require.NoError(t, err)
// 		list, err := c.ListPGP()
// 		require.NoError(t, err)
// 		var found bool
// 		for _, key := range list {
// 			if key.Publickey == pubkey2 {
// 				found = true
// 			}
// 		}
// 		require.True(t, found)
// 	})
// }

// func TestConfirmPGP(t *testing.T) {
// 	cfg := config.New()
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	t.Run("get pgp list", func(t *testing.T) {
// 		err := c.ConfirmPGP(pubkey2)
// 		require.NoError(t, err)
// 		list, err := c.ListPGP()
// 		require.NoError(t, err)
// 		for _, key := range list {
// 			if key.Publickey == pubkey2 {
// 				require.True(t, key.Confirmed)
// 			}
// 		}

// 	})
// }

// func TestRevokePGP(t *testing.T) {
// 	cfg := config.New()
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	t.Run("get pgp list", func(t *testing.T) {
// 		err := c.RevokePGP(pubkey2)
// 		require.NoError(t, err)
// 		list, err := c.ListPGP()
// 		require.NoError(t, err)
// 		var found bool
// 		for _, key := range list {
// 			if key.Publickey == pubkey2 {
// 				found = true
// 			}
// 		}
// 		require.False(t, found)

// 	})
// }

// func TestListPGP(t *testing.T) {
// 	cfg := config.New()
// 	log := zerolog.New()
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	t.Run("get pgp list", func(t *testing.T) {
// 		keys, err := c.ListPGP()
// 		require.NoError(t, err)
// 		require.NotEmpty(t, keys)
// 		log.Tracef("%+v", nil, keys)
// 	})
// }

// func TestWSClient(t *testing.T) {
// 	cfg := config.New()
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	err := c.Login("test", "password", pubKey)
// 	require.NoError(t, err)
// 	t.Run("websocket create", func(t *testing.T) {
// 		go func() {
// 			time.Sleep(3 * time.Second)
// 			err := c.AddPGP(pubkey3)
// 			require.NoError(t, err)
// 			time.Sleep(3 * time.Second)
// 			err = c.Push(content2, "29385608746501983456")
// 			require.NoError(t, err)
// 			time.Sleep(3 * time.Second)
// 			c.Close()
// 		}()
// 		err := c.NewStream()
// 		require.NoError(t, err)

// 	})
// }

// func TestDelete(t *testing.T) {
// 	cfg := config.New()
// 	zerolog.New().SetLevel(logging.TraceLevel)
// 	cfg.RemoteHTTP = "https://127.0.0.1:8443"
// 	c := New()
// 	require.NotNil(t, c)
// 	err := c.Login("test", "password", pubKey)
// 	require.NoError(t, err)
// 	t.Run("Delete user", func(t *testing.T) {
// 		err := c.Delete()
// 		require.NoError(t, err)

// 	})
// }
