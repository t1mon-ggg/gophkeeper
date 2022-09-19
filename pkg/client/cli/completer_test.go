package cli

// import (
// 	"reflect"
// 	"sync"
// 	"testing"

// 	prompt "github.com/c-bata/go-prompt"

// 	"github.com/t1mon-ggg/gophkeeper/pkg/client/config"
// 	"github.com/t1mon-ggg/gophkeeper/pkg/client/openpgp"
// 	"github.com/t1mon-ggg/gophkeeper/pkg/client/storage"
// 	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
// 	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
// )

// func TestCLI_completer(t *testing.T) {
// 	zerolog.Initialize()
// 	type fields struct {
// 		wg      *sync.WaitGroup
// 		storage storage.Storage
// 		config  *config.Config
// 		crypto  *openpgp.OpenPGP
// 		logger  logging.Logger
// 	}
// 	type args struct {
// 		in prompt.Document
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   []prompt.Suggest
// 	}{
// 		{
// 			name: "init suggest",
// 			,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &CLI{
// 				wg:      tt.fields.wg,
// 				storage: tt.fields.storage,
// 				config:  tt.fields.config,
// 				crypto:  tt.fields.crypto,
// 				logger:  tt.fields.logger,
// 			}
// 			if got := c.completer(tt.args.in); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("CLI.completer() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
