// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: hi.proto

package hi

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type HiReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	In string `protobuf:"bytes,1,opt,name=in,proto3" json:"in,omitempty"`
}

func (x *HiReq) Reset() {
	*x = HiReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HiReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HiReq) ProtoMessage() {
  //func (*HiReq) ProtoMessage() 
}

func (x *HiReq) ProtoReflect() protoreflect.Message {
	mi := &file_hi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HiReq.ProtoReflect.Descriptor instead.
func (*HiReq) Descriptor() ([]byte, []int) {
	return filehiprotorawDescGZIP(), []int{0}
}

func (x *HiReq) GetIn() string {
	if x != nil {
		return x.In
	}
	return ""
}

type HelloReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	In string `protobuf:"bytes,1,opt,name=in,proto3" json:"in,omitempty"`
}

func (x *HelloReq) Reset() {
	*x = HelloReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloReq) ProtoMessage() {
  //func (*HelloReq) ProtoMessage() 
}

func (x *HelloReq) ProtoReflect() protoreflect.Message {
	mi := &file_hi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloReq.ProtoReflect.Descriptor instead.
func (*HelloReq) Descriptor() ([]byte, []int) {
	return filehiprotorawDescGZIP(), []int{1}
}

func (x *HelloReq) GetIn() string {
	if x != nil {
		return x.In
	}
	return ""
}

type HiResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *HiResp) Reset() {
	*x = HiResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hi_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HiResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HiResp) ProtoMessage() {
  //func (*HiResp) ProtoMessage() 
}

func (x *HiResp) ProtoReflect() protoreflect.Message {
	mi := &file_hi_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HiResp.ProtoReflect.Descriptor instead.
func (*HiResp) Descriptor() ([]byte, []int) {
	return filehiprotorawDescGZIP(), []int{2}
}

func (x *HiResp) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type HelloResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *HelloResp) Reset() {
	*x = HelloResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hi_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloResp) ProtoMessage() {
  //func (*HelloResp) ProtoMessage() 
}

func (x *HelloResp) ProtoReflect() protoreflect.Message {
	mi := &file_hi_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloResp.ProtoReflect.Descriptor instead.
func (*HelloResp) Descriptor() ([]byte, []int) {
	return filehiprotorawDescGZIP(), []int{3}
}

func (x *HelloResp) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type EventReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EventReq) Reset() {
	*x = EventReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hi_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventReq) ProtoMessage() {
  //func (*EventReq) ProtoMessage() 
}

func (x *EventReq) ProtoReflect() protoreflect.Message {
	mi := &file_hi_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventReq.ProtoReflect.Descriptor instead.
func (*EventReq) Descriptor() ([]byte, []int) {
	return filehiprotorawDescGZIP(), []int{4}
}

type EventResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EventResp) Reset() {
	*x = EventResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hi_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventResp) ProtoMessage() {
  //func (*EventResp) ProtoMessage() 
}

func (x *EventResp) ProtoReflect() protoreflect.Message {
	mi := &file_hi_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventResp.ProtoReflect.Descriptor instead.
func (*EventResp) Descriptor() ([]byte, []int) {
	return filehiprotorawDescGZIP(), []int{5}
}

var File_hi_proto protoreflect.FileDescriptor

var file_hi_proto_rawDesc = []byte{
	0x0a, 0x08, 0x68, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x68, 0x69, 0x22, 0x17,
	0x0a, 0x05, 0x48, 0x69, 0x52, 0x65, 0x71, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x6e, 0x22, 0x1a, 0x0a, 0x08, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
	0x52, 0x65, 0x71, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x6e, 0x22, 0x1a, 0x0a, 0x06, 0x48, 0x69, 0x52, 0x65, 0x73, 0x70, 0x12, 0x10, 0x0a,
	0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22,
	0x1d, 0x0a, 0x09, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x12, 0x10, 0x0a, 0x03,
	0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x0a,
	0x0a, 0x08, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x22, 0x0b, 0x0a, 0x09, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x32, 0x50, 0x0a, 0x05, 0x47, 0x72, 0x65, 0x65, 0x74,
	0x12, 0x1e, 0x0a, 0x05, 0x53, 0x61, 0x79, 0x48, 0x69, 0x12, 0x09, 0x2e, 0x68, 0x69, 0x2e, 0x48,
	0x69, 0x52, 0x65, 0x71, 0x1a, 0x0a, 0x2e, 0x68, 0x69, 0x2e, 0x48, 0x69, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x27, 0x0a, 0x08, 0x53, 0x61, 0x79, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x0c, 0x2e, 0x68,
	0x69, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71, 0x1a, 0x0d, 0x2e, 0x68, 0x69, 0x2e,
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x32, 0x33, 0x0a, 0x05, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x2a, 0x0a, 0x0b, 0x41, 0x73, 0x6b, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x0c, 0x2e, 0x68, 0x69, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x1a,
	0x0d, 0x2e, 0x68, 0x69, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x42, 0x06,
	0x5a, 0x04, 0x2e, 0x2f, 0x68, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_hi_proto_rawDescOnce sync.Once
	file_hi_proto_rawDescData = file_hi_proto_rawDesc
)

func filehiprotorawDescGZIP() []byte {
	file_hi_proto_rawDescOnce.Do(func() {
		file_hi_proto_rawDescData = protoimpl.X.CompressGZIP(file_hi_proto_rawDescData)
	})
	return file_hi_proto_rawDescData
}

var file_hi_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_hi_proto_goTypes = []any{
	(*HiReq)(nil),     // 0: hi.HiReq
	(*HelloReq)(nil),  // 1: hi.HelloReq
	(*HiResp)(nil),    // 2: hi.HiResp
	(*HelloResp)(nil), // 3: hi.HelloResp
	(*EventReq)(nil),  // 4: hi.EventReq
	(*EventResp)(nil), // 5: hi.EventResp
}
var file_hi_proto_depIdxs = []int32{
	0, // 0: hi.Greet.SayHi:input_type -> hi.HiReq
	1, // 1: hi.Greet.SayHello:input_type -> hi.HelloReq
	4, // 2: hi.Event.AskQuestion:input_type -> hi.EventReq
	2, // 3: hi.Greet.SayHi:output_type -> hi.HiResp
	3, // 4: hi.Greet.SayHello:output_type -> hi.HelloResp
	5, // 5: hi.Event.AskQuestion:output_type -> hi.EventResp
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { filehiprotoinit() }
func filehiprotoinit() {//NOSONAR
	if File_hi_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_hi_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*HiReq); i {
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
		file_hi_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*HelloReq); i {
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
		file_hi_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*HiResp); i {
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
		file_hi_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*HelloResp); i {
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
		file_hi_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*EventReq); i {
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
		file_hi_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*EventResp); i {
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
			RawDescriptor: file_hi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_hi_proto_goTypes,
		DependencyIndexes: file_hi_proto_depIdxs,
		MessageInfos:      file_hi_proto_msgTypes,
	}.Build()
	File_hi_proto = out.File
	file_hi_proto_rawDesc = nil
	file_hi_proto_goTypes = nil
	file_hi_proto_depIdxs = nil
}
