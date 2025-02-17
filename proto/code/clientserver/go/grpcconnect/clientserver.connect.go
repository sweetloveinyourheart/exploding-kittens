// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: clientserver.proto

package grpcconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	_go "github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// ClientServerName is the fully-qualified name of the ClientServer service.
	ClientServerName = "com.sweetloveinyourheart.kittens.clients.ClientServer"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ClientServerCreateNewGuestUserProcedure is the fully-qualified name of the ClientServer's
	// CreateNewGuestUser RPC.
	ClientServerCreateNewGuestUserProcedure = "/com.sweetloveinyourheart.kittens.clients.ClientServer/CreateNewGuestUser"
	// ClientServerGuestLoginProcedure is the fully-qualified name of the ClientServer's GuestLogin RPC.
	ClientServerGuestLoginProcedure = "/com.sweetloveinyourheart.kittens.clients.ClientServer/GuestLogin"
	// ClientServerGetPlayerProfileProcedure is the fully-qualified name of the ClientServer's
	// GetPlayerProfile RPC.
	ClientServerGetPlayerProfileProcedure = "/com.sweetloveinyourheart.kittens.clients.ClientServer/GetPlayerProfile"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	clientServerServiceDescriptor                  = _go.File_clientserver_proto.Services().ByName("ClientServer")
	clientServerCreateNewGuestUserMethodDescriptor = clientServerServiceDescriptor.Methods().ByName("CreateNewGuestUser")
	clientServerGuestLoginMethodDescriptor         = clientServerServiceDescriptor.Methods().ByName("GuestLogin")
	clientServerGetPlayerProfileMethodDescriptor   = clientServerServiceDescriptor.Methods().ByName("GetPlayerProfile")
)

// ClientServerClient is a client for the com.sweetloveinyourheart.kittens.clients.ClientServer
// service.
type ClientServerClient interface {
	CreateNewGuestUser(context.Context, *connect.Request[_go.CreateNewGuestUserRequest]) (*connect.Response[_go.CreateNewGuestUserResponse], error)
	GuestLogin(context.Context, *connect.Request[_go.GuestLoginRequest]) (*connect.Response[_go.GuestLoginResponse], error)
	GetPlayerProfile(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[_go.PlayerProfileResponse], error)
}

// NewClientServerClient constructs a client for the
// com.sweetloveinyourheart.kittens.clients.ClientServer service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewClientServerClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ClientServerClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &clientServerClient{
		createNewGuestUser: connect.NewClient[_go.CreateNewGuestUserRequest, _go.CreateNewGuestUserResponse](
			httpClient,
			baseURL+ClientServerCreateNewGuestUserProcedure,
			connect.WithSchema(clientServerCreateNewGuestUserMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		guestLogin: connect.NewClient[_go.GuestLoginRequest, _go.GuestLoginResponse](
			httpClient,
			baseURL+ClientServerGuestLoginProcedure,
			connect.WithSchema(clientServerGuestLoginMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getPlayerProfile: connect.NewClient[emptypb.Empty, _go.PlayerProfileResponse](
			httpClient,
			baseURL+ClientServerGetPlayerProfileProcedure,
			connect.WithSchema(clientServerGetPlayerProfileMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// clientServerClient implements ClientServerClient.
type clientServerClient struct {
	createNewGuestUser *connect.Client[_go.CreateNewGuestUserRequest, _go.CreateNewGuestUserResponse]
	guestLogin         *connect.Client[_go.GuestLoginRequest, _go.GuestLoginResponse]
	getPlayerProfile   *connect.Client[emptypb.Empty, _go.PlayerProfileResponse]
}

// CreateNewGuestUser calls
// com.sweetloveinyourheart.kittens.clients.ClientServer.CreateNewGuestUser.
func (c *clientServerClient) CreateNewGuestUser(ctx context.Context, req *connect.Request[_go.CreateNewGuestUserRequest]) (*connect.Response[_go.CreateNewGuestUserResponse], error) {
	return c.createNewGuestUser.CallUnary(ctx, req)
}

// GuestLogin calls com.sweetloveinyourheart.kittens.clients.ClientServer.GuestLogin.
func (c *clientServerClient) GuestLogin(ctx context.Context, req *connect.Request[_go.GuestLoginRequest]) (*connect.Response[_go.GuestLoginResponse], error) {
	return c.guestLogin.CallUnary(ctx, req)
}

// GetPlayerProfile calls com.sweetloveinyourheart.kittens.clients.ClientServer.GetPlayerProfile.
func (c *clientServerClient) GetPlayerProfile(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[_go.PlayerProfileResponse], error) {
	return c.getPlayerProfile.CallUnary(ctx, req)
}

// ClientServerHandler is an implementation of the
// com.sweetloveinyourheart.kittens.clients.ClientServer service.
type ClientServerHandler interface {
	CreateNewGuestUser(context.Context, *connect.Request[_go.CreateNewGuestUserRequest]) (*connect.Response[_go.CreateNewGuestUserResponse], error)
	GuestLogin(context.Context, *connect.Request[_go.GuestLoginRequest]) (*connect.Response[_go.GuestLoginResponse], error)
	GetPlayerProfile(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[_go.PlayerProfileResponse], error)
}

// NewClientServerHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewClientServerHandler(svc ClientServerHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	clientServerCreateNewGuestUserHandler := connect.NewUnaryHandler(
		ClientServerCreateNewGuestUserProcedure,
		svc.CreateNewGuestUser,
		connect.WithSchema(clientServerCreateNewGuestUserMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	clientServerGuestLoginHandler := connect.NewUnaryHandler(
		ClientServerGuestLoginProcedure,
		svc.GuestLogin,
		connect.WithSchema(clientServerGuestLoginMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	clientServerGetPlayerProfileHandler := connect.NewUnaryHandler(
		ClientServerGetPlayerProfileProcedure,
		svc.GetPlayerProfile,
		connect.WithSchema(clientServerGetPlayerProfileMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/com.sweetloveinyourheart.kittens.clients.ClientServer/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ClientServerCreateNewGuestUserProcedure:
			clientServerCreateNewGuestUserHandler.ServeHTTP(w, r)
		case ClientServerGuestLoginProcedure:
			clientServerGuestLoginHandler.ServeHTTP(w, r)
		case ClientServerGetPlayerProfileProcedure:
			clientServerGetPlayerProfileHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedClientServerHandler returns CodeUnimplemented from all methods.
type UnimplementedClientServerHandler struct{}

func (UnimplementedClientServerHandler) CreateNewGuestUser(context.Context, *connect.Request[_go.CreateNewGuestUserRequest]) (*connect.Response[_go.CreateNewGuestUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("com.sweetloveinyourheart.kittens.clients.ClientServer.CreateNewGuestUser is not implemented"))
}

func (UnimplementedClientServerHandler) GuestLogin(context.Context, *connect.Request[_go.GuestLoginRequest]) (*connect.Response[_go.GuestLoginResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("com.sweetloveinyourheart.kittens.clients.ClientServer.GuestLogin is not implemented"))
}

func (UnimplementedClientServerHandler) GetPlayerProfile(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[_go.PlayerProfileResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("com.sweetloveinyourheart.kittens.clients.ClientServer.GetPlayerProfile is not implemented"))
}
