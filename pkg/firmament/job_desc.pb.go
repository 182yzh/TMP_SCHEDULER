/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: job_desc.proto

package firmament

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type JobDescriptor_JobState int32

const (
	JobDescriptor_NEW       JobDescriptor_JobState = 0
	JobDescriptor_CREATED   JobDescriptor_JobState = 1
	JobDescriptor_RUNNING   JobDescriptor_JobState = 2
	JobDescriptor_COMPLETED JobDescriptor_JobState = 3
	JobDescriptor_FAILED    JobDescriptor_JobState = 4
	JobDescriptor_ABORTED   JobDescriptor_JobState = 5
	JobDescriptor_UNKNOWN   JobDescriptor_JobState = 6
)

var JobDescriptor_JobState_name = map[int32]string{
	0: "NEW",
	1: "CREATED",
	2: "RUNNING",
	3: "COMPLETED",
	4: "FAILED",
	5: "ABORTED",
	6: "UNKNOWN",
}
var JobDescriptor_JobState_value = map[string]int32{
	"NEW":       0,
	"CREATED":   1,
	"RUNNING":   2,
	"COMPLETED": 3,
	"FAILED":    4,
	"ABORTED":   5,
	"UNKNOWN":   6,
}

func (x JobDescriptor_JobState) String() string {
	return proto.EnumName(JobDescriptor_JobState_name, int32(x))
}
func (JobDescriptor_JobState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_job_desc_32b15611d30ec88a, []int{0, 0}
}

type JobDescriptor struct {
	Uuid      string                 `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Name      string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	State     JobDescriptor_JobState `protobuf:"varint,3,opt,name=state,proto3,enum=firmament.JobDescriptor_JobState" json:"state,omitempty"`
	RootTask  *TaskDescriptor        `protobuf:"bytes,4,opt,name=root_task,json=rootTask,proto3" json:"root_task,omitempty"`
	OutputIds [][]byte               `protobuf:"bytes,5,rep,name=output_ids,json=outputIds,proto3" json:"output_ids,omitempty"`
	// For gang scheduling jobs.
	MinNumberOfTasks     uint64   `protobuf:"varint,6,opt,name=min_number_of_tasks,json=minNumberOfTasks,proto3" json:"min_number_of_tasks,omitempty"`
	ScheduledTasksCount  uint64   `protobuf:"varint,7,opt,name=scheduled_tasks_count,json=scheduledTasksCount,proto3" json:"scheduled_tasks_count,omitempty"`
	IsGangSchedulingJob  bool     `protobuf:"varint,8,opt,name=is_gang_scheduling_job,json=isGangSchedulingJob,proto3" json:"is_gang_scheduling_job,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JobDescriptor) Reset()         { *m = JobDescriptor{} }
func (m *JobDescriptor) String() string { return proto.CompactTextString(m) }
func (*JobDescriptor) ProtoMessage()    {}
func (*JobDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_job_desc_32b15611d30ec88a, []int{0}
}
func (m *JobDescriptor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JobDescriptor.Unmarshal(m, b)
}
func (m *JobDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JobDescriptor.Marshal(b, m, deterministic)
}
func (dst *JobDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JobDescriptor.Merge(dst, src)
}
func (m *JobDescriptor) XXX_Size() int {
	return xxx_messageInfo_JobDescriptor.Size(m)
}
func (m *JobDescriptor) XXX_DiscardUnknown() {
	xxx_messageInfo_JobDescriptor.DiscardUnknown(m)
}

var xxx_messageInfo_JobDescriptor proto.InternalMessageInfo

func (m *JobDescriptor) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *JobDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *JobDescriptor) GetState() JobDescriptor_JobState {
	if m != nil {
		return m.State
	}
	return JobDescriptor_NEW
}

func (m *JobDescriptor) GetRootTask() *TaskDescriptor {
	if m != nil {
		return m.RootTask
	}
	return nil
}

func (m *JobDescriptor) GetOutputIds() [][]byte {
	if m != nil {
		return m.OutputIds
	}
	return nil
}

func (m *JobDescriptor) GetMinNumberOfTasks() uint64 {
	if m != nil {
		return m.MinNumberOfTasks
	}
	return 0
}

func (m *JobDescriptor) GetScheduledTasksCount() uint64 {
	if m != nil {
		return m.ScheduledTasksCount
	}
	return 0
}

func (m *JobDescriptor) GetIsGangSchedulingJob() bool {
	if m != nil {
		return m.IsGangSchedulingJob
	}
	return false
}

func init() {
	proto.RegisterType((*JobDescriptor)(nil), "firmament.JobDescriptor")
	proto.RegisterEnum("firmament.JobDescriptor_JobState", JobDescriptor_JobState_name, JobDescriptor_JobState_value)
}

func init() { proto.RegisterFile("job_desc.proto", fileDescriptor_job_desc_32b15611d30ec88a) }

var fileDescriptor_job_desc_32b15611d30ec88a = []byte{
	// 369 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0xcf, 0x8e, 0xd3, 0x30,
	0x10, 0x87, 0xc9, 0x26, 0x4d, 0x93, 0x59, 0x76, 0x89, 0x5c, 0x81, 0x0c, 0x12, 0x52, 0xd8, 0x53,
	0x2e, 0xf4, 0xd0, 0x95, 0xe0, 0x5c, 0xda, 0xb0, 0xea, 0xb2, 0x38, 0xc8, 0x6d, 0xd5, 0xa3, 0x95,
	0x3f, 0x6e, 0x70, 0x4b, 0xec, 0x2a, 0x76, 0x9e, 0x90, 0x17, 0x43, 0x4e, 0xa0, 0x74, 0x6f, 0x33,
	0xbf, 0xef, 0x9b, 0xb1, 0x65, 0xc3, 0xed, 0x41, 0x15, 0xac, 0xe2, 0xba, 0x9c, 0x9e, 0x5a, 0x65,
	0x14, 0x0a, 0xf7, 0xa2, 0x6d, 0xf2, 0x86, 0x4b, 0xf3, 0xee, 0x95, 0xc9, 0xf5, 0xf1, 0x82, 0xdd,
	0xfd, 0x76, 0xe1, 0xe6, 0x51, 0x15, 0x4b, 0xae, 0xcb, 0x56, 0x9c, 0x8c, 0x6a, 0x11, 0x02, 0xaf,
	0xeb, 0x44, 0x85, 0x9d, 0xd8, 0x49, 0x42, 0xda, 0xd7, 0x36, 0x93, 0x79, 0xc3, 0xf1, 0xd5, 0x90,
	0xd9, 0x1a, 0x7d, 0x86, 0x91, 0x36, 0xb9, 0xe1, 0xd8, 0x8d, 0x9d, 0xe4, 0x76, 0xf6, 0x61, 0x7a,
	0x3e, 0x65, 0xfa, 0x6c, 0xa1, 0xed, 0xd6, 0x56, 0xa4, 0x83, 0x8f, 0x3e, 0x41, 0xd8, 0x2a, 0x65,
	0x98, 0xbd, 0x0a, 0xf6, 0x62, 0x27, 0xb9, 0x9e, 0xbd, 0xbd, 0x18, 0xde, 0xe4, 0xfa, 0xf8, 0x7f,
	0x9a, 0x06, 0xd6, 0xb5, 0x19, 0x7a, 0x0f, 0xa0, 0x3a, 0x73, 0xea, 0x0c, 0x13, 0x95, 0xc6, 0xa3,
	0xd8, 0x4d, 0x5e, 0xd2, 0x70, 0x48, 0x56, 0x95, 0x46, 0x1f, 0x61, 0xd2, 0x08, 0xc9, 0x64, 0xd7,
	0x14, 0xbc, 0x65, 0x6a, 0xdf, 0xef, 0xd7, 0xd8, 0x8f, 0x9d, 0xc4, 0xa3, 0x51, 0x23, 0x24, 0xe9,
	0x49, 0xb6, 0xb7, 0xcb, 0x34, 0x9a, 0xc1, 0x6b, 0x5d, 0xfe, 0xe4, 0x55, 0xf7, 0x8b, 0x57, 0x83,
	0xca, 0x4a, 0xd5, 0x49, 0x83, 0xc7, 0xfd, 0xc0, 0xe4, 0x0c, 0x7b, 0x7d, 0x61, 0x11, 0xba, 0x87,
	0x37, 0x42, 0xb3, 0x3a, 0x97, 0x35, 0xfb, 0x8b, 0x85, 0xac, 0xd9, 0x41, 0x15, 0x38, 0x88, 0x9d,
	0x24, 0xa0, 0x13, 0xa1, 0x1f, 0x72, 0x59, 0xaf, 0xcf, 0xec, 0x51, 0x15, 0x77, 0x05, 0x04, 0xff,
	0x5e, 0x00, 0x8d, 0xc1, 0x25, 0xe9, 0x2e, 0x7a, 0x81, 0xae, 0x61, 0xbc, 0xa0, 0xe9, 0x7c, 0x93,
	0x2e, 0x23, 0xc7, 0x36, 0x74, 0x4b, 0xc8, 0x8a, 0x3c, 0x44, 0x57, 0xe8, 0x06, 0xc2, 0x45, 0xf6,
	0xfd, 0xc7, 0x53, 0x6a, 0x99, 0x8b, 0x00, 0xfc, 0xaf, 0xf3, 0xd5, 0x53, 0xba, 0x8c, 0x3c, 0xeb,
	0xcd, 0xbf, 0x64, 0xd4, 0x82, 0x91, 0x6d, 0xb6, 0xe4, 0x1b, 0xc9, 0x76, 0x24, 0xf2, 0x0b, 0xbf,
	0xff, 0xcc, 0xfb, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x80, 0x38, 0x6b, 0x10, 0xfa, 0x01, 0x00,
	0x00,
}
