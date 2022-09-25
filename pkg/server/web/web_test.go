package web

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
	"github.com/t1mon-ggg/gophkeeper/pkg/models"
	mockStorage "github.com/t1mon-ggg/gophkeeper/pkg/server/storage/mock_storage"
	ssl "github.com/t1mon-ggg/gophkeeper/pkg/server/tls"
	"github.com/t1mon-ggg/gophkeeper/pkg/server/web/static"
)

func clean(t *testing.T) {
	list, err := os.ReadDir("./ssl")
	require.NoError(t, err)
	for _, item := range list {
		err := os.Remove(fmt.Sprintf("./ssl/%s", item.Name()))
		require.NoError(t, err)
	}
	err = os.RemoveAll("./ssl")
	require.NoError(t, err)
	os.Unsetenv("WEB_ADDRESS")
}

func client(t *testing.T) *resty.Client {
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := resty.New()
	client.SetCookieJar(jar)
	client.SetTLSClientConfig(&tls.Config{
		Rand:               rand.Reader,
		InsecureSkipVerify: true,
	})
	return client
}

func login() (*resty.Response, error) {
	return nil, nil
}

func TestServer(t *testing.T) {
	rbi, err := rand.Int(rand.Reader, big.NewInt(60000))
	require.NoError(t, err)
	port := int(rbi.Int64())
	if port < 10000 {
		port += 10000
	}
	webbind := fmt.Sprintf("127.0.0.1:%d", port)
	os.Setenv("WEB_ADDRESS", webbind)
	ssl.Prepare()
	defer clean(t)
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	ip := net.ParseIP("127.0.0.1")
	gomock.InOrder(
		//signin
		db.EXPECT().SignIn(models.User{Username: "test", Password: "test", PublicKey: "test"}, ip).Return(nil),
		db.EXPECT().ListPGP("test", ip).Return([]models.PGP{{Date: time.Now(), Publickey: "test", Confirmed: true}}, nil),
	)

	wg := new(sync.WaitGroup)
	s := &Server{
		echo: echo.New(),
		log:  zerolog.New().WithPrefix("web-server-test"),
		db:   db,
		sig:  make(chan struct{}),
		msg:  make(map[string]chan models.Message),
	}
	s.applyMiddlewares()
	static.ApplyStatic(s.echo)
	s.createRouter()
	go func() {
		err := s.Start(wg)
		require.NoError(t, err)
	}()
	time.Sleep(2 * time.Second)
	t.Run("test login", func(t *testing.T) {
		c := client(t)
		c.SetBaseURL(fmt.Sprintf("https://%s", webbind))

		type test struct {
			name string
			body models.User

			want bool
			code int
		}

		tests := []test{
			{
				name: "successfull login",
				body: models.User{
					Username:  "test",
					Password:  "test",
					PublicKey: "test",
				},

				want: false,
				code: http.StatusOK,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var got bool
				response, err := c.R().
					SetHeader("Content-Type", "application/json").
					SetBody(tt.body).
					Post("api/v1/signin")
				if err != nil {
					got = true
				}
				require.Equal(t, tt.want, got)
				if !tt.want {
					require.Equal(t, tt.code, response.StatusCode())
					require.NotEmpty(t, response.Cookies())
					var token bool
					for _, cookie := range response.Cookies() {
						if cookie.Name == "token" {
							token = true
							require.NotEmpty(t, cookie.Value)
						}
					}
					require.True(t, token)
				}
			})
		}

	})
	err = s.Stop()
	require.NoError(t, err)
	s.wg.Wait()
}
