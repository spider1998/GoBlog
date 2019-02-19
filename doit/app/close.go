package app

import (
	"golang.org/x/net/context"
	"net/http"
	"sync"
)

var (
	serverMutex sync.Mutex
	server      *http.Server
)

/*-----关闭服务-----*/
func ServerClose() error {
	Logger.Info().Msg("try to shutdown http server.")
	defer Logger.Info().Msg("already shutdown http server.")
	serverMutex.Lock()
	defer serverMutex.Unlock()
	if server != nil {
		return server.Shutdown(context.Background())
	}
	return nil
}
