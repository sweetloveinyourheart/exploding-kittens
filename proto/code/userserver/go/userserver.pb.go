// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: userserver.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateUserRequest_AuthProvider int32

const (
	CreateUserRequest_GUEST  CreateUserRequest_AuthProvider = 0 // Guest user
	CreateUserRequest_GOOGLE CreateUserRequest_AuthProvider = 1 // Google SSO user
)

// Enum value maps for CreateUserRequest_AuthProvider.
var (
	CreateUserRequest_AuthProvider_name = map[int32]string{
		0: "GUEST",
		1: "GOOGLE",
	}
	CreateUserRequest_AuthProvider_value = map[string]int32{
		"GUEST":  0,
		"GOOGLE": 1,
	}
)

func (x CreateUserRequest_AuthProvider) Enum() *CreateUserRequest_AuthProvider {
	p := new(CreateUserRequest_AuthProvider)
	*p = x
	return p
}

func (x CreateUserRequest_AuthProvider) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CreateUserRequest_AuthProvider) Descriptor() protoreflect.EnumDescriptor {
	return file_userserver_proto_enumTypes[0].Descriptor()
}

func (CreateUserRequest_AuthProvider) Type() protoreflect.EnumType {
	return &file_userserver_proto_enumTypes[0]
}

func (x CreateUserRequest_AuthProvider) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CreateUserRequest_AuthProvider.Descriptor instead.
func (CreateUserRequest_AuthProvider) EnumDescriptor() ([]byte, []int) {
	return file_userserver_proto_rawDescGZIP(), []int{3, 0}
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Username string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	FullName string `protobuf:"bytes,3,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Status   int32  `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_userserver_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_userserver_proto_msgTypes[0]
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
	return file_userserver_proto_rawDescGZIP(), []int{0}
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

type GetUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetUserRequest) Reset() {
	*x = GetUserRequest{}
	mi := &file_userserver_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRequest) ProtoMessage() {}

func (x *GetUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_userserver_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserRequest.ProtoReflect.Descriptor instead.
func (*GetUserRequest) Descriptor() ([]byte, []int) {
	return file_userserver_proto_rawDescGZIP(), []int{1}
}

func (x *GetUserRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *GetUserResponse) Reset() {
	*x = GetUserResponse{}
	mi := &file_userserver_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserResponse) ProtoMessage() {}

func (x *GetUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_userserver_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserResponse.ProtoReflect.Descriptor instead.
func (*GetUserResponse) Descriptor() ([]byte, []int) {
	return file_userserver_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type CreateUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username     string                         `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	FullName     string                         `protobuf:"bytes,2,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	AuthProvider CreateUserRequest_AuthProvider `protobuf:"varint,3,opt,name=auth_provider,json=authProvider,proto3,enum=com.sweetloveinyourheart.pocker.users.CreateUserRequest_AuthProvider" json:"auth_provider,omitempty"`
	Meta         *string                        `protobuf:"bytes,4,opt,name=meta,proto3,oneof" json:"meta,omitempty"`
}

func (x *CreateUserRequest) Reset() {
	*x = CreateUserRequest{}
	mi := &file_userserver_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserRequest) ProtoMessage() {}

func (x *CreateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_userserver_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserRequest.ProtoReflect.Descriptor instead.
func (*CreateUserRequest) Descriptor() ([]byte, []int) {
	return file_userserver_proto_rawDescGZIP(), []int{3}
}

func (x *CreateUserRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *CreateUserRequest) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *CreateUserRequest) GetAuthProvider() CreateUserRequest_AuthProvider {
	if x != nil {
		return x.AuthProvider
	}
	return CreateUserRequest_GUEST
}

func (x *CreateUserRequest) GetMeta() string {
	if x != nil && x.Meta != nil {
		return *x.Meta
	}
	return ""
}

type CreateUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *CreateUserResponse) Reset() {
	*x = CreateUserResponse{}
	mi := &file_userserver_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserResponse) ProtoMessage() {}

func (x *CreateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_userserver_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserResponse.ProtoReflect.Descriptor instead.
func (*CreateUserResponse) Descriptor() ([]byte, []int) {
	return file_userserver_proto_rawDescGZIP(), []int{4}
}

func (x *CreateUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type SignInRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *SignInRequest) Reset() {
	*x = SignInRequest{}
	mi := &file_userserver_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignInRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignInRequest) ProtoMessage() {}

func (x *SignInRequest) ProtoReflect() protoreflect.Message {
	mi := &file_userserver_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignInRequest.ProtoReflect.Descriptor instead.
func (*SignInRequest) Descriptor() ([]byte, []int) {
	return file_userserver_proto_rawDescGZIP(), []int{5}
}

func (x *SignInRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type SignInResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // The database id for this user (UUID).
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`                 // The session token for this user.
}

func (x *SignInResponse) Reset() {
	*x = SignInResponse{}
	mi := &file_userserver_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignInResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignInResponse) ProtoMessage() {}

func (x *SignInResponse) ProtoReflect() protoreflect.Message {
	mi := &file_userserver_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignInResponse.ProtoReflect.Descriptor instead.
func (*SignInResponse) Descriptor() ([]byte, []int) {
	return file_userserver_proto_rawDescGZIP(), []int{6}
}

func (x *SignInResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *SignInResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

var File_userserver_proto protoreflect.FileDescriptor

var file_userserver_proto_rawDesc = []byte{
	0x0a, 0x10, 0x75, 0x73, 0x65, 0x72, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x25, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76,
	0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x70, 0x6f, 0x63,
	0x6b, 0x65, 0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x70, 0x0a, 0x04, 0x55, 0x73, 0x65,
	0x72, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x75, 0x6c, 0x6c, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x29, 0x0a, 0x0e, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x52, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77,
	0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61,
	0x72, 0x74, 0x2e, 0x70, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x81, 0x02, 0x0a, 0x11, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x66, 0x75, 0x6c, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x6a, 0x0a, 0x0d, 0x61, 0x75, 0x74,
	0x68, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x45, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65,
	0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x70, 0x6f, 0x63, 0x6b,
	0x65, 0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x50,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x0c, 0x61, 0x75, 0x74, 0x68, 0x50, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x17, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x88, 0x01, 0x01, 0x22, 0x25,
	0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x09,
	0x0a, 0x05, 0x47, 0x55, 0x45, 0x53, 0x54, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x47, 0x4f, 0x4f,
	0x47, 0x4c, 0x45, 0x10, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x22, 0x55,
	0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f,
	0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x70, 0x6f,
	0x63, 0x6b, 0x65, 0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x28, 0x0a, 0x0d, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22,
	0x3f, 0x0a, 0x0e, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x32, 0x84, 0x03, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12,
	0x78, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x35, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72,
	0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x70, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x36, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76,
	0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x70, 0x6f, 0x63,
	0x6b, 0x65, 0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x84, 0x01, 0x0a, 0x0d, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x77, 0x55, 0x73, 0x65, 0x72, 0x12, 0x38, 0x2e, 0x63, 0x6f,
	0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75,
	0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x70, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x39, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65,
	0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74,
	0x2e, 0x70, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x75, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x12, 0x34, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72,
	0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x70, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x35, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65,
	0x69, 0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2e, 0x70, 0x6f, 0x63, 0x6b,
	0x65, 0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x4f, 0x5a, 0x4d, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x77, 0x65, 0x65, 0x74, 0x6c, 0x6f, 0x76, 0x65, 0x69,
	0x6e, 0x79, 0x6f, 0x75, 0x72, 0x68, 0x65, 0x61, 0x72, 0x74, 0x2f, 0x70, 0x6c, 0x61, 0x6e, 0x6e,
	0x69, 0x6e, 0x67, 0x2d, 0x70, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2f, 0x67, 0x6f, 0x3b, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_userserver_proto_rawDescOnce sync.Once
	file_userserver_proto_rawDescData = file_userserver_proto_rawDesc
)

func file_userserver_proto_rawDescGZIP() []byte {
	file_userserver_proto_rawDescOnce.Do(func() {
		file_userserver_proto_rawDescData = protoimpl.X.CompressGZIP(file_userserver_proto_rawDescData)
	})
	return file_userserver_proto_rawDescData
}

var file_userserver_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_userserver_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_userserver_proto_goTypes = []any{
	(CreateUserRequest_AuthProvider)(0), // 0: com.sweetloveinyourheart.pocker.users.CreateUserRequest.AuthProvider
	(*User)(nil),                        // 1: com.sweetloveinyourheart.pocker.users.User
	(*GetUserRequest)(nil),              // 2: com.sweetloveinyourheart.pocker.users.GetUserRequest
	(*GetUserResponse)(nil),             // 3: com.sweetloveinyourheart.pocker.users.GetUserResponse
	(*CreateUserRequest)(nil),           // 4: com.sweetloveinyourheart.pocker.users.CreateUserRequest
	(*CreateUserResponse)(nil),          // 5: com.sweetloveinyourheart.pocker.users.CreateUserResponse
	(*SignInRequest)(nil),               // 6: com.sweetloveinyourheart.pocker.users.SignInRequest
	(*SignInResponse)(nil),              // 7: com.sweetloveinyourheart.pocker.users.SignInResponse
}
var file_userserver_proto_depIdxs = []int32{
	1, // 0: com.sweetloveinyourheart.pocker.users.GetUserResponse.user:type_name -> com.sweetloveinyourheart.pocker.users.User
	0, // 1: com.sweetloveinyourheart.pocker.users.CreateUserRequest.auth_provider:type_name -> com.sweetloveinyourheart.pocker.users.CreateUserRequest.AuthProvider
	1, // 2: com.sweetloveinyourheart.pocker.users.CreateUserResponse.user:type_name -> com.sweetloveinyourheart.pocker.users.User
	2, // 3: com.sweetloveinyourheart.pocker.users.UserServer.GetUser:input_type -> com.sweetloveinyourheart.pocker.users.GetUserRequest
	4, // 4: com.sweetloveinyourheart.pocker.users.UserServer.CreateNewUser:input_type -> com.sweetloveinyourheart.pocker.users.CreateUserRequest
	6, // 5: com.sweetloveinyourheart.pocker.users.UserServer.SignIn:input_type -> com.sweetloveinyourheart.pocker.users.SignInRequest
	3, // 6: com.sweetloveinyourheart.pocker.users.UserServer.GetUser:output_type -> com.sweetloveinyourheart.pocker.users.GetUserResponse
	5, // 7: com.sweetloveinyourheart.pocker.users.UserServer.CreateNewUser:output_type -> com.sweetloveinyourheart.pocker.users.CreateUserResponse
	7, // 8: com.sweetloveinyourheart.pocker.users.UserServer.SignIn:output_type -> com.sweetloveinyourheart.pocker.users.SignInResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_userserver_proto_init() }
func file_userserver_proto_init() {
	if File_userserver_proto != nil {
		return
	}
	file_userserver_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_userserver_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_userserver_proto_goTypes,
		DependencyIndexes: file_userserver_proto_depIdxs,
		EnumInfos:         file_userserver_proto_enumTypes,
		MessageInfos:      file_userserver_proto_msgTypes,
	}.Build()
	File_userserver_proto = out.File
	file_userserver_proto_rawDesc = nil
	file_userserver_proto_goTypes = nil
	file_userserver_proto_depIdxs = nil
}
