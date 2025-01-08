package interceptors

import (
	"connectrpc.com/connect"
	auth_interceptors "github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors/auth"
)

func CommonConnectInterceptors(serviceName string, signingKey string, authFunc auth_interceptors.AuthFunc, authOpts ...auth_interceptors.Option) []connect.Interceptor {
	if authFunc == nil {
		authFunc = ConnectAuthHandler(signingKey)
	}

	return []connect.Interceptor{
		auth_interceptors.NewAuthInterceptor(authFunc, authOpts...),
	}
}
