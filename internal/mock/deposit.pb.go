// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.8.0
// source: deposit.proto

package mock

import (
	context "context"
	reflect "reflect"
	sync "sync"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DepositRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount float32 `protobuf:"fixed32,1,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *DepositRequest) Reset() {
	*x = DepositRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deposit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DepositRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DepositRequest) ProtoMessage() {
  //func (*DepositRequest) ProtoMessage() 
}

func (x *DepositRequest) ProtoReflect() protoreflect.Message {
	mi := &file_deposit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DepositRequest.ProtoReflect.Descriptor instead.
func (*DepositRequest) Descriptor() ([]byte, []int) {
	return fileDepositProtoRawDescGZIP(), []int{0}
}

func (x *DepositRequest) GetAmount() float32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type DepositResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *DepositResponse) Reset() {
	*x = DepositResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deposit_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DepositResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DepositResponse) ProtoMessage() {
  //func (*DepositResponse) ProtoMessage() 
}

func (x *DepositResponse) ProtoReflect() protoreflect.Message {
	mi := &file_deposit_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DepositResponse.ProtoReflect.Descriptor instead.
func (*DepositResponse) Descriptor() ([]byte, []int) {
	return fileDepositProtoRawDescGZIP(), []int{1}
}

func (x *DepositResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

var File_deposit_proto protoreflect.FileDescriptor

var file_deposit_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x04, 0x6d, 0x6f, 0x63, 0x6b, 0x22, 0x28, 0x0a, 0x0e, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x21, 0x0a, 0x0f, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02,
	0x6f, 0x6b, 0x32, 0x48, 0x0a, 0x0e, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x12,
	0x14, 0x2e, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6d, 0x6f, 0x63, 0x6b, 0x2e, 0x44, 0x65, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x08, 0x5a, 0x06,
	0x2e, 0x3b, 0x6d, 0x6f, 0x63, 0x6b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_deposit_proto_rawDescOnce sync.Once
	file_deposit_proto_rawDescData = file_deposit_proto_rawDesc
)

func fileDepositProtoRawDescGZIP() []byte {
	file_deposit_proto_rawDescOnce.Do(func() {
		file_deposit_proto_rawDescData = protoimpl.X.CompressGZIP(file_deposit_proto_rawDescData)
	})
	return file_deposit_proto_rawDescData
}

var (
	file_deposit_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
	file_deposit_proto_goTypes  = []any{
		(*DepositRequest)(nil),  // 0: mock.DepositRequest
		(*DepositResponse)(nil), // 1: mock.DepositResponse
	}
)

var file_deposit_proto_depIdxs = []int32{
	0, // 0: mock.DepositService.Deposit:input_type -> mock.DepositRequest
	1, // 1: mock.DepositService.Deposit:output_type -> mock.DepositResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { filedepositprotoinit() }
func filedepositprotoinit() {
	if File_deposit_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_deposit_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*DepositRequest); i {
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
		file_deposit_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*DepositResponse); i {
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
			RawDescriptor: file_deposit_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_deposit_proto_goTypes,
		DependencyIndexes: file_deposit_proto_depIdxs,
		MessageInfos:      file_deposit_proto_msgTypes,
	}.Build()
	File_deposit_proto = out.File
	file_deposit_proto_rawDesc = nil
	file_deposit_proto_goTypes = nil
	file_deposit_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ context.Context
	_ grpc.ClientConnInterface
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DepositServiceClient is the client API for DepositService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DepositServiceClient interface {
	Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositResponse, error)
}

type depositServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDepositServiceClient(cc grpc.ClientConnInterface) DepositServiceClient {
	return &depositServiceClient{cc}
}

func (c *depositServiceClient) Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositResponse, error) {
	out := new(DepositResponse)
	err := c.cc.Invoke(ctx, "/mock.DepositService/Deposit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DepositServiceServer is the server API for DepositService service.
type DepositServiceServer interface {
	Deposit(context.Context, *DepositRequest) (*DepositResponse, error)
}

// UnimplementedDepositServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDepositServiceServer struct{}

func (*UnimplementedDepositServiceServer) Deposit(context.Context, *DepositRequest) (*DepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposit not implemented")
}

func RegisterDepositServiceServer(s *grpc.Server, srv DepositServiceServer) {
	s.RegisterService(&_DepositService_serviceDesc, srv)
}

func DepositServiceDepositHandler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(DepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServiceServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mock.DepositService/Deposit",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(DepositServiceServer).Deposit(ctx, req.(*DepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DepositService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mock.DepositService",
	HandlerType: (*DepositServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deposit",
			Handler:    DepositServiceDepositHandler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "deposit.proto",
}
