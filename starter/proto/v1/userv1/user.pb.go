package userv1

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/descriptorpb"
)

func strPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32 { return &i }

var (
	FileDescriptor protoreflect.FileDescriptor
	reqDesc        protoreflect.MessageDescriptor
	respDesc       protoreflect.MessageDescriptor
)

func init() {
	syntax := "proto3"
	pkg := "user.v1"
	name := "starter/proto/v1/user.proto"

	fd := &descriptorpb.FileDescriptorProto{
		Name:    &name,
		Package: &pkg,
		Syntax:  &syntax,
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: strPtr("GetProfileRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     strPtr("user_id"),
						Number:   int32Ptr(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: strPtr("userId"),
					},
				},
			},
			{
				Name: strPtr("GetProfileResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     strPtr("user_id"),
						Number:   int32Ptr(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: strPtr("userId"),
					},
					{
						Name:     strPtr("name"),
						Number:   int32Ptr(2),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: strPtr("name"),
					},
					{
						Name:     strPtr("email"),
						Number:   int32Ptr(3),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: strPtr("email"),
					},
					{
						Name:     strPtr("role"),
						Number:   int32Ptr(4),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: strPtr("role"),
					},
					{
						Name:     strPtr("tenant_slug"),
						Number:   int32Ptr(5),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: strPtr("tenantSlug"),
					},
				},
			},
		},
	}

	var err error
	FileDescriptor, err = protodesc.NewFile(fd, nil)
	if err != nil {
		panic(err)
	}
	reqDesc = FileDescriptor.Messages().ByName("GetProfileRequest")
	respDesc = FileDescriptor.Messages().ByName("GetProfileResponse")
}

// -----------------------------------------------------------------------------
// GetProfileRequest
// -----------------------------------------------------------------------------

type GetProfileRequest struct {
	UserId string `json:"user_id,omitempty"`
}

func (x *GetProfileRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetProfileRequest) ProtoMessage() {}
func (x *GetProfileRequest) Reset()        { *x = GetProfileRequest{} }
func (x *GetProfileRequest) String() string {
	b, _ := protojson.Marshal(x)
	return string(b)
}

func (x *GetProfileRequest) ProtoReflect() protoreflect.Message {
	return &getProfileRequestReflect{req: x}
}

type getProfileRequestReflect struct {
	req *GetProfileRequest
}

func (r *getProfileRequestReflect) Descriptor() protoreflect.MessageDescriptor { return reqDesc }
func (r *getProfileRequestReflect) Type() protoreflect.MessageType             { return nil }
func (r *getProfileRequestReflect) New() protoreflect.Message {
	return (&GetProfileRequest{}).ProtoReflect()
}
func (r *getProfileRequestReflect) Interface() protoreflect.ProtoMessage { return r.req }
func (r *getProfileRequestReflect) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if r.req.UserId != "" {
		fd := reqDesc.Fields().ByName("user_id")
		f(fd, protoreflect.ValueOfString(r.req.UserId))
	}
}
func (r *getProfileRequestReflect) Has(fd protoreflect.FieldDescriptor) bool {
	return r.req.UserId != ""
}
func (r *getProfileRequestReflect) Clear(fd protoreflect.FieldDescriptor) {
	r.req.UserId = ""
}
func (r *getProfileRequestReflect) Get(fd protoreflect.FieldDescriptor) protoreflect.Value {
	if fd.Name() == "user_id" {
		return protoreflect.ValueOfString(r.req.UserId)
	}
	return protoreflect.ValueOfString("")
}
func (r *getProfileRequestReflect) Set(fd protoreflect.FieldDescriptor, v protoreflect.Value) {
	if fd.Name() == "user_id" {
		r.req.UserId = v.String()
	}
}
func (r *getProfileRequestReflect) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	panic("mutable not supported for scalar string")
}
func (r *getProfileRequestReflect) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	panic("newfield not supported for scalar string")
}
func (r *getProfileRequestReflect) WhichOneof(od protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	return nil
}
func (r *getProfileRequestReflect) GetUnknown() protoreflect.RawFields               { return nil }
func (r *getProfileRequestReflect) SetUnknown(b protoreflect.RawFields)             {}
func (r *getProfileRequestReflect) IsValid() bool                                   { return r.req != nil }
func (r *getProfileRequestReflect) ProtoMethods() *protoiface.Methods               { return nil }

// -----------------------------------------------------------------------------
// GetProfileResponse
// -----------------------------------------------------------------------------

type GetProfileResponse struct {
	UserId     string `json:"user_id,omitempty"`
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
	Role       string `json:"role,omitempty"`
	TenantSlug string `json:"tenant_slug,omitempty"`
}

func (x *GetProfileResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetProfileResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetProfileResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *GetProfileResponse) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *GetProfileResponse) GetTenantSlug() string {
	if x != nil {
		return x.TenantSlug
	}
	return ""
}

func (x *GetProfileResponse) ProtoMessage() {}
func (x *GetProfileResponse) Reset()        { *x = GetProfileResponse{} }
func (x *GetProfileResponse) String() string {
	b, _ := protojson.Marshal(x)
	return string(b)
}

func (x *GetProfileResponse) ProtoReflect() protoreflect.Message {
	return &getProfileResponseReflect{resp: x}
}

type getProfileResponseReflect struct {
	resp *GetProfileResponse
}

func (r *getProfileResponseReflect) Descriptor() protoreflect.MessageDescriptor { return respDesc }
func (r *getProfileResponseReflect) Type() protoreflect.MessageType             { return nil }
func (r *getProfileResponseReflect) New() protoreflect.Message {
	return (&GetProfileResponse{}).ProtoReflect()
}
func (r *getProfileResponseReflect) Interface() protoreflect.ProtoMessage { return r.resp }
func (r *getProfileResponseReflect) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	fields := respDesc.Fields()
	if r.resp.UserId != "" {
		f(fields.ByName("user_id"), protoreflect.ValueOfString(r.resp.UserId))
	}
	if r.resp.Name != "" {
		f(fields.ByName("name"), protoreflect.ValueOfString(r.resp.Name))
	}
	if r.resp.Email != "" {
		f(fields.ByName("email"), protoreflect.ValueOfString(r.resp.Email))
	}
	if r.resp.Role != "" {
		f(fields.ByName("role"), protoreflect.ValueOfString(r.resp.Role))
	}
	if r.resp.TenantSlug != "" {
		f(fields.ByName("tenant_slug"), protoreflect.ValueOfString(r.resp.TenantSlug))
	}
}
func (r *getProfileResponseReflect) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.Name() {
	case "user_id":
		return r.resp.UserId != ""
	case "name":
		return r.resp.Name != ""
	case "email":
		return r.resp.Email != ""
	case "role":
		return r.resp.Role != ""
	case "tenant_slug":
		return r.resp.TenantSlug != ""
	}
	return false
}
func (r *getProfileResponseReflect) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.Name() {
	case "user_id":
		r.resp.UserId = ""
	case "name":
		r.resp.Name = ""
	case "email":
		r.resp.Email = ""
	case "role":
		r.resp.Role = ""
	case "tenant_slug":
		r.resp.TenantSlug = ""
	}
}
func (r *getProfileResponseReflect) Get(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.Name() {
	case "user_id":
		return protoreflect.ValueOfString(r.resp.UserId)
	case "name":
		return protoreflect.ValueOfString(r.resp.Name)
	case "email":
		return protoreflect.ValueOfString(r.resp.Email)
	case "role":
		return protoreflect.ValueOfString(r.resp.Role)
	case "tenant_slug":
		return protoreflect.ValueOfString(r.resp.TenantSlug)
	}
	return protoreflect.ValueOfString("")
}
func (r *getProfileResponseReflect) Set(fd protoreflect.FieldDescriptor, v protoreflect.Value) {
	switch fd.Name() {
	case "user_id":
		r.resp.UserId = v.String()
	case "name":
		r.resp.Name = v.String()
	case "email":
		r.resp.Email = v.String()
	case "role":
		r.resp.Role = v.String()
	case "tenant_slug":
		r.resp.TenantSlug = v.String()
	}
}
func (r *getProfileResponseReflect) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	panic("mutable not supported for scalar string")
}
func (r *getProfileResponseReflect) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	panic("newfield not supported for scalar string")
}
func (r *getProfileResponseReflect) WhichOneof(od protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	return nil
}
func (r *getProfileResponseReflect) GetUnknown() protoreflect.RawFields               { return nil }
func (r *getProfileResponseReflect) SetUnknown(b protoreflect.RawFields)             {}
func (r *getProfileResponseReflect) IsValid() bool                                   { return r.resp != nil }
func (r *getProfileResponseReflect) ProtoMethods() *protoiface.Methods               { return nil }
