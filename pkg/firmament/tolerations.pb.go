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
// source: tolerations.proto

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

// The pod this Toleration is attached to tolerates any taint that matches
// the triple <key,value,effect> using the matching operator <operator>.
type Toleration struct {
	// Key is the taint key that the toleration applies to. Empty means match all taint keys.
	// If the key is empty, operator must be Exists; this combination means to match all values and all keys.
	// +optional
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Operator represents a key's relationship to the value.
	// Valid operators are Exists and Equal. Defaults to Equal.
	// Exists is equivalent to wildcard for value, so that a pod can
	// tolerate all taints of a particular category.
	// +optional
	Operator string `protobuf:"bytes,2,opt,name=operator,proto3" json:"operator,omitempty"`
	// Value is the taint value the toleration matches to.
	// If the operator is Exists, the value should be empty, otherwise just a regular string.
	// +optional
	Value string `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	// Effect indicates the taint effect to match. Empty means match all taint effects.
	// When specified, allowed values are NoSchedule, PreferNoSchedule and NoExecute.
	// +optional
	Effect string `protobuf:"bytes,4,opt,name=effect,proto3" json:"effect,omitempty"`
	// TolerationSeconds represents the period of time the toleration (which must be
	// of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default,
	// it is not set, which means tolerate the taint forever (do not evict). Zero and
	// negative values will be treated as 0 (evict immediately) by the system.
	// +optional
	TolerationSeconds    int64    `protobuf:"varint,5,opt,name=tolerationSeconds,proto3" json:"tolerationSeconds,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Toleration) Reset()         { *m = Toleration{} }
func (m *Toleration) String() string { return proto.CompactTextString(m) }
func (*Toleration) ProtoMessage()    {}
func (*Toleration) Descriptor() ([]byte, []int) {
	return fileDescriptor_tolerations_aa2c6496805744bf, []int{0}
}
func (m *Toleration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Toleration.Unmarshal(m, b)
}
func (m *Toleration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Toleration.Marshal(b, m, deterministic)
}
func (dst *Toleration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Toleration.Merge(dst, src)
}
func (m *Toleration) XXX_Size() int {
	return xxx_messageInfo_Toleration.Size(m)
}
func (m *Toleration) XXX_DiscardUnknown() {
	xxx_messageInfo_Toleration.DiscardUnknown(m)
}

var xxx_messageInfo_Toleration proto.InternalMessageInfo

func (m *Toleration) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Toleration) GetOperator() string {
	if m != nil {
		return m.Operator
	}
	return ""
}

func (m *Toleration) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *Toleration) GetEffect() string {
	if m != nil {
		return m.Effect
	}
	return ""
}

func (m *Toleration) GetTolerationSeconds() int64 {
	if m != nil {
		return m.TolerationSeconds
	}
	return 0
}

func init() {
	proto.RegisterType((*Toleration)(nil), "firmament.Toleration")
}

func init() { proto.RegisterFile("tolerations.proto", fileDescriptor_tolerations_aa2c6496805744bf) }

var fileDescriptor_tolerations_aa2c6496805744bf = []byte{
	// 149 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0xc9, 0xcf, 0x49,
	0x2d, 0x4a, 0x2c, 0xc9, 0xcc, 0xcf, 0x2b, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4c,
	0xcb, 0x2c, 0xca, 0x4d, 0xcc, 0x4d, 0xcd, 0x2b, 0x51, 0x9a, 0xc6, 0xc8, 0xc5, 0x15, 0x02, 0x57,
	0x20, 0x24, 0xc0, 0xc5, 0x9c, 0x9d, 0x5a, 0x29, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0x62,
	0x0a, 0x49, 0x71, 0x71, 0xe4, 0x17, 0x80, 0xa4, 0xf3, 0x8b, 0x24, 0x98, 0xc0, 0xc2, 0x70, 0xbe,
	0x90, 0x08, 0x17, 0x6b, 0x59, 0x62, 0x4e, 0x69, 0xaa, 0x04, 0x33, 0x58, 0x02, 0xc2, 0x11, 0x12,
	0xe3, 0x62, 0x4b, 0x4d, 0x4b, 0x4b, 0x4d, 0x2e, 0x91, 0x60, 0x01, 0x0b, 0x43, 0x79, 0x42, 0x3a,
	0xc8, 0x4e, 0x09, 0x4e, 0x4d, 0xce, 0xcf, 0x4b, 0x29, 0x96, 0x60, 0x55, 0x60, 0xd4, 0x60, 0x0e,
	0xc2, 0x94, 0x48, 0x62, 0x03, 0x3b, 0xd5, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x7b, 0xda, 0xcd,
	0x97, 0xbf, 0x00, 0x00, 0x00,
}