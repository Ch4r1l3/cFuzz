// Code generated by protoc-gen-go. DO NOT EDIT.
// source: fuzzer_interface.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ErrorResponse struct {
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorResponse) Reset()         { *m = ErrorResponse{} }
func (m *ErrorResponse) String() string { return proto.CompactTextString(m) }
func (*ErrorResponse) ProtoMessage()    {}
func (*ErrorResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_41ff9b483743270f, []int{0}
}

func (m *ErrorResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrorResponse.Unmarshal(m, b)
}
func (m *ErrorResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrorResponse.Marshal(b, m, deterministic)
}
func (m *ErrorResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorResponse.Merge(m, src)
}
func (m *ErrorResponse) XXX_Size() int {
	return xxx_messageInfo_ErrorResponse.Size(m)
}
func (m *ErrorResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorResponse proto.InternalMessageInfo

func (m *ErrorResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type StringMap struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StringMap) Reset()         { *m = StringMap{} }
func (m *StringMap) String() string { return proto.CompactTextString(m) }
func (*StringMap) ProtoMessage()    {}
func (*StringMap) Descriptor() ([]byte, []int) {
	return fileDescriptor_41ff9b483743270f, []int{1}
}

func (m *StringMap) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StringMap.Unmarshal(m, b)
}
func (m *StringMap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StringMap.Marshal(b, m, deterministic)
}
func (m *StringMap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StringMap.Merge(m, src)
}
func (m *StringMap) XXX_Size() int {
	return xxx_messageInfo_StringMap.Size(m)
}
func (m *StringMap) XXX_DiscardUnknown() {
	xxx_messageInfo_StringMap.DiscardUnknown(m)
}

var xxx_messageInfo_StringMap proto.InternalMessageInfo

func (m *StringMap) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *StringMap) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type PrepareArg struct {
	CorpusDir            string            `protobuf:"bytes,1,opt,name=corpusDir,proto3" json:"corpusDir,omitempty"`
	TargetPath           string            `protobuf:"bytes,2,opt,name=targetPath,proto3" json:"targetPath,omitempty"`
	Arguments            map[string]string `protobuf:"bytes,3,rep,name=arguments,proto3" json:"arguments,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Environments         []string          `protobuf:"bytes,4,rep,name=environments,proto3" json:"environments,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *PrepareArg) Reset()         { *m = PrepareArg{} }
func (m *PrepareArg) String() string { return proto.CompactTextString(m) }
func (*PrepareArg) ProtoMessage()    {}
func (*PrepareArg) Descriptor() ([]byte, []int) {
	return fileDescriptor_41ff9b483743270f, []int{2}
}

func (m *PrepareArg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrepareArg.Unmarshal(m, b)
}
func (m *PrepareArg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrepareArg.Marshal(b, m, deterministic)
}
func (m *PrepareArg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrepareArg.Merge(m, src)
}
func (m *PrepareArg) XXX_Size() int {
	return xxx_messageInfo_PrepareArg.Size(m)
}
func (m *PrepareArg) XXX_DiscardUnknown() {
	xxx_messageInfo_PrepareArg.DiscardUnknown(m)
}

var xxx_messageInfo_PrepareArg proto.InternalMessageInfo

func (m *PrepareArg) GetCorpusDir() string {
	if m != nil {
		return m.CorpusDir
	}
	return ""
}

func (m *PrepareArg) GetTargetPath() string {
	if m != nil {
		return m.TargetPath
	}
	return ""
}

func (m *PrepareArg) GetArguments() map[string]string {
	if m != nil {
		return m.Arguments
	}
	return nil
}

func (m *PrepareArg) GetEnvironments() []string {
	if m != nil {
		return m.Environments
	}
	return nil
}

type FuzzArg struct {
	MaxTime              int32    `protobuf:"varint,1,opt,name=maxTime,proto3" json:"maxTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FuzzArg) Reset()         { *m = FuzzArg{} }
func (m *FuzzArg) String() string { return proto.CompactTextString(m) }
func (*FuzzArg) ProtoMessage()    {}
func (*FuzzArg) Descriptor() ([]byte, []int) {
	return fileDescriptor_41ff9b483743270f, []int{3}
}

func (m *FuzzArg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FuzzArg.Unmarshal(m, b)
}
func (m *FuzzArg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FuzzArg.Marshal(b, m, deterministic)
}
func (m *FuzzArg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FuzzArg.Merge(m, src)
}
func (m *FuzzArg) XXX_Size() int {
	return xxx_messageInfo_FuzzArg.Size(m)
}
func (m *FuzzArg) XXX_DiscardUnknown() {
	xxx_messageInfo_FuzzArg.DiscardUnknown(m)
}

var xxx_messageInfo_FuzzArg proto.InternalMessageInfo

func (m *FuzzArg) GetMaxTime() int32 {
	if m != nil {
		return m.MaxTime
	}
	return 0
}

type Crash struct {
	InputPath            string   `protobuf:"bytes,1,opt,name=inputPath,proto3" json:"inputPath,omitempty"`
	ReproduceArg         []string `protobuf:"bytes,2,rep,name=reproduceArg,proto3" json:"reproduceArg,omitempty"`
	Environments         []string `protobuf:"bytes,3,rep,name=environments,proto3" json:"environments,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Crash) Reset()         { *m = Crash{} }
func (m *Crash) String() string { return proto.CompactTextString(m) }
func (*Crash) ProtoMessage()    {}
func (*Crash) Descriptor() ([]byte, []int) {
	return fileDescriptor_41ff9b483743270f, []int{4}
}

func (m *Crash) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Crash.Unmarshal(m, b)
}
func (m *Crash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Crash.Marshal(b, m, deterministic)
}
func (m *Crash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Crash.Merge(m, src)
}
func (m *Crash) XXX_Size() int {
	return xxx_messageInfo_Crash.Size(m)
}
func (m *Crash) XXX_DiscardUnknown() {
	xxx_messageInfo_Crash.DiscardUnknown(m)
}

var xxx_messageInfo_Crash proto.InternalMessageInfo

func (m *Crash) GetInputPath() string {
	if m != nil {
		return m.InputPath
	}
	return ""
}

func (m *Crash) GetReproduceArg() []string {
	if m != nil {
		return m.ReproduceArg
	}
	return nil
}

func (m *Crash) GetEnvironments() []string {
	if m != nil {
		return m.Environments
	}
	return nil
}

type FuzzResult struct {
	Command              []string          `protobuf:"bytes,1,rep,name=command,proto3" json:"command,omitempty"`
	Crashes              []*Crash          `protobuf:"bytes,2,rep,name=crashes,proto3" json:"crashes,omitempty"`
	Stats                map[string]string `protobuf:"bytes,3,rep,name=stats,proto3" json:"stats,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TimeExecuted         int32             `protobuf:"varint,4,opt,name=timeExecuted,proto3" json:"timeExecuted,omitempty"`
	Error                string            `protobuf:"bytes,5,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *FuzzResult) Reset()         { *m = FuzzResult{} }
func (m *FuzzResult) String() string { return proto.CompactTextString(m) }
func (*FuzzResult) ProtoMessage()    {}
func (*FuzzResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_41ff9b483743270f, []int{5}
}

func (m *FuzzResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FuzzResult.Unmarshal(m, b)
}
func (m *FuzzResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FuzzResult.Marshal(b, m, deterministic)
}
func (m *FuzzResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FuzzResult.Merge(m, src)
}
func (m *FuzzResult) XXX_Size() int {
	return xxx_messageInfo_FuzzResult.Size(m)
}
func (m *FuzzResult) XXX_DiscardUnknown() {
	xxx_messageInfo_FuzzResult.DiscardUnknown(m)
}

var xxx_messageInfo_FuzzResult proto.InternalMessageInfo

func (m *FuzzResult) GetCommand() []string {
	if m != nil {
		return m.Command
	}
	return nil
}

func (m *FuzzResult) GetCrashes() []*Crash {
	if m != nil {
		return m.Crashes
	}
	return nil
}

func (m *FuzzResult) GetStats() map[string]string {
	if m != nil {
		return m.Stats
	}
	return nil
}

func (m *FuzzResult) GetTimeExecuted() int32 {
	if m != nil {
		return m.TimeExecuted
	}
	return 0
}

func (m *FuzzResult) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type ReproduceArg struct {
	InputPath            string   `protobuf:"bytes,1,opt,name=inputPath,proto3" json:"inputPath,omitempty"`
	MaxTime              int32    `protobuf:"varint,2,opt,name=maxTime,proto3" json:"maxTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReproduceArg) Reset()         { *m = ReproduceArg{} }
func (m *ReproduceArg) String() string { return proto.CompactTextString(m) }
func (*ReproduceArg) ProtoMessage()    {}
func (*ReproduceArg) Descriptor() ([]byte, []int) {
	return fileDescriptor_41ff9b483743270f, []int{6}
}

func (m *ReproduceArg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReproduceArg.Unmarshal(m, b)
}
func (m *ReproduceArg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReproduceArg.Marshal(b, m, deterministic)
}
func (m *ReproduceArg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReproduceArg.Merge(m, src)
}
func (m *ReproduceArg) XXX_Size() int {
	return xxx_messageInfo_ReproduceArg.Size(m)
}
func (m *ReproduceArg) XXX_DiscardUnknown() {
	xxx_messageInfo_ReproduceArg.DiscardUnknown(m)
}

var xxx_messageInfo_ReproduceArg proto.InternalMessageInfo

func (m *ReproduceArg) GetInputPath() string {
	if m != nil {
		return m.InputPath
	}
	return ""
}

func (m *ReproduceArg) GetMaxTime() int32 {
	if m != nil {
		return m.MaxTime
	}
	return 0
}

type ReproduceResult struct {
	Command              []string `protobuf:"bytes,1,rep,name=command,proto3" json:"command,omitempty"`
	ReturnCode           int32    `protobuf:"varint,2,opt,name=returnCode,proto3" json:"returnCode,omitempty"`
	TimeExecuted         int32    `protobuf:"varint,3,opt,name=timeExecuted,proto3" json:"timeExecuted,omitempty"`
	Output               []string `protobuf:"bytes,4,rep,name=output,proto3" json:"output,omitempty"`
	Error                string   `protobuf:"bytes,5,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReproduceResult) Reset()         { *m = ReproduceResult{} }
func (m *ReproduceResult) String() string { return proto.CompactTextString(m) }
func (*ReproduceResult) ProtoMessage()    {}
func (*ReproduceResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_41ff9b483743270f, []int{7}
}

func (m *ReproduceResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReproduceResult.Unmarshal(m, b)
}
func (m *ReproduceResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReproduceResult.Marshal(b, m, deterministic)
}
func (m *ReproduceResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReproduceResult.Merge(m, src)
}
func (m *ReproduceResult) XXX_Size() int {
	return xxx_messageInfo_ReproduceResult.Size(m)
}
func (m *ReproduceResult) XXX_DiscardUnknown() {
	xxx_messageInfo_ReproduceResult.DiscardUnknown(m)
}

var xxx_messageInfo_ReproduceResult proto.InternalMessageInfo

func (m *ReproduceResult) GetCommand() []string {
	if m != nil {
		return m.Command
	}
	return nil
}

func (m *ReproduceResult) GetReturnCode() int32 {
	if m != nil {
		return m.ReturnCode
	}
	return 0
}

func (m *ReproduceResult) GetTimeExecuted() int32 {
	if m != nil {
		return m.TimeExecuted
	}
	return 0
}

func (m *ReproduceResult) GetOutput() []string {
	if m != nil {
		return m.Output
	}
	return nil
}

func (m *ReproduceResult) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type MinimizeCorpusArg struct {
	InputDir             string   `protobuf:"bytes,1,opt,name=inputDir,proto3" json:"inputDir,omitempty"`
	OutputDir            string   `protobuf:"bytes,2,opt,name=outputDir,proto3" json:"outputDir,omitempty"`
	MaxTime              int32    `protobuf:"varint,3,opt,name=maxTime,proto3" json:"maxTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MinimizeCorpusArg) Reset()         { *m = MinimizeCorpusArg{} }
func (m *MinimizeCorpusArg) String() string { return proto.CompactTextString(m) }
func (*MinimizeCorpusArg) ProtoMessage()    {}
func (*MinimizeCorpusArg) Descriptor() ([]byte, []int) {
	return fileDescriptor_41ff9b483743270f, []int{8}
}

func (m *MinimizeCorpusArg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MinimizeCorpusArg.Unmarshal(m, b)
}
func (m *MinimizeCorpusArg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MinimizeCorpusArg.Marshal(b, m, deterministic)
}
func (m *MinimizeCorpusArg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MinimizeCorpusArg.Merge(m, src)
}
func (m *MinimizeCorpusArg) XXX_Size() int {
	return xxx_messageInfo_MinimizeCorpusArg.Size(m)
}
func (m *MinimizeCorpusArg) XXX_DiscardUnknown() {
	xxx_messageInfo_MinimizeCorpusArg.DiscardUnknown(m)
}

var xxx_messageInfo_MinimizeCorpusArg proto.InternalMessageInfo

func (m *MinimizeCorpusArg) GetInputDir() string {
	if m != nil {
		return m.InputDir
	}
	return ""
}

func (m *MinimizeCorpusArg) GetOutputDir() string {
	if m != nil {
		return m.OutputDir
	}
	return ""
}

func (m *MinimizeCorpusArg) GetMaxTime() int32 {
	if m != nil {
		return m.MaxTime
	}
	return 0
}

type MinimizeCorpusResult struct {
	Command              []string          `protobuf:"bytes,1,rep,name=command,proto3" json:"command,omitempty"`
	Stats                map[string]string `protobuf:"bytes,2,rep,name=stats,proto3" json:"stats,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TimeExecuted         int32             `protobuf:"varint,3,opt,name=timeExecuted,proto3" json:"timeExecuted,omitempty"`
	Error                string            `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *MinimizeCorpusResult) Reset()         { *m = MinimizeCorpusResult{} }
func (m *MinimizeCorpusResult) String() string { return proto.CompactTextString(m) }
func (*MinimizeCorpusResult) ProtoMessage()    {}
func (*MinimizeCorpusResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_41ff9b483743270f, []int{9}
}

func (m *MinimizeCorpusResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MinimizeCorpusResult.Unmarshal(m, b)
}
func (m *MinimizeCorpusResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MinimizeCorpusResult.Marshal(b, m, deterministic)
}
func (m *MinimizeCorpusResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MinimizeCorpusResult.Merge(m, src)
}
func (m *MinimizeCorpusResult) XXX_Size() int {
	return xxx_messageInfo_MinimizeCorpusResult.Size(m)
}
func (m *MinimizeCorpusResult) XXX_DiscardUnknown() {
	xxx_messageInfo_MinimizeCorpusResult.DiscardUnknown(m)
}

var xxx_messageInfo_MinimizeCorpusResult proto.InternalMessageInfo

func (m *MinimizeCorpusResult) GetCommand() []string {
	if m != nil {
		return m.Command
	}
	return nil
}

func (m *MinimizeCorpusResult) GetStats() map[string]string {
	if m != nil {
		return m.Stats
	}
	return nil
}

func (m *MinimizeCorpusResult) GetTimeExecuted() int32 {
	if m != nil {
		return m.TimeExecuted
	}
	return 0
}

func (m *MinimizeCorpusResult) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_41ff9b483743270f, []int{10}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ErrorResponse)(nil), "proto.ErrorResponse")
	proto.RegisterType((*StringMap)(nil), "proto.StringMap")
	proto.RegisterType((*PrepareArg)(nil), "proto.PrepareArg")
	proto.RegisterMapType((map[string]string)(nil), "proto.PrepareArg.ArgumentsEntry")
	proto.RegisterType((*FuzzArg)(nil), "proto.FuzzArg")
	proto.RegisterType((*Crash)(nil), "proto.Crash")
	proto.RegisterType((*FuzzResult)(nil), "proto.FuzzResult")
	proto.RegisterMapType((map[string]string)(nil), "proto.FuzzResult.StatsEntry")
	proto.RegisterType((*ReproduceArg)(nil), "proto.ReproduceArg")
	proto.RegisterType((*ReproduceResult)(nil), "proto.ReproduceResult")
	proto.RegisterType((*MinimizeCorpusArg)(nil), "proto.MinimizeCorpusArg")
	proto.RegisterType((*MinimizeCorpusResult)(nil), "proto.MinimizeCorpusResult")
	proto.RegisterMapType((map[string]string)(nil), "proto.MinimizeCorpusResult.StatsEntry")
	proto.RegisterType((*Empty)(nil), "proto.Empty")
}

func init() { proto.RegisterFile("fuzzer_interface.proto", fileDescriptor_41ff9b483743270f) }

var fileDescriptor_41ff9b483743270f = []byte{
	// 612 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x41, 0x6b, 0xdb, 0x4c,
	0x10, 0x45, 0x96, 0x15, 0x7f, 0x9a, 0xe4, 0x4b, 0xeb, 0xad, 0x31, 0xc2, 0x0d, 0xc1, 0xa8, 0x34,
	0xb8, 0x14, 0x7c, 0x70, 0x2e, 0xa6, 0x84, 0x42, 0x70, 0xed, 0x9e, 0x02, 0x41, 0xe9, 0xbd, 0x6c,
	0xe5, 0x89, 0x23, 0x6a, 0xad, 0xc4, 0x6a, 0x37, 0xc4, 0xfe, 0x23, 0x3d, 0xf5, 0xcf, 0xf5, 0xdc,
	0x73, 0xcf, 0x65, 0x57, 0x2b, 0x4b, 0xaa, 0x8d, 0x93, 0x42, 0x4f, 0xd2, 0xbc, 0x9d, 0x99, 0x9d,
	0xf7, 0xf6, 0xed, 0x42, 0xf7, 0x56, 0xae, 0xd7, 0xc8, 0x3f, 0x47, 0x4c, 0x20, 0xbf, 0xa5, 0x21,
	0x0e, 0x53, 0x9e, 0x88, 0x84, 0x38, 0xfa, 0xe3, 0xbf, 0x86, 0xff, 0xa7, 0x9c, 0x27, 0x3c, 0xc0,
	0x2c, 0x4d, 0x58, 0x86, 0xa4, 0x03, 0x0e, 0x2a, 0xc0, 0xb3, 0xfa, 0xd6, 0xc0, 0x0d, 0xf2, 0xc0,
	0x3f, 0x07, 0xf7, 0x46, 0xf0, 0x88, 0x2d, 0xae, 0x68, 0x4a, 0x9e, 0x83, 0xfd, 0x15, 0x57, 0x26,
	0x41, 0xfd, 0xaa, 0xa2, 0x7b, 0xba, 0x94, 0xe8, 0x35, 0xf2, 0x22, 0x1d, 0xf8, 0x3f, 0x2d, 0x80,
	0x6b, 0x8e, 0x29, 0xe5, 0x78, 0xc9, 0x17, 0xe4, 0x04, 0xdc, 0x30, 0xe1, 0xa9, 0xcc, 0x3e, 0x44,
	0x45, 0xf7, 0x12, 0x20, 0xa7, 0x00, 0x82, 0xf2, 0x05, 0x8a, 0x6b, 0x2a, 0xee, 0x4c, 0x9f, 0x0a,
	0x42, 0xde, 0x83, 0x4b, 0xf9, 0x42, 0xc6, 0xc8, 0x44, 0xe6, 0xd9, 0x7d, 0x7b, 0x70, 0x38, 0xea,
	0xe7, 0x54, 0x86, 0xe5, 0x1e, 0xc3, 0xcb, 0x22, 0x65, 0xca, 0x04, 0x5f, 0x05, 0x65, 0x09, 0xf1,
	0xe1, 0x08, 0xd9, 0x7d, 0xc4, 0x13, 0x96, 0xb7, 0x68, 0xf6, 0xed, 0x81, 0x1b, 0xd4, 0xb0, 0xde,
	0x05, 0x1c, 0xd7, 0x1b, 0x3c, 0x95, 0xea, 0xbb, 0xc6, 0xd8, 0xf2, 0x5f, 0x41, 0x6b, 0x26, 0xd7,
	0x6b, 0x45, 0xd5, 0x83, 0x56, 0x4c, 0x1f, 0x3e, 0x45, 0x31, 0xea, 0x52, 0x27, 0x28, 0x42, 0x3f,
	0x06, 0x67, 0xc2, 0x69, 0x76, 0xa7, 0xd4, 0x88, 0x58, 0x2a, 0x73, 0xba, 0x46, 0x8d, 0x0d, 0xa0,
	0xa6, 0xe5, 0x98, 0xf2, 0x64, 0x2e, 0x43, 0xc5, 0xcb, 0x6b, 0xe4, 0xd3, 0x56, 0xb1, 0x2d, 0x46,
	0xf6, 0x36, 0x23, 0xff, 0x97, 0x05, 0xa0, 0x86, 0x0a, 0x30, 0x93, 0x4b, 0xa1, 0xe6, 0x0a, 0x93,
	0x38, 0xa6, 0x6c, 0xee, 0x59, 0x3a, 0xbb, 0x08, 0xc9, 0x19, 0xb4, 0x42, 0x35, 0x17, 0x66, 0x7a,
	0xaf, 0xc3, 0xd1, 0x91, 0x11, 0x57, 0x4f, 0x1b, 0x14, 0x8b, 0x64, 0x04, 0x4e, 0x26, 0xe8, 0xe6,
	0x08, 0x4e, 0x4c, 0x56, 0xb9, 0xc7, 0xf0, 0x46, 0x2d, 0xe7, 0xf2, 0xe7, 0xa9, 0x6a, 0x50, 0x11,
	0xc5, 0x38, 0x7d, 0xc0, 0x50, 0x0a, 0x9c, 0x7b, 0x4d, 0x2d, 0x49, 0x0d, 0x2b, 0x6d, 0xe7, 0x54,
	0x6c, 0xd7, 0x1b, 0x03, 0x94, 0xed, 0xfe, 0xea, 0x30, 0x66, 0x70, 0x14, 0x54, 0xc5, 0xda, 0x2f,
	0x77, 0xe5, 0xbc, 0x1a, 0xf5, 0xf3, 0xfa, 0x6e, 0xc1, 0xb3, 0x4d, 0xa3, 0x47, 0x55, 0x3c, 0x05,
	0xe0, 0x28, 0x24, 0x67, 0x93, 0x64, 0x5e, 0xb4, 0xaa, 0x20, 0x5b, 0x4a, 0xd8, 0x3b, 0x94, 0xe8,
	0xc2, 0x41, 0x22, 0x45, 0x2a, 0x85, 0xb1, 0xa8, 0x89, 0x76, 0x2b, 0xe4, 0x2f, 0xa0, 0x7d, 0x15,
	0xb1, 0x28, 0x8e, 0xd6, 0x38, 0xd1, 0x77, 0x49, 0x91, 0xed, 0xc1, 0x7f, 0x9a, 0x5b, 0x79, 0xd1,
	0x36, 0xb1, 0x12, 0x22, 0x6f, 0xa8, 0x16, 0x73, 0xd9, 0x4a, 0xa0, 0x2a, 0x84, 0x5d, 0x17, 0xe2,
	0x87, 0x05, 0x9d, 0xfa, 0x4e, 0x8f, 0xaa, 0x71, 0x51, 0x78, 0x25, 0x77, 0xd4, 0x99, 0xf1, 0xca,
	0xae, 0x2e, 0x4f, 0x70, 0x8d, 0xbd, 0xcf, 0x35, 0xcd, 0x7f, 0xe3, 0x9a, 0x16, 0x38, 0xd3, 0x38,
	0x15, 0xab, 0xd1, 0xb7, 0x06, 0x1c, 0xcc, 0xf4, 0xc3, 0x49, 0x46, 0xd0, 0x32, 0x0f, 0x0c, 0x69,
	0x6f, 0x3d, 0x38, 0xbd, 0x8e, 0x81, 0xea, 0x8f, 0xe8, 0x1b, 0x68, 0xaa, 0x6a, 0x72, 0x5c, 0xb9,
	0x1e, 0x2a, 0xbb, 0xbd, 0x75, 0x5d, 0xc8, 0x18, 0xdc, 0x8d, 0xbf, 0xc8, 0x0b, 0xb3, 0x5e, 0xb5,
	0x6e, 0xaf, 0xfb, 0x27, 0x68, 0x2a, 0x3f, 0xc2, 0x71, 0x5d, 0x4a, 0xe2, 0xed, 0x54, 0x58, 0xf5,
	0x78, 0xb9, 0x47, 0x7b, 0xf2, 0x16, 0x9c, 0xc9, 0x12, 0x29, 0x23, 0xc5, 0x9d, 0xd7, 0x1a, 0xec,
	0xa6, 0xf6, 0xe5, 0x40, 0x83, 0xe7, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x4f, 0xa3, 0x75, 0xcf,
	0x58, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FuzzerClient is the client API for Fuzzer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FuzzerClient interface {
	Prepare(ctx context.Context, in *PrepareArg, opts ...grpc.CallOption) (*ErrorResponse, error)
	Fuzz(ctx context.Context, in *FuzzArg, opts ...grpc.CallOption) (*FuzzResult, error)
	Reproduce(ctx context.Context, in *ReproduceArg, opts ...grpc.CallOption) (*ReproduceResult, error)
	MinimizeCorpus(ctx context.Context, in *MinimizeCorpusArg, opts ...grpc.CallOption) (*MinimizeCorpusResult, error)
	Clean(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ErrorResponse, error)
}

type fuzzerClient struct {
	cc *grpc.ClientConn
}

func NewFuzzerClient(cc *grpc.ClientConn) FuzzerClient {
	return &fuzzerClient{cc}
}

func (c *fuzzerClient) Prepare(ctx context.Context, in *PrepareArg, opts ...grpc.CallOption) (*ErrorResponse, error) {
	out := new(ErrorResponse)
	err := c.cc.Invoke(ctx, "/proto.Fuzzer/Prepare", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fuzzerClient) Fuzz(ctx context.Context, in *FuzzArg, opts ...grpc.CallOption) (*FuzzResult, error) {
	out := new(FuzzResult)
	err := c.cc.Invoke(ctx, "/proto.Fuzzer/Fuzz", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fuzzerClient) Reproduce(ctx context.Context, in *ReproduceArg, opts ...grpc.CallOption) (*ReproduceResult, error) {
	out := new(ReproduceResult)
	err := c.cc.Invoke(ctx, "/proto.Fuzzer/Reproduce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fuzzerClient) MinimizeCorpus(ctx context.Context, in *MinimizeCorpusArg, opts ...grpc.CallOption) (*MinimizeCorpusResult, error) {
	out := new(MinimizeCorpusResult)
	err := c.cc.Invoke(ctx, "/proto.Fuzzer/MinimizeCorpus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fuzzerClient) Clean(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ErrorResponse, error) {
	out := new(ErrorResponse)
	err := c.cc.Invoke(ctx, "/proto.Fuzzer/Clean", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FuzzerServer is the server API for Fuzzer service.
type FuzzerServer interface {
	Prepare(context.Context, *PrepareArg) (*ErrorResponse, error)
	Fuzz(context.Context, *FuzzArg) (*FuzzResult, error)
	Reproduce(context.Context, *ReproduceArg) (*ReproduceResult, error)
	MinimizeCorpus(context.Context, *MinimizeCorpusArg) (*MinimizeCorpusResult, error)
	Clean(context.Context, *Empty) (*ErrorResponse, error)
}

// UnimplementedFuzzerServer can be embedded to have forward compatible implementations.
type UnimplementedFuzzerServer struct {
}

func (*UnimplementedFuzzerServer) Prepare(ctx context.Context, req *PrepareArg) (*ErrorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Prepare not implemented")
}
func (*UnimplementedFuzzerServer) Fuzz(ctx context.Context, req *FuzzArg) (*FuzzResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fuzz not implemented")
}
func (*UnimplementedFuzzerServer) Reproduce(ctx context.Context, req *ReproduceArg) (*ReproduceResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reproduce not implemented")
}
func (*UnimplementedFuzzerServer) MinimizeCorpus(ctx context.Context, req *MinimizeCorpusArg) (*MinimizeCorpusResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MinimizeCorpus not implemented")
}
func (*UnimplementedFuzzerServer) Clean(ctx context.Context, req *Empty) (*ErrorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Clean not implemented")
}

func RegisterFuzzerServer(s *grpc.Server, srv FuzzerServer) {
	s.RegisterService(&_Fuzzer_serviceDesc, srv)
}

func _Fuzzer_Prepare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrepareArg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FuzzerServer).Prepare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Fuzzer/Prepare",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FuzzerServer).Prepare(ctx, req.(*PrepareArg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fuzzer_Fuzz_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FuzzArg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FuzzerServer).Fuzz(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Fuzzer/Fuzz",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FuzzerServer).Fuzz(ctx, req.(*FuzzArg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fuzzer_Reproduce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReproduceArg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FuzzerServer).Reproduce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Fuzzer/Reproduce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FuzzerServer).Reproduce(ctx, req.(*ReproduceArg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fuzzer_MinimizeCorpus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MinimizeCorpusArg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FuzzerServer).MinimizeCorpus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Fuzzer/MinimizeCorpus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FuzzerServer).MinimizeCorpus(ctx, req.(*MinimizeCorpusArg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fuzzer_Clean_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FuzzerServer).Clean(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Fuzzer/Clean",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FuzzerServer).Clean(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Fuzzer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Fuzzer",
	HandlerType: (*FuzzerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Prepare",
			Handler:    _Fuzzer_Prepare_Handler,
		},
		{
			MethodName: "Fuzz",
			Handler:    _Fuzzer_Fuzz_Handler,
		},
		{
			MethodName: "Reproduce",
			Handler:    _Fuzzer_Reproduce_Handler,
		},
		{
			MethodName: "MinimizeCorpus",
			Handler:    _Fuzzer_MinimizeCorpus_Handler,
		},
		{
			MethodName: "Clean",
			Handler:    _Fuzzer_Clean_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fuzzer_interface.proto",
}
