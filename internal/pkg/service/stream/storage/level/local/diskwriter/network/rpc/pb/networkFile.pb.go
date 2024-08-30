// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: networkFile.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OpenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SourceNodeId  string `protobuf:"bytes,1,opt,name=source_node_id,json=sourceNodeId,proto3" json:"source_node_id,omitempty"`
	SliceDataJson []byte `protobuf:"bytes,2,opt,name=slice_data_json,json=sliceDataJson,proto3" json:"slice_data_json,omitempty"`
}

func (x *OpenRequest) Reset() {
	*x = OpenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networkFile_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpenRequest) ProtoMessage() {}

func (x *OpenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_networkFile_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpenRequest.ProtoReflect.Descriptor instead.
func (*OpenRequest) Descriptor() ([]byte, []int) {
	return file_networkFile_proto_rawDescGZIP(), []int{0}
}

func (x *OpenRequest) GetSourceNodeId() string {
	if x != nil {
		return x.SourceNodeId
	}
	return ""
}

func (x *OpenRequest) GetSliceDataJson() []byte {
	if x != nil {
		return x.SliceDataJson
	}
	return nil
}

type OpenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId uint64 `protobuf:"varint,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
}

func (x *OpenResponse) Reset() {
	*x = OpenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networkFile_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpenResponse) ProtoMessage() {}

func (x *OpenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_networkFile_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpenResponse.ProtoReflect.Descriptor instead.
func (*OpenResponse) Descriptor() ([]byte, []int) {
	return file_networkFile_proto_rawDescGZIP(), []int{1}
}

func (x *OpenResponse) GetFileId() uint64 {
	if x != nil {
		return x.FileId
	}
	return 0
}

type WaitForServerTerminationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *WaitForServerTerminationRequest) Reset() {
	*x = WaitForServerTerminationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networkFile_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WaitForServerTerminationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WaitForServerTerminationRequest) ProtoMessage() {}

func (x *WaitForServerTerminationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_networkFile_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WaitForServerTerminationRequest.ProtoReflect.Descriptor instead.
func (*WaitForServerTerminationRequest) Descriptor() ([]byte, []int) {
	return file_networkFile_proto_rawDescGZIP(), []int{2}
}

type ServerIsTerminatingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ServerIsTerminatingResponse) Reset() {
	*x = ServerIsTerminatingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networkFile_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerIsTerminatingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerIsTerminatingResponse) ProtoMessage() {}

func (x *ServerIsTerminatingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_networkFile_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerIsTerminatingResponse.ProtoReflect.Descriptor instead.
func (*ServerIsTerminatingResponse) Descriptor() ([]byte, []int) {
	return file_networkFile_proto_rawDescGZIP(), []int{3}
}

type WriteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId uint64 `protobuf:"varint,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	// The "aligned" is true after the pipeline Flush or Close operation.
	// It means that a complete block of data has been written so far.
	Aligned bool   `protobuf:"varint,2,opt,name=aligned,proto3" json:"aligned,omitempty"`
	Data    []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *WriteRequest) Reset() {
	*x = WriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networkFile_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteRequest) ProtoMessage() {}

func (x *WriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_networkFile_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteRequest.ProtoReflect.Descriptor instead.
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return file_networkFile_proto_rawDescGZIP(), []int{4}
}

func (x *WriteRequest) GetFileId() uint64 {
	if x != nil {
		return x.FileId
	}
	return 0
}

func (x *WriteRequest) GetAligned() bool {
	if x != nil {
		return x.Aligned
	}
	return false
}

func (x *WriteRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type WriteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	N int64 `protobuf:"varint,2,opt,name=n,proto3" json:"n,omitempty"`
}

func (x *WriteResponse) Reset() {
	*x = WriteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networkFile_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteResponse) ProtoMessage() {}

func (x *WriteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_networkFile_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteResponse.ProtoReflect.Descriptor instead.
func (*WriteResponse) Descriptor() ([]byte, []int) {
	return file_networkFile_proto_rawDescGZIP(), []int{5}
}

func (x *WriteResponse) GetN() int64 {
	if x != nil {
		return x.N
	}
	return 0
}

type SyncRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId uint64 `protobuf:"varint,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
}

func (x *SyncRequest) Reset() {
	*x = SyncRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networkFile_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncRequest) ProtoMessage() {}

func (x *SyncRequest) ProtoReflect() protoreflect.Message {
	mi := &file_networkFile_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncRequest.ProtoReflect.Descriptor instead.
func (*SyncRequest) Descriptor() ([]byte, []int) {
	return file_networkFile_proto_rawDescGZIP(), []int{6}
}

func (x *SyncRequest) GetFileId() uint64 {
	if x != nil {
		return x.FileId
	}
	return 0
}

type SyncResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SyncResponse) Reset() {
	*x = SyncResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networkFile_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncResponse) ProtoMessage() {}

func (x *SyncResponse) ProtoReflect() protoreflect.Message {
	mi := &file_networkFile_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncResponse.ProtoReflect.Descriptor instead.
func (*SyncResponse) Descriptor() ([]byte, []int) {
	return file_networkFile_proto_rawDescGZIP(), []int{7}
}

type CloseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId uint64 `protobuf:"varint,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
}

func (x *CloseRequest) Reset() {
	*x = CloseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networkFile_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloseRequest) ProtoMessage() {}

func (x *CloseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_networkFile_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloseRequest.ProtoReflect.Descriptor instead.
func (*CloseRequest) Descriptor() ([]byte, []int) {
	return file_networkFile_proto_rawDescGZIP(), []int{8}
}

func (x *CloseRequest) GetFileId() uint64 {
	if x != nil {
		return x.FileId
	}
	return 0
}

type CloseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CloseResponse) Reset() {
	*x = CloseResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networkFile_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloseResponse) ProtoMessage() {}

func (x *CloseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_networkFile_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloseResponse.ProtoReflect.Descriptor instead.
func (*CloseResponse) Descriptor() ([]byte, []int) {
	return file_networkFile_proto_rawDescGZIP(), []int{9}
}

type SliceKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId int64                  `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	BranchId  int64                  `protobuf:"varint,2,opt,name=branch_id,json=branchId,proto3" json:"branch_id,omitempty"`
	SourceId  string                 `protobuf:"bytes,3,opt,name=source_id,json=sourceId,proto3" json:"source_id,omitempty"`
	SinkId    string                 `protobuf:"bytes,4,opt,name=sink_id,json=sinkId,proto3" json:"sink_id,omitempty"`
	FileId    *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	VolumeId  string                 `protobuf:"bytes,6,opt,name=volume_id,json=volumeId,proto3" json:"volume_id,omitempty"`
	SliceId   *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=slice_id,json=sliceId,proto3" json:"slice_id,omitempty"`
}

func (x *SliceKey) Reset() {
	*x = SliceKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networkFile_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SliceKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SliceKey) ProtoMessage() {}

func (x *SliceKey) ProtoReflect() protoreflect.Message {
	mi := &file_networkFile_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SliceKey.ProtoReflect.Descriptor instead.
func (*SliceKey) Descriptor() ([]byte, []int) {
	return file_networkFile_proto_rawDescGZIP(), []int{10}
}

func (x *SliceKey) GetProjectId() int64 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *SliceKey) GetBranchId() int64 {
	if x != nil {
		return x.BranchId
	}
	return 0
}

func (x *SliceKey) GetSourceId() string {
	if x != nil {
		return x.SourceId
	}
	return ""
}

func (x *SliceKey) GetSinkId() string {
	if x != nil {
		return x.SinkId
	}
	return ""
}

func (x *SliceKey) GetFileId() *timestamppb.Timestamp {
	if x != nil {
		return x.FileId
	}
	return nil
}

func (x *SliceKey) GetVolumeId() string {
	if x != nil {
		return x.VolumeId
	}
	return ""
}

func (x *SliceKey) GetSliceId() *timestamppb.Timestamp {
	if x != nil {
		return x.SliceId
	}
	return nil
}

var File_networkFile_proto protoreflect.FileDescriptor

var file_networkFile_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x46, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5b, 0x0a, 0x0b, 0x4f, 0x70, 0x65, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x26, 0x0a,
	0x0f, 0x73, 0x6c, 0x69, 0x63, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x6a, 0x73, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x73, 0x6c, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x4a, 0x73, 0x6f, 0x6e, 0x22, 0x27, 0x0a, 0x0c, 0x4f, 0x70, 0x65, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x21,
	0x0a, 0x1f, 0x57, 0x61, 0x69, 0x74, 0x46, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x54,
	0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x1d, 0x0a, 0x1b, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x73, 0x54, 0x65, 0x72,
	0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x55, 0x0a, 0x0c, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x6c, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x61, 0x6c, 0x69, 0x67,
	0x6e, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x1d, 0x0a, 0x0d, 0x57, 0x72, 0x69, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x01, 0x6e, 0x22, 0x26, 0x0a, 0x0b, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x0e,
	0x0a, 0x0c, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x27,
	0x0a, 0x0c, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x0f, 0x0a, 0x0d, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x85, 0x02, 0x0a, 0x08, 0x53, 0x6c, 0x69,
	0x63, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x49,
	0x64, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x17,
	0x0a, 0x07, 0x73, 0x69, 0x6e, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x69, 0x6e, 0x6b, 0x49, 0x64, 0x12, 0x33, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09,
	0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x49, 0x64, 0x12, 0x35, 0x0a, 0x08, 0x73, 0x6c, 0x69,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x73, 0x6c, 0x69, 0x63, 0x65, 0x49, 0x64,
	0x32, 0xab, 0x02, 0x0a, 0x0b, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x46, 0x69, 0x6c, 0x65,
	0x12, 0x2b, 0x0a, 0x04, 0x4f, 0x70, 0x65, 0x6e, 0x12, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x4f, 0x70,
	0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x4f,
	0x70, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x62, 0x0a,
	0x18, 0x57, 0x61, 0x69, 0x74, 0x46, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x54, 0x65,
	0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x2e, 0x70, 0x62, 0x2e, 0x57,
	0x61, 0x69, 0x74, 0x46, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x54, 0x65, 0x72, 0x6d,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f,
	0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x73, 0x54, 0x65, 0x72, 0x6d,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30,
	0x01, 0x12, 0x2e, 0x0a, 0x05, 0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x10, 0x2e, 0x70, 0x62, 0x2e,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70,
	0x62, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x2b, 0x0a, 0x04, 0x53, 0x79, 0x6e, 0x63, 0x12, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x53,
	0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x70, 0x62, 0x2e,
	0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2e,
	0x0a, 0x05, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6c, 0x6f,
	0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x43,
	0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x05,
	0x5a, 0x03, 0x70, 0x62, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_networkFile_proto_rawDescOnce sync.Once
	file_networkFile_proto_rawDescData = file_networkFile_proto_rawDesc
)

func file_networkFile_proto_rawDescGZIP() []byte {
	file_networkFile_proto_rawDescOnce.Do(func() {
		file_networkFile_proto_rawDescData = protoimpl.X.CompressGZIP(file_networkFile_proto_rawDescData)
	})
	return file_networkFile_proto_rawDescData
}

var file_networkFile_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_networkFile_proto_goTypes = []any{
	(*OpenRequest)(nil),                     // 0: pb.OpenRequest
	(*OpenResponse)(nil),                    // 1: pb.OpenResponse
	(*WaitForServerTerminationRequest)(nil), // 2: pb.WaitForServerTerminationRequest
	(*ServerIsTerminatingResponse)(nil),     // 3: pb.ServerIsTerminatingResponse
	(*WriteRequest)(nil),                    // 4: pb.WriteRequest
	(*WriteResponse)(nil),                   // 5: pb.WriteResponse
	(*SyncRequest)(nil),                     // 6: pb.SyncRequest
	(*SyncResponse)(nil),                    // 7: pb.SyncResponse
	(*CloseRequest)(nil),                    // 8: pb.CloseRequest
	(*CloseResponse)(nil),                   // 9: pb.CloseResponse
	(*SliceKey)(nil),                        // 10: pb.SliceKey
	(*timestamppb.Timestamp)(nil),           // 11: google.protobuf.Timestamp
}
var file_networkFile_proto_depIdxs = []int32{
	11, // 0: pb.SliceKey.file_id:type_name -> google.protobuf.Timestamp
	11, // 1: pb.SliceKey.slice_id:type_name -> google.protobuf.Timestamp
	0,  // 2: pb.NetworkFile.Open:input_type -> pb.OpenRequest
	2,  // 3: pb.NetworkFile.WaitForServerTermination:input_type -> pb.WaitForServerTerminationRequest
	4,  // 4: pb.NetworkFile.Write:input_type -> pb.WriteRequest
	6,  // 5: pb.NetworkFile.Sync:input_type -> pb.SyncRequest
	8,  // 6: pb.NetworkFile.Close:input_type -> pb.CloseRequest
	1,  // 7: pb.NetworkFile.Open:output_type -> pb.OpenResponse
	3,  // 8: pb.NetworkFile.WaitForServerTermination:output_type -> pb.ServerIsTerminatingResponse
	5,  // 9: pb.NetworkFile.Write:output_type -> pb.WriteResponse
	7,  // 10: pb.NetworkFile.Sync:output_type -> pb.SyncResponse
	9,  // 11: pb.NetworkFile.Close:output_type -> pb.CloseResponse
	7,  // [7:12] is the sub-list for method output_type
	2,  // [2:7] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_networkFile_proto_init() }
func file_networkFile_proto_init() {
	if File_networkFile_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_networkFile_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*OpenRequest); i {
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
		file_networkFile_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*OpenResponse); i {
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
		file_networkFile_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*WaitForServerTerminationRequest); i {
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
		file_networkFile_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ServerIsTerminatingResponse); i {
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
		file_networkFile_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*WriteRequest); i {
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
		file_networkFile_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*WriteResponse); i {
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
		file_networkFile_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*SyncRequest); i {
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
		file_networkFile_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*SyncResponse); i {
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
		file_networkFile_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*CloseRequest); i {
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
		file_networkFile_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*CloseResponse); i {
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
		file_networkFile_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*SliceKey); i {
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
			RawDescriptor: file_networkFile_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_networkFile_proto_goTypes,
		DependencyIndexes: file_networkFile_proto_depIdxs,
		MessageInfos:      file_networkFile_proto_msgTypes,
	}.Build()
	File_networkFile_proto = out.File
	file_networkFile_proto_rawDesc = nil
	file_networkFile_proto_goTypes = nil
	file_networkFile_proto_depIdxs = nil
}
