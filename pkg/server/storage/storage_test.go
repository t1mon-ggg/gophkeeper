package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
)

var testDSN = "postgresql://gophkeeper:gophkeeper@127.0.0.1:5432/gophkeeper"

func TestNew(t *testing.T) {
	logger := zerolog.New()
	type args struct {
		dsn    string
		logger logging.Logger
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid dsn",
			args: args{
				logger: logger,
				dsn:    testDSN,
			},
			want: true,
		},
		{
			name: "invalid dsn",
			args: args{
				logger: logger,
				dsn:    "postgresql://gophkeeper:gophkeeper@127.0.0.1:6432/gophkeeper",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := true
			_, err := New(tt.args.dsn, tt.args.logger)
			if err != nil {
				e = false
			}
			require.Equal(t, tt.want, e)
		})
	}
}

func TestStorage_log(t *testing.T) {
	logger := zerolog.New()
	s, err := New(testDSN, logger)
	require.NoError(t, err)
	t.Run("logger test", func(t *testing.T) {
		require.NotNil(t, s.log())
	})
}

func TestStorage_SignUp(t *testing.T) {
	logger := zerolog.New()
	s, err := New(testDSN, logger)
	require.NoError(t, err)
	type args struct {
		u string
		p string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "create first user",
			args: args{
				u: "first",
				p: "password",
			},
		},
		{
			name: "create second user",
			args: args{
				u: "second",
				p: "password",
			},
		},
	}
	defer func() {
		for _, user := range tests {
			err := s.DeleteUser(user.args.u)
			require.NoError(t, err)
		}
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.SignUp(tt.args.u, tt.args.p)
			require.NoError(t, err)
		})
	}
}

// func Test_SignIn(t *testing.T) {
// 	logger := zerolog.New()
// 	s, err := New(testDSN, logger)
// 	require.NoError(t, err)
// 	type args struct {
// 		u  string
// 		p  string
// 		tp string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bool
// 	}{
// 		{
// 			name: "normal login",
// 			args: args{
// 				u:  "first",
// 				p:  "password",
// 				tp: "password",
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "wrong password",
// 			args: args{
// 				u:  "second",
// 				p:  "password",
// 				tp: "wrong_password",
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "wrong username",
// 			args: args{
// 				u:  "third",
// 				p:  "password",
// 				tp: "password",
// 			},
// 			want: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		if tt.args.u == "third" {
// 			continue
// 		}
// 		err := s.SignUp(tt.args.u, tt.args.p)
// 		require.NoError(t, err)
// 	}
// 	defer func() {
// 		for _, user := range tests {
// 			if user.args.u == "third" {
// 				continue
// 			}
// 			err := s.DeleteUser(user.args.u)
// 			require.NoError(t, err)
// 		}
// 	}()
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			login := true
// 			err := s.SignIn(tt.args.u, tt.args.tp)
// 			if err != nil {
// 				login = false
// 				s.log().Trace(err)
// 			}
// 			require.Equal(t, tt.want, login)
// 		})
// 	}
// }

func Test_DeleteUser(t *testing.T) {
	logger := zerolog.New()
	s, err := New(testDSN, logger)
	require.NoError(t, err)
	type args struct {
		u string
		p string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "delete existing user",
			args: args{
				u: "second",
				p: "password",
			},
			want: true,
		},
		{
			name: "delete not existing user",
			args: args{
				u: "third",
				p: "password",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		if tt.args.u == "third" {
			continue
		}
		err := s.SignUp(tt.args.u, tt.args.p)
		require.NoError(t, err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := true
			err := s.DeleteUser(tt.args.u)
			if err != nil {
				got = false
				s.log().Trace(err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_Ping(t *testing.T) {
	logger := zerolog.New()
	s, err := New(testDSN, logger)
	require.NoError(t, err)
	t.Run("test ping function", func(t *testing.T) {
		err := s.Ping()
		require.NoError(t, err)
	})
}

func Test_Close(t *testing.T) {
	logger := zerolog.New()
	s, err := New(testDSN, logger)
	require.NoError(t, err)
	t.Run("test close function", func(t *testing.T) {
		s.Close()
		err := s.Ping()
		require.Error(t, err)
		s.log().Trace(err)
	})
}
