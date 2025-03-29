// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: gameengineserver.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	GameEngineServer_GetCards_FullMethodName = "/com.sweetloveinyourheart.kittens.gameengines.GameEngineServer/GetCards"
)

// GameEngineServerClient is the client API for GameEngineServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GameEngineServerClient interface {
	// Get cards
	GetCards(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetCardsResponse, error)
}

type gameEngineServerClient struct {
	cc grpc.ClientConnInterface
}

func NewGameEngineServerClient(cc grpc.ClientConnInterface) GameEngineServerClient {
	return &gameEngineServerClient{cc}
}

func (c *gameEngineServerClient) GetCards(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetCardsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCardsResponse)
	err := c.cc.Invoke(ctx, GameEngineServer_GetCards_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GameEngineServerServer is the server API for GameEngineServer service.
// All implementations should embed UnimplementedGameEngineServerServer
// for forward compatibility.
type GameEngineServerServer interface {
	// Get cards
	GetCards(context.Context, *emptypb.Empty) (*GetCardsResponse, error)
}

// UnimplementedGameEngineServerServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGameEngineServerServer struct{}

func (UnimplementedGameEngineServerServer) GetCards(context.Context, *emptypb.Empty) (*GetCardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCards not implemented")
}
func (UnimplementedGameEngineServerServer) testEmbeddedByValue() {}

// UnsafeGameEngineServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GameEngineServerServer will
// result in compilation errors.
type UnsafeGameEngineServerServer interface {
	mustEmbedUnimplementedGameEngineServerServer()
}

func RegisterGameEngineServerServer(s grpc.ServiceRegistrar, srv GameEngineServerServer) {
	// If the following call pancis, it indicates UnimplementedGameEngineServerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GameEngineServer_ServiceDesc, srv)
}

func _GameEngineServer_GetCards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameEngineServerServer).GetCards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GameEngineServer_GetCards_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameEngineServerServer).GetCards(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// GameEngineServer_ServiceDesc is the grpc.ServiceDesc for GameEngineServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GameEngineServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.sweetloveinyourheart.kittens.gameengines.GameEngineServer",
	HandlerType: (*GameEngineServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCards",
			Handler:    _GameEngineServer_GetCards_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gameengineserver.proto",
}
