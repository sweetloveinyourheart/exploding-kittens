// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: clientserver.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type User struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	FullName      string                 `protobuf:"bytes,3,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Status        int32                  `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_clientserver_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_clientserver_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_clientserver_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *User) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *User) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

// Message for creating a new guest user
type CreateNewGuestUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`                 // Required: Username of the guest user
	FullName      string                 `protobuf:"bytes,2,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"` // Required: Full name of the guest user
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateNewGuestUserRequest) Reset() {
	*x = CreateNewGuestUserRequest{}
	mi := &file_clientserver_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateNewGuestUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateNewGuestUserRequest) ProtoMessage() {}

func (x *CreateNewGuestUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_clientserver_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateNewGuestUserRequest.ProtoReflect.Descriptor instead.
func (*CreateNewGuestUserRequest) Descriptor() ([]byte, []int) {
	return file_clientserver_proto_rawDescGZIP(), []int{1}
}

func (x *CreateNewGuestUserRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *CreateNewGuestUserRequest) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

type CreateNewGuestUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"` // The user basic info
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateNewGuestUserResponse) Reset() {
	*x = CreateNewGuestUserResponse{}
	mi := &file_clientserver_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateNewGuestUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateNewGuestUserResponse) ProtoMessage() {}

func (x *CreateNewGuestUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_clientserver_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateNewGuestUserResponse.ProtoReflect.Descriptor instead.
func (*CreateNewGuestUserResponse) Descriptor() ([]byte, []int) {
	return file_clientserver_proto_rawDescGZIP(), []int{2}
}

func (x *CreateNewGuestUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

// Message for guest login
type GuestLoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // Required: UUID of the guest user
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GuestLoginRequest) Reset() {
	*x = GuestLoginRequest{}
	mi := &file_clientserver_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GuestLoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuestLoginRequest) ProtoMessage() {}

func (x *GuestLoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_clientserver_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuestLoginRequest.ProtoReflect.Descriptor instead.
func (*GuestLoginRequest) Descriptor() ([]byte, []int) {
	return file_clientserver_proto_rawDescGZIP(), []int{3}
}

func (x *GuestLoginRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GuestLoginResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // The database id for this user (UUID).
	Token         string                 `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`                 // The session token for this user.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GuestLoginResponse) Reset() {
	*x = GuestLoginResponse{}
	mi := &file_clientserver_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GuestLoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuestLoginResponse) ProtoMessage() {}

func (x *GuestLoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_clientserver_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuestLoginResponse.ProtoReflect.Descriptor instead.
func (*GuestLoginResponse) Descriptor() ([]byte, []int) {
	return file_clientserver_proto_rawDescGZIP(), []int{4}
}

func (x *GuestLoginResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GuestLoginResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// Message for player profile
type PlayerProfileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PlayerProfileResponse) Reset() {
	*x = PlayerProfileResponse{}
	mi := &file_clientserver_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PlayerProfileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerProfileResponse) ProtoMessage() {}

func (x *PlayerProfileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_clientserver_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerProfileResponse.ProtoReflect.Descriptor instead.
func (*PlayerProfileResponse) Descriptor() ([]byte, []int) {
	return file_clientserver_proto_rawDescGZIP(), []int{5}
}

func (x *PlayerProfileResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type Lobby struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	LobbyId       string                 `protobuf:"bytes,1,opt,name=lobby_id,json=lobbyId,proto3" json:"lobby_id,omitempty"`
	LobbyCode     string                 `protobuf:"bytes,2,opt,name=lobby_code,json=lobbyCode,proto3" json:"lobby_code,omitempty"`
	LobbyName     string                 `protobuf:"bytes,3,opt,name=lobby_name,json=lobbyName,proto3" json:"lobby_name,omitempty"`
	HostUserId    string                 `protobuf:"bytes,4,opt,name=host_user_id,json=hostUserId,proto3" json:"host_user_id,omitempty"`
	Participants  []string               `protobuf:"bytes,5,rep,name=participants,proto3" json:"participants,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Lobby) Reset() {
	*x = Lobby{}
	mi := &file_clientserver_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Lobby) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Lobby) ProtoMessage() {}

func (x *Lobby) ProtoReflect() protoreflect.Message {
	mi := &file_clientserver_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Lobby.ProtoReflect.Descriptor instead.
func (*Lobby) Descriptor() ([]byte, []int) {
	return file_clientserver_proto_rawDescGZIP(), []int{6}
}

func (x *Lobby) GetLobbyId() string {
	if x != nil {
		return x.LobbyId
	}
	return ""
}

func (x *Lobby) GetLobbyCode() string {
	if x != nil {
		return x.LobbyCode
	}
	return ""
}

func (x *Lobby) GetLobbyName() string {
	if x != nil {
		return x.LobbyName
	}
	return ""
}

func (x *Lobby) GetHostUserId() string {
	if x != nil {
		return x.HostUserId
	}
	return ""
}

func (x *Lobby) GetParticipants() []string {
	if x != nil {
		return x.Participants
	}
	return nil
}

type CreateLobbyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	LobbyName     string                 `protobuf:"bytes,1,opt,name=lobby_name,json=lobbyName,proto3" json:"lobby_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateLobbyRequest) Reset() {
	*x = CreateLobbyRequest{}
	mi := &file_clientserver_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateLobbyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLobbyRequest) ProtoMessage() {}

func (x *CreateLobbyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_clientserver_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLobbyRequest.ProtoReflect.Descriptor instead.
func (*CreateLobbyRequest) Descriptor() ([]byte, []int) {
	return file_clientserver_proto_rawDescGZIP(), []int{7}
}

func (x *CreateLobbyRequest) GetLobbyName() string {
	if x != nil {
		return x.LobbyName
	}
	return ""
}

type CreateLobbyResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	LobbyId       string                 `protobuf:"bytes,1,opt,name=lobby_id,json=lobbyId,proto3" json:"lobby_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateLobbyResponse) Reset() {
	*x = CreateLobbyResponse{}
	mi := &file_clientserver_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateLobbyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLobbyResponse) ProtoMessage() {}

func (x *CreateLobbyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_clientserver_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLobbyResponse.ProtoReflect.Descriptor instead.
func (*CreateLobbyResponse) Descriptor() ([]byte, []int) {
	return file_clientserver_proto_rawDescGZIP(), []int{8}
}

func (x *CreateLobbyResponse) GetLobbyId() string {
	if x != nil {
		return x.LobbyId
	}
	return ""
}

type GetLobbyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	LobbyId       string                 `protobuf:"bytes,1,opt,name=lobby_id,json=lobbyId,proto3" json:"lobby_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetLobbyRequest) Reset() {
	*x = GetLobbyRequest{}
	mi := &file_clientserver_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetLobbyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLobbyRequest) ProtoMessage() {}

func (x *GetLobbyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_clientserver_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLobbyRequest.ProtoReflect.Descriptor instead.
func (*GetLobbyRequest) Descriptor() ([]byte, []int) {
	return file_clientserver_proto_rawDescGZIP(), []int{9}
}

func (x *GetLobbyRequest) GetLobbyId() string {
	if x != nil {
		return x.LobbyId
	}
	return ""
}

type GetLobbyReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Lobby         *Lobby                 `protobuf:"bytes,1,opt,name=lobby,proto3" json:"lobby,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetLobbyReply) Reset() {
	*x = GetLobbyReply{}
	mi := &file_clientserver_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetLobbyReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLobbyReply) ProtoMessage() {}

func (x *GetLobbyReply) ProtoReflect() protoreflect.Message {
	mi := &file_clientserver_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLobbyReply.ProtoReflect.Descriptor instead.
func (*GetLobbyReply) Descriptor() ([]byte, []int) {
	return file_clientserver_proto_rawDescGZIP(), []int{10}
}

func (x *GetLobbyReply) GetLobby() *Lobby {
	if x != nil {
		return x.Lobby
	}
	return nil
}

var File_clientserver_proto protoreflect.FileDescriptor

var file_clientserver_proto_rawDesc = string([]byte{
	0x0a, 0x12, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x28, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c,
	0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x6b,
	0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x70, 0x0a, 0x04, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x75, 0x6c, 0x6c,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6c,
	0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x54, 0x0a,
	0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x77, 0x47, 0x75, 0x65, 0x73, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x75, 0x6c, 0x6c, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0x60, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x77,
	0x47, 0x75, 0x65, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x42, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2e, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69,
	0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x6b, 0x69, 0x74, 0x74, 0x65,
	0x6e, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x2c, 0x0a, 0x11, 0x47, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x22, 0x43, 0x0a, 0x12, 0x47, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x5b, 0x0a, 0x15, 0x50, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x42, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2e, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69,
	0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x6b, 0x69, 0x74, 0x74, 0x65,
	0x6e, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0xa6, 0x01, 0x0a, 0x05, 0x4c, 0x6f, 0x62, 0x62, 0x79, 0x12,
	0x19, 0x0a, 0x08, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f,
	0x62, 0x62, 0x79, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x6c, 0x6f, 0x62, 0x62, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x62,
	0x62, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6c,
	0x6f, 0x62, 0x62, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0c, 0x68, 0x6f, 0x73, 0x74,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x68, 0x6f, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x61,
	0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0c, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x22, 0x33,
	0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x62, 0x62, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0x30, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x62,
	0x62, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x6f,
	0x62, 0x62, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x6f,
	0x62, 0x62, 0x79, 0x49, 0x64, 0x22, 0x2c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x62, 0x62,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x6f, 0x62, 0x62,
	0x79, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x6f, 0x62, 0x62,
	0x79, 0x49, 0x64, 0x22, 0x56, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x62, 0x62, 0x79, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x05, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c,
	0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x6b,
	0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x4c,
	0x6f, 0x62, 0x62, 0x79, 0x52, 0x05, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x32, 0xba, 0x05, 0x0a, 0x0c,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x9f, 0x01, 0x0a,
	0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x77, 0x47, 0x75, 0x65, 0x73, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x43, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c,
	0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x6b,
	0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x77, 0x47, 0x75, 0x65, 0x73, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x44, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73,
	0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65,
	0x61, 0x72, 0x74, 0x2e, 0x6b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x77, 0x47, 0x75, 0x65,
	0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x87,
	0x01, 0x0a, 0x0a, 0x47, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x3b, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79,
	0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x6b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73,
	0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x47, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3c, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72,
	0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x6b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x47, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6b, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x50,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x3f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74,
	0x6c, 0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e,
	0x6b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x8a, 0x01, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x4c, 0x6f, 0x62, 0x62, 0x79, 0x12, 0x3c, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65,
	0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74,
	0x2e, 0x6b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x62, 0x62, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x3d, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c,
	0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x6b,
	0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x62, 0x62, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x83, 0x01, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4c, 0x6f, 0x62,
	0x62, 0x79, 0x12, 0x39, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f,
	0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x6b, 0x69,
	0x74, 0x74, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x47, 0x65,
	0x74, 0x4c, 0x6f, 0x62, 0x62, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x37, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79,
	0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x6b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73,
	0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x62, 0x62,
	0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x30, 0x01, 0x42, 0x53, 0x5a, 0x51, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65,
	0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2f, 0x65, 0x78, 0x70, 0x6c,
	0x6f, 0x64, 0x69, 0x6e, 0x67, 0x2d, 0x6b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x67, 0x6f, 0x3b, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_clientserver_proto_rawDescOnce sync.Once
	file_clientserver_proto_rawDescData []byte
)

func file_clientserver_proto_rawDescGZIP() []byte {
	file_clientserver_proto_rawDescOnce.Do(func() {
		file_clientserver_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_clientserver_proto_rawDesc), len(file_clientserver_proto_rawDesc)))
	})
	return file_clientserver_proto_rawDescData
}

var file_clientserver_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_clientserver_proto_goTypes = []any{
	(*User)(nil),                       // 0: com.sweetloveinyourheart.kittens.clients.User
	(*CreateNewGuestUserRequest)(nil),  // 1: com.sweetloveinyourheart.kittens.clients.CreateNewGuestUserRequest
	(*CreateNewGuestUserResponse)(nil), // 2: com.sweetloveinyourheart.kittens.clients.CreateNewGuestUserResponse
	(*GuestLoginRequest)(nil),          // 3: com.sweetloveinyourheart.kittens.clients.GuestLoginRequest
	(*GuestLoginResponse)(nil),         // 4: com.sweetloveinyourheart.kittens.clients.GuestLoginResponse
	(*PlayerProfileResponse)(nil),      // 5: com.sweetloveinyourheart.kittens.clients.PlayerProfileResponse
	(*Lobby)(nil),                      // 6: com.sweetloveinyourheart.kittens.clients.Lobby
	(*CreateLobbyRequest)(nil),         // 7: com.sweetloveinyourheart.kittens.clients.CreateLobbyRequest
	(*CreateLobbyResponse)(nil),        // 8: com.sweetloveinyourheart.kittens.clients.CreateLobbyResponse
	(*GetLobbyRequest)(nil),            // 9: com.sweetloveinyourheart.kittens.clients.GetLobbyRequest
	(*GetLobbyReply)(nil),              // 10: com.sweetloveinyourheart.kittens.clients.GetLobbyReply
	(*emptypb.Empty)(nil),              // 11: google.protobuf.Empty
}
var file_clientserver_proto_depIdxs = []int32{
	0,  // 0: com.sweetloveinyourheart.kittens.clients.CreateNewGuestUserResponse.user:type_name -> com.sweetloveinyourheart.kittens.clients.User
	0,  // 1: com.sweetloveinyourheart.kittens.clients.PlayerProfileResponse.user:type_name -> com.sweetloveinyourheart.kittens.clients.User
	6,  // 2: com.sweetloveinyourheart.kittens.clients.GetLobbyReply.lobby:type_name -> com.sweetloveinyourheart.kittens.clients.Lobby
	1,  // 3: com.sweetloveinyourheart.kittens.clients.ClientServer.CreateNewGuestUser:input_type -> com.sweetloveinyourheart.kittens.clients.CreateNewGuestUserRequest
	3,  // 4: com.sweetloveinyourheart.kittens.clients.ClientServer.GuestLogin:input_type -> com.sweetloveinyourheart.kittens.clients.GuestLoginRequest
	11, // 5: com.sweetloveinyourheart.kittens.clients.ClientServer.GetPlayerProfile:input_type -> google.protobuf.Empty
	7,  // 6: com.sweetloveinyourheart.kittens.clients.ClientServer.CreateLobby:input_type -> com.sweetloveinyourheart.kittens.clients.CreateLobbyRequest
	9,  // 7: com.sweetloveinyourheart.kittens.clients.ClientServer.StreamLobby:input_type -> com.sweetloveinyourheart.kittens.clients.GetLobbyRequest
	2,  // 8: com.sweetloveinyourheart.kittens.clients.ClientServer.CreateNewGuestUser:output_type -> com.sweetloveinyourheart.kittens.clients.CreateNewGuestUserResponse
	4,  // 9: com.sweetloveinyourheart.kittens.clients.ClientServer.GuestLogin:output_type -> com.sweetloveinyourheart.kittens.clients.GuestLoginResponse
	5,  // 10: com.sweetloveinyourheart.kittens.clients.ClientServer.GetPlayerProfile:output_type -> com.sweetloveinyourheart.kittens.clients.PlayerProfileResponse
	8,  // 11: com.sweetloveinyourheart.kittens.clients.ClientServer.CreateLobby:output_type -> com.sweetloveinyourheart.kittens.clients.CreateLobbyResponse
	10, // 12: com.sweetloveinyourheart.kittens.clients.ClientServer.StreamLobby:output_type -> com.sweetloveinyourheart.kittens.clients.GetLobbyReply
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_clientserver_proto_init() }
func file_clientserver_proto_init() {
	if File_clientserver_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_clientserver_proto_rawDesc), len(file_clientserver_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_clientserver_proto_goTypes,
		DependencyIndexes: file_clientserver_proto_depIdxs,
		MessageInfos:      file_clientserver_proto_msgTypes,
	}.Build()
	File_clientserver_proto = out.File
	file_clientserver_proto_goTypes = nil
	file_clientserver_proto_depIdxs = nil
}
