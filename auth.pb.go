// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: khepri/horus/auth.proto

package horus

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

type BasicSignUpRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *BasicSignUpRequest) Reset() {
	*x = BasicSignUpRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicSignUpRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicSignUpRequest) ProtoMessage() {}

func (x *BasicSignUpRequest) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicSignUpRequest.ProtoReflect.Descriptor instead.
func (*BasicSignUpRequest) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{0}
}

func (x *BasicSignUpRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *BasicSignUpRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type BasicSignUpResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token *Token `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *BasicSignUpResponse) Reset() {
	*x = BasicSignUpResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicSignUpResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicSignUpResponse) ProtoMessage() {}

func (x *BasicSignUpResponse) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicSignUpResponse.ProtoReflect.Descriptor instead.
func (*BasicSignUpResponse) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{1}
}

func (x *BasicSignUpResponse) GetToken() *Token {
	if x != nil {
		return x.Token
	}
	return nil
}

type BasicSignInRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *BasicSignInRequest) Reset() {
	*x = BasicSignInRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicSignInRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicSignInRequest) ProtoMessage() {}

func (x *BasicSignInRequest) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicSignInRequest.ProtoReflect.Descriptor instead.
func (*BasicSignInRequest) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{2}
}

func (x *BasicSignInRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *BasicSignInRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type BasicSignInResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token *Token `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"` // Access token.
}

func (x *BasicSignInResponse) Reset() {
	*x = BasicSignInResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicSignInResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicSignInResponse) ProtoMessage() {}

func (x *BasicSignInResponse) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicSignInResponse.ProtoReflect.Descriptor instead.
func (*BasicSignInResponse) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{3}
}

func (x *BasicSignInResponse) GetToken() *Token {
	if x != nil {
		return x.Token
	}
	return nil
}

type TokenSignInRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"` // Access token.
}

func (x *TokenSignInRequest) Reset() {
	*x = TokenSignInRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenSignInRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenSignInRequest) ProtoMessage() {}

func (x *TokenSignInRequest) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenSignInRequest.ProtoReflect.Descriptor instead.
func (*TokenSignInRequest) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{4}
}

func (x *TokenSignInRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type TokenSignInResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token *Token `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"` // Access token.
}

func (x *TokenSignInResponse) Reset() {
	*x = TokenSignInResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenSignInResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenSignInResponse) ProtoMessage() {}

func (x *TokenSignInResponse) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenSignInResponse.ProtoReflect.Descriptor instead.
func (*TokenSignInResponse) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{5}
}

func (x *TokenSignInResponse) GetToken() *Token {
	if x != nil {
		return x.Token
	}
	return nil
}

type RefreshRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"` // Refresh token.
}

func (x *RefreshRequest) Reset() {
	*x = RefreshRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshRequest) ProtoMessage() {}

func (x *RefreshRequest) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshRequest.ProtoReflect.Descriptor instead.
func (*RefreshRequest) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{6}
}

func (x *RefreshRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type RefreshResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token *Token `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"` // Access token.
}

func (x *RefreshResponse) Reset() {
	*x = RefreshResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshResponse) ProtoMessage() {}

func (x *RefreshResponse) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshResponse.ProtoReflect.Descriptor instead.
func (*RefreshResponse) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{7}
}

func (x *RefreshResponse) GetToken() *Token {
	if x != nil {
		return x.Token
	}
	return nil
}

type VerifyOtpRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *VerifyOtpRequest) Reset() {
	*x = VerifyOtpRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyOtpRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyOtpRequest) ProtoMessage() {}

func (x *VerifyOtpRequest) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyOtpRequest.ProtoReflect.Descriptor instead.
func (*VerifyOtpRequest) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{8}
}

func (x *VerifyOtpRequest) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type VerifyOtpResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token *Token `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *VerifyOtpResponse) Reset() {
	*x = VerifyOtpResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyOtpResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyOtpResponse) ProtoMessage() {}

func (x *VerifyOtpResponse) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyOtpResponse.ProtoReflect.Descriptor instead.
func (*VerifyOtpResponse) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{9}
}

func (x *VerifyOtpResponse) GetToken() *Token {
	if x != nil {
		return x.Token
	}
	return nil
}

type SingOutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"` // Access token.
}

func (x *SingOutRequest) Reset() {
	*x = SingOutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SingOutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SingOutRequest) ProtoMessage() {}

func (x *SingOutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SingOutRequest.ProtoReflect.Descriptor instead.
func (*SingOutRequest) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{10}
}

func (x *SingOutRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type SingOutResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SingOutResponse) Reset() {
	*x = SingOutResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_khepri_horus_auth_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SingOutResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SingOutResponse) ProtoMessage() {}

func (x *SingOutResponse) ProtoReflect() protoreflect.Message {
	mi := &file_khepri_horus_auth_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SingOutResponse.ProtoReflect.Descriptor instead.
func (*SingOutResponse) Descriptor() ([]byte, []int) {
	return file_khepri_horus_auth_proto_rawDescGZIP(), []int{11}
}

var File_khepri_horus_auth_proto protoreflect.FileDescriptor

var file_khepri_horus_auth_proto_rawDesc = []byte{
	0x0a, 0x17, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2f, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x2f, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6b, 0x68, 0x65, 0x70, 0x72,
	0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x1a, 0x19, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2f,
	0x68, 0x6f, 0x72, 0x75, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x4c, 0x0a, 0x12, 0x42, 0x61, 0x73, 0x69, 0x63, 0x53, 0x69, 0x67, 0x6e, 0x55,
	0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x22, 0x40, 0x0a, 0x13, 0x42, 0x61, 0x73, 0x69, 0x63, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e,
	0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x05, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x22, 0x4c, 0x0a, 0x12, 0x42, 0x61, 0x73, 0x69, 0x63, 0x53, 0x69, 0x67, 0x6e, 0x49,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x22, 0x40, 0x0a, 0x13, 0x42, 0x61, 0x73, 0x69, 0x63, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e,
	0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x05, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x22, 0x2a, 0x0a, 0x12, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x53, 0x69, 0x67, 0x6e, 0x49,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x40,
	0x0a, 0x13, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f,
	0x72, 0x75, 0x73, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x22, 0x26, 0x0a, 0x0e, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x3c, 0x0a, 0x0f, 0x52, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6b, 0x68, 0x65,
	0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x28, 0x0a, 0x10, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x4f, 0x74, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x3e, 0x0a, 0x11, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x4f, 0x74, 0x70, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f,
	0x72, 0x75, 0x73, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x22, 0x26, 0x0a, 0x0e, 0x53, 0x69, 0x6e, 0x67, 0x4f, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x11, 0x0a, 0x0f, 0x53, 0x69, 0x6e, 0x67,
	0x4f, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xe7, 0x03, 0x0a, 0x0b,
	0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x52, 0x0a, 0x0b, 0x42,
	0x61, 0x73, 0x69, 0x63, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70, 0x12, 0x20, 0x2e, 0x6b, 0x68, 0x65,
	0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x53,
	0x69, 0x67, 0x6e, 0x55, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x6b,
	0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e, 0x42, 0x61, 0x73, 0x69,
	0x63, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x52, 0x0a, 0x0b, 0x42, 0x61, 0x73, 0x69, 0x63, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x12, 0x20,
	0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e, 0x42, 0x61,
	0x73, 0x69, 0x63, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x21, 0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e,
	0x42, 0x61, 0x73, 0x69, 0x63, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x52, 0x0a, 0x0b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x53, 0x69, 0x67, 0x6e,
	0x49, 0x6e, 0x12, 0x20, 0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75,
	0x73, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f,
	0x72, 0x75, 0x73, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x07, 0x52, 0x65, 0x66, 0x72, 0x65,
	0x73, 0x68, 0x12, 0x1c, 0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75,
	0x73, 0x2e, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e,
	0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x4c, 0x0a, 0x09, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x4f, 0x74, 0x70, 0x12, 0x1e, 0x2e, 0x6b,
	0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e, 0x56, 0x65, 0x72, 0x69,
	0x66, 0x79, 0x4f, 0x74, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x6b,
	0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e, 0x56, 0x65, 0x72, 0x69,
	0x66, 0x79, 0x4f, 0x74, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a,
	0x07, 0x53, 0x69, 0x67, 0x6e, 0x4f, 0x75, 0x74, 0x12, 0x1c, 0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72,
	0x69, 0x2e, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x4f, 0x75, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e,
	0x68, 0x6f, 0x72, 0x75, 0x73, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x4f, 0x75, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x12, 0x5a, 0x10, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e,
	0x64, 0x65, 0x76, 0x2f, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_khepri_horus_auth_proto_rawDescOnce sync.Once
	file_khepri_horus_auth_proto_rawDescData = file_khepri_horus_auth_proto_rawDesc
)

func file_khepri_horus_auth_proto_rawDescGZIP() []byte {
	file_khepri_horus_auth_proto_rawDescOnce.Do(func() {
		file_khepri_horus_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_khepri_horus_auth_proto_rawDescData)
	})
	return file_khepri_horus_auth_proto_rawDescData
}

var file_khepri_horus_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_khepri_horus_auth_proto_goTypes = []any{
	(*BasicSignUpRequest)(nil),  // 0: khepri.horus.BasicSignUpRequest
	(*BasicSignUpResponse)(nil), // 1: khepri.horus.BasicSignUpResponse
	(*BasicSignInRequest)(nil),  // 2: khepri.horus.BasicSignInRequest
	(*BasicSignInResponse)(nil), // 3: khepri.horus.BasicSignInResponse
	(*TokenSignInRequest)(nil),  // 4: khepri.horus.TokenSignInRequest
	(*TokenSignInResponse)(nil), // 5: khepri.horus.TokenSignInResponse
	(*RefreshRequest)(nil),      // 6: khepri.horus.RefreshRequest
	(*RefreshResponse)(nil),     // 7: khepri.horus.RefreshResponse
	(*VerifyOtpRequest)(nil),    // 8: khepri.horus.VerifyOtpRequest
	(*VerifyOtpResponse)(nil),   // 9: khepri.horus.VerifyOtpResponse
	(*SingOutRequest)(nil),      // 10: khepri.horus.SingOutRequest
	(*SingOutResponse)(nil),     // 11: khepri.horus.SingOutResponse
	(*Token)(nil),               // 12: khepri.horus.Token
}
var file_khepri_horus_auth_proto_depIdxs = []int32{
	12, // 0: khepri.horus.BasicSignUpResponse.token:type_name -> khepri.horus.Token
	12, // 1: khepri.horus.BasicSignInResponse.token:type_name -> khepri.horus.Token
	12, // 2: khepri.horus.TokenSignInResponse.token:type_name -> khepri.horus.Token
	12, // 3: khepri.horus.RefreshResponse.token:type_name -> khepri.horus.Token
	12, // 4: khepri.horus.VerifyOtpResponse.token:type_name -> khepri.horus.Token
	0,  // 5: khepri.horus.AuthService.BasicSignUp:input_type -> khepri.horus.BasicSignUpRequest
	2,  // 6: khepri.horus.AuthService.BasicSignIn:input_type -> khepri.horus.BasicSignInRequest
	4,  // 7: khepri.horus.AuthService.TokenSignIn:input_type -> khepri.horus.TokenSignInRequest
	6,  // 8: khepri.horus.AuthService.Refresh:input_type -> khepri.horus.RefreshRequest
	8,  // 9: khepri.horus.AuthService.VerifyOtp:input_type -> khepri.horus.VerifyOtpRequest
	10, // 10: khepri.horus.AuthService.SignOut:input_type -> khepri.horus.SingOutRequest
	1,  // 11: khepri.horus.AuthService.BasicSignUp:output_type -> khepri.horus.BasicSignUpResponse
	3,  // 12: khepri.horus.AuthService.BasicSignIn:output_type -> khepri.horus.BasicSignInResponse
	5,  // 13: khepri.horus.AuthService.TokenSignIn:output_type -> khepri.horus.TokenSignInResponse
	7,  // 14: khepri.horus.AuthService.Refresh:output_type -> khepri.horus.RefreshResponse
	9,  // 15: khepri.horus.AuthService.VerifyOtp:output_type -> khepri.horus.VerifyOtpResponse
	11, // 16: khepri.horus.AuthService.SignOut:output_type -> khepri.horus.SingOutResponse
	11, // [11:17] is the sub-list for method output_type
	5,  // [5:11] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_khepri_horus_auth_proto_init() }
func file_khepri_horus_auth_proto_init() {
	if File_khepri_horus_auth_proto != nil {
		return
	}
	file_khepri_horus_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_khepri_horus_auth_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*BasicSignUpRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_khepri_horus_auth_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*BasicSignUpResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_khepri_horus_auth_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*BasicSignInRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_khepri_horus_auth_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*BasicSignInResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_khepri_horus_auth_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*TokenSignInRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_khepri_horus_auth_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*TokenSignInResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_khepri_horus_auth_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*RefreshRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_khepri_horus_auth_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*RefreshResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_khepri_horus_auth_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*VerifyOtpRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_khepri_horus_auth_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*VerifyOtpResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_khepri_horus_auth_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*SingOutRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_khepri_horus_auth_proto_msgTypes[11].Exporter = func(v any, i int) any {
			switch v := v.(*SingOutResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_khepri_horus_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_khepri_horus_auth_proto_goTypes,
		DependencyIndexes: file_khepri_horus_auth_proto_depIdxs,
		MessageInfos:      file_khepri_horus_auth_proto_msgTypes,
	}.Build()
	File_khepri_horus_auth_proto = out.File
	file_khepri_horus_auth_proto_rawDesc = nil
	file_khepri_horus_auth_proto_goTypes = nil
	file_khepri_horus_auth_proto_depIdxs = nil
}
