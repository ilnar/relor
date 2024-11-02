// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.21.12
// source: api/job.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListenOperation int32

const (
	ListenOperation_UNSPECIFIED ListenOperation = 0
	ListenOperation_CREATE      ListenOperation = 1
	ListenOperation_COMPLETE    ListenOperation = 2
)

// Enum value maps for ListenOperation.
var (
	ListenOperation_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "CREATE",
		2: "COMPLETE",
	}
	ListenOperation_value = map[string]int32{
		"UNSPECIFIED": 0,
		"CREATE":      1,
		"COMPLETE":    2,
	}
)

func (x ListenOperation) Enum() *ListenOperation {
	p := new(ListenOperation)
	*p = x
	return p
}

func (x ListenOperation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ListenOperation) Descriptor() protoreflect.EnumDescriptor {
	return file_api_job_proto_enumTypes[0].Descriptor()
}

func (ListenOperation) Type() protoreflect.EnumType {
	return &file_api_job_proto_enumTypes[0]
}

func (x ListenOperation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ListenOperation.Descriptor instead.
func (ListenOperation) EnumDescriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{0}
}

type Reference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkflowId     string `protobuf:"bytes,1,opt,name=workflow_id,json=workflowId,proto3" json:"workflow_id,omitempty"`
	WorkflowAction string `protobuf:"bytes,2,opt,name=workflow_action,json=workflowAction,proto3" json:"workflow_action,omitempty"`
	TransitionId   string `protobuf:"bytes,3,opt,name=transition_id,json=transitionId,proto3" json:"transition_id,omitempty"`
}

func (x *Reference) Reset() {
	*x = Reference{}
	mi := &file_api_job_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Reference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reference) ProtoMessage() {}

func (x *Reference) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reference.ProtoReflect.Descriptor instead.
func (*Reference) Descriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{0}
}

func (x *Reference) GetWorkflowId() string {
	if x != nil {
		return x.WorkflowId
	}
	return ""
}

func (x *Reference) GetWorkflowAction() string {
	if x != nil {
		return x.WorkflowAction
	}
	return ""
}

func (x *Reference) GetTransitionId() string {
	if x != nil {
		return x.TransitionId
	}
	return ""
}

type CreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateResponse) Reset() {
	*x = CreateResponse{}
	mi := &file_api_job_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateResponse) ProtoMessage() {}

func (x *CreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateResponse.ProtoReflect.Descriptor instead.
func (*CreateResponse) Descriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{1}
}

type CreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Reference    *Reference           `protobuf:"bytes,2,opt,name=reference,proto3" json:"reference,omitempty"`
	ResultLabels []string             `protobuf:"bytes,3,rep,name=result_labels,json=resultLabels,proto3" json:"result_labels,omitempty"`
	Ttl          *durationpb.Duration `protobuf:"bytes,4,opt,name=ttl,proto3" json:"ttl,omitempty"`
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	mi := &file_api_job_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{2}
}

func (x *CreateRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CreateRequest) GetReference() *Reference {
	if x != nil {
		return x.Reference
	}
	return nil
}

func (x *CreateRequest) GetResultLabels() []string {
	if x != nil {
		return x.ResultLabels
	}
	return nil
}

func (x *CreateRequest) GetTtl() *durationpb.Duration {
	if x != nil {
		return x.Ttl
	}
	return nil
}

type ClaimRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ClaimRequest) Reset() {
	*x = ClaimRequest{}
	mi := &file_api_job_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClaimRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClaimRequest) ProtoMessage() {}

func (x *ClaimRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClaimRequest.ProtoReflect.Descriptor instead.
func (*ClaimRequest) Descriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{3}
}

func (x *ClaimRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ClaimResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActionId     string   `protobuf:"bytes,1,opt,name=action_id,json=actionId,proto3" json:"action_id,omitempty"`
	ResultLabels []string `protobuf:"bytes,2,rep,name=result_labels,json=resultLabels,proto3" json:"result_labels,omitempty"`
}

func (x *ClaimResponse) Reset() {
	*x = ClaimResponse{}
	mi := &file_api_job_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClaimResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClaimResponse) ProtoMessage() {}

func (x *ClaimResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClaimResponse.ProtoReflect.Descriptor instead.
func (*ClaimResponse) Descriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{4}
}

func (x *ClaimResponse) GetActionId() string {
	if x != nil {
		return x.ActionId
	}
	return ""
}

func (x *ClaimResponse) GetResultLabels() []string {
	if x != nil {
		return x.ResultLabels
	}
	return nil
}

type ReleaseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ReleaseRequest) Reset() {
	*x = ReleaseRequest{}
	mi := &file_api_job_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReleaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReleaseRequest) ProtoMessage() {}

func (x *ReleaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReleaseRequest.ProtoReflect.Descriptor instead.
func (*ReleaseRequest) Descriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{5}
}

func (x *ReleaseRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ReleaseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReleaseResponse) Reset() {
	*x = ReleaseResponse{}
	mi := &file_api_job_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReleaseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReleaseResponse) ProtoMessage() {}

func (x *ReleaseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReleaseResponse.ProtoReflect.Descriptor instead.
func (*ReleaseResponse) Descriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{6}
}

type CompleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ResultLabel  string `protobuf:"bytes,2,opt,name=result_label,json=resultLabel,proto3" json:"result_label,omitempty"`
	TransitionId string `protobuf:"bytes,3,opt,name=transition_id,json=transitionId,proto3" json:"transition_id,omitempty"`
}

func (x *CompleteRequest) Reset() {
	*x = CompleteRequest{}
	mi := &file_api_job_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteRequest) ProtoMessage() {}

func (x *CompleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteRequest.ProtoReflect.Descriptor instead.
func (*CompleteRequest) Descriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{7}
}

func (x *CompleteRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CompleteRequest) GetResultLabel() string {
	if x != nil {
		return x.ResultLabel
	}
	return ""
}

func (x *CompleteRequest) GetTransitionId() string {
	if x != nil {
		return x.TransitionId
	}
	return ""
}

type CompleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CompleteResponse) Reset() {
	*x = CompleteResponse{}
	mi := &file_api_job_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteResponse) ProtoMessage() {}

func (x *CompleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteResponse.ProtoReflect.Descriptor instead.
func (*CompleteResponse) Descriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{8}
}

type ListenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkerId string          `protobuf:"bytes,1,opt,name=worker_id,json=workerId,proto3" json:"worker_id,omitempty"`
	ListenTo ListenOperation `protobuf:"varint,2,opt,name=listen_to,json=listenTo,proto3,enum=api.ListenOperation" json:"listen_to,omitempty"`
}

func (x *ListenRequest) Reset() {
	*x = ListenRequest{}
	mi := &file_api_job_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenRequest) ProtoMessage() {}

func (x *ListenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenRequest.ProtoReflect.Descriptor instead.
func (*ListenRequest) Descriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{9}
}

func (x *ListenRequest) GetWorkerId() string {
	if x != nil {
		return x.WorkerId
	}
	return ""
}

func (x *ListenRequest) GetListenTo() ListenOperation {
	if x != nil {
		return x.ListenTo
	}
	return ListenOperation_UNSPECIFIED
}

type Job struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Reference       *Reference `protobuf:"bytes,2,opt,name=reference,proto3" json:"reference,omitempty"`
	AvailableLabels []string   `protobuf:"bytes,3,rep,name=available_labels,json=availableLabels,proto3" json:"available_labels,omitempty"`
	ResultLabel     string     `protobuf:"bytes,4,opt,name=result_label,json=resultLabel,proto3" json:"result_label,omitempty"`
}

func (x *Job) Reset() {
	*x = Job{}
	mi := &file_api_job_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_api_job_proto_rawDescGZIP(), []int{10}
}

func (x *Job) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Job) GetReference() *Reference {
	if x != nil {
		return x.Reference
	}
	return nil
}

func (x *Job) GetAvailableLabels() []string {
	if x != nil {
		return x.AvailableLabels
	}
	return nil
}

func (x *Job) GetResultLabel() string {
	if x != nil {
		return x.ResultLabel
	}
	return ""
}

var File_api_job_proto protoreflect.FileDescriptor

var file_api_job_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x61, 0x70, 0x69, 0x2f, 0x6a, 0x6f, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x61, 0x70, 0x69, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7a, 0x0a, 0x09, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x77, 0x6f, 0x72,
	0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x22, 0x10, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x9f, 0x01, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x2c, 0x0a, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x2b, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x03, 0x74, 0x74, 0x6c, 0x22, 0x1e, 0x0a, 0x0c, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x51, 0x0a, 0x0d, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x22, 0x20, 0x0a, 0x0e, 0x52, 0x65, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x11, 0x0a, 0x0f, 0x52, 0x65, 0x6c,
	0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x69, 0x0a, 0x0f,
	0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x21, 0x0a, 0x0c, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x12, 0x0a, 0x10, 0x43, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x5f, 0x0a, 0x0d, 0x4c,
	0x69, 0x73, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09,
	0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x09, 0x6c, 0x69, 0x73,
	0x74, 0x65, 0x6e, 0x5f, 0x74, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x54, 0x6f, 0x22, 0x91, 0x01, 0x0a,
	0x03, 0x4a, 0x6f, 0x62, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x2c, 0x0a, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x5f,
	0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x76,
	0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x21, 0x0a,
	0x0c, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x2a, 0x3c, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x10, 0x01,
	0x12, 0x0c, 0x0a, 0x08, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x02, 0x32, 0x92,
	0x02, 0x0a, 0x0a, 0x4a, 0x6f, 0x62, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x33, 0x0a,
	0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x12, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x30, 0x0a, 0x05, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x12, 0x11, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x07, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12,
	0x13, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x08,
	0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43,
	0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2a, 0x0a, 0x06, 0x4c, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x12, 0x12, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x08, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4a, 0x6f, 0x62, 0x22,
	0x00, 0x30, 0x01, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x69, 0x6c, 0x6e, 0x61, 0x72, 0x2f, 0x77, 0x66, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70,
	0x62, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_job_proto_rawDescOnce sync.Once
	file_api_job_proto_rawDescData = file_api_job_proto_rawDesc
)

func file_api_job_proto_rawDescGZIP() []byte {
	file_api_job_proto_rawDescOnce.Do(func() {
		file_api_job_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_job_proto_rawDescData)
	})
	return file_api_job_proto_rawDescData
}

var file_api_job_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_job_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_api_job_proto_goTypes = []any{
	(ListenOperation)(0),        // 0: api.ListenOperation
	(*Reference)(nil),           // 1: api.Reference
	(*CreateResponse)(nil),      // 2: api.CreateResponse
	(*CreateRequest)(nil),       // 3: api.CreateRequest
	(*ClaimRequest)(nil),        // 4: api.ClaimRequest
	(*ClaimResponse)(nil),       // 5: api.ClaimResponse
	(*ReleaseRequest)(nil),      // 6: api.ReleaseRequest
	(*ReleaseResponse)(nil),     // 7: api.ReleaseResponse
	(*CompleteRequest)(nil),     // 8: api.CompleteRequest
	(*CompleteResponse)(nil),    // 9: api.CompleteResponse
	(*ListenRequest)(nil),       // 10: api.ListenRequest
	(*Job)(nil),                 // 11: api.Job
	(*durationpb.Duration)(nil), // 12: google.protobuf.Duration
}
var file_api_job_proto_depIdxs = []int32{
	1,  // 0: api.CreateRequest.reference:type_name -> api.Reference
	12, // 1: api.CreateRequest.ttl:type_name -> google.protobuf.Duration
	0,  // 2: api.ListenRequest.listen_to:type_name -> api.ListenOperation
	1,  // 3: api.Job.reference:type_name -> api.Reference
	3,  // 4: api.JobService.Create:input_type -> api.CreateRequest
	4,  // 5: api.JobService.Claim:input_type -> api.ClaimRequest
	6,  // 6: api.JobService.Release:input_type -> api.ReleaseRequest
	8,  // 7: api.JobService.Complete:input_type -> api.CompleteRequest
	10, // 8: api.JobService.Listen:input_type -> api.ListenRequest
	2,  // 9: api.JobService.Create:output_type -> api.CreateResponse
	5,  // 10: api.JobService.Claim:output_type -> api.ClaimResponse
	7,  // 11: api.JobService.Release:output_type -> api.ReleaseResponse
	9,  // 12: api.JobService.Complete:output_type -> api.CompleteResponse
	11, // 13: api.JobService.Listen:output_type -> api.Job
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_api_job_proto_init() }
func file_api_job_proto_init() {
	if File_api_job_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_job_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_job_proto_goTypes,
		DependencyIndexes: file_api_job_proto_depIdxs,
		EnumInfos:         file_api_job_proto_enumTypes,
		MessageInfos:      file_api_job_proto_msgTypes,
	}.Build()
	File_api_job_proto = out.File
	file_api_job_proto_rawDesc = nil
	file_api_job_proto_goTypes = nil
	file_api_job_proto_depIdxs = nil
}
