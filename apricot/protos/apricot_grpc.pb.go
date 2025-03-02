// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package apricotpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ApricotClient is the client API for Apricot service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ApricotClient interface {
	NewRunNumber(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*RunNumberResponse, error)
	GetDefaults(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StringMap, error)
	GetVars(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StringMap, error)
	RawGetRecursive(ctx context.Context, in *RawGetRecursiveRequest, opts ...grpc.CallOption) (*ComponentResponse, error)
	// Detectors and host inventories
	ListDetectors(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*DetectorsResponse, error)
	GetHostInventory(ctx context.Context, in *HostGetRequest, opts ...grpc.CallOption) (*HostEntriesResponse, error)
	GetDetectorForHost(ctx context.Context, in *HostRequest, opts ...grpc.CallOption) (*DetectorResponse, error)
	GetDetectorsForHosts(ctx context.Context, in *HostsRequest, opts ...grpc.CallOption) (*DetectorsResponse, error)
	GetCRUCardsForHost(ctx context.Context, in *HostRequest, opts ...grpc.CallOption) (*CRUCardsResponse, error)
	GetEndpointsForCRUCard(ctx context.Context, in *CardRequest, opts ...grpc.CallOption) (*CRUCardEndpointResponse, error)
	// Runtime KV calls
	GetRuntimeEntry(ctx context.Context, in *GetRuntimeEntryRequest, opts ...grpc.CallOption) (*ComponentResponse, error)
	SetRuntimeEntry(ctx context.Context, in *SetRuntimeEntryRequest, opts ...grpc.CallOption) (*Empty, error)
	ListRuntimeEntries(ctx context.Context, in *ListRuntimeEntriesRequest, opts ...grpc.CallOption) (*ComponentEntriesResponse, error)
	// Component configuration calls
	ListComponents(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ComponentEntriesResponse, error)
	ListComponentEntries(ctx context.Context, in *ListComponentEntriesRequest, opts ...grpc.CallOption) (*ComponentEntriesResponse, error)
	ListComponentEntryHistory(ctx context.Context, in *ComponentQuery, opts ...grpc.CallOption) (*ComponentEntriesResponse, error)
	GetComponentConfiguration(ctx context.Context, in *ComponentRequest, opts ...grpc.CallOption) (*ComponentResponse, error)
	GetComponentConfigurationWithLastIndex(ctx context.Context, in *ComponentRequest, opts ...grpc.CallOption) (*ComponentResponseWithLastIndex, error)
	ImportComponentConfiguration(ctx context.Context, in *ImportComponentConfigurationRequest, opts ...grpc.CallOption) (*ImportComponentConfigurationResponse, error)
}

type apricotClient struct {
	cc grpc.ClientConnInterface
}

func NewApricotClient(cc grpc.ClientConnInterface) ApricotClient {
	return &apricotClient{cc}
}

func (c *apricotClient) NewRunNumber(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*RunNumberResponse, error) {
	out := new(RunNumberResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/NewRunNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) GetDefaults(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StringMap, error) {
	out := new(StringMap)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/GetDefaults", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) GetVars(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StringMap, error) {
	out := new(StringMap)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/GetVars", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) RawGetRecursive(ctx context.Context, in *RawGetRecursiveRequest, opts ...grpc.CallOption) (*ComponentResponse, error) {
	out := new(ComponentResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/RawGetRecursive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) ListDetectors(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*DetectorsResponse, error) {
	out := new(DetectorsResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/ListDetectors", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) GetHostInventory(ctx context.Context, in *HostGetRequest, opts ...grpc.CallOption) (*HostEntriesResponse, error) {
	out := new(HostEntriesResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/GetHostInventory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) GetDetectorForHost(ctx context.Context, in *HostRequest, opts ...grpc.CallOption) (*DetectorResponse, error) {
	out := new(DetectorResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/GetDetectorForHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) GetDetectorsForHosts(ctx context.Context, in *HostsRequest, opts ...grpc.CallOption) (*DetectorsResponse, error) {
	out := new(DetectorsResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/GetDetectorsForHosts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) GetCRUCardsForHost(ctx context.Context, in *HostRequest, opts ...grpc.CallOption) (*CRUCardsResponse, error) {
	out := new(CRUCardsResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/GetCRUCardsForHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) GetEndpointsForCRUCard(ctx context.Context, in *CardRequest, opts ...grpc.CallOption) (*CRUCardEndpointResponse, error) {
	out := new(CRUCardEndpointResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/GetEndpointsForCRUCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) GetRuntimeEntry(ctx context.Context, in *GetRuntimeEntryRequest, opts ...grpc.CallOption) (*ComponentResponse, error) {
	out := new(ComponentResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/GetRuntimeEntry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) SetRuntimeEntry(ctx context.Context, in *SetRuntimeEntryRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/SetRuntimeEntry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) ListRuntimeEntries(ctx context.Context, in *ListRuntimeEntriesRequest, opts ...grpc.CallOption) (*ComponentEntriesResponse, error) {
	out := new(ComponentEntriesResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/ListRuntimeEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) ListComponents(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ComponentEntriesResponse, error) {
	out := new(ComponentEntriesResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/ListComponents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) ListComponentEntries(ctx context.Context, in *ListComponentEntriesRequest, opts ...grpc.CallOption) (*ComponentEntriesResponse, error) {
	out := new(ComponentEntriesResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/ListComponentEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) ListComponentEntryHistory(ctx context.Context, in *ComponentQuery, opts ...grpc.CallOption) (*ComponentEntriesResponse, error) {
	out := new(ComponentEntriesResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/ListComponentEntryHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) GetComponentConfiguration(ctx context.Context, in *ComponentRequest, opts ...grpc.CallOption) (*ComponentResponse, error) {
	out := new(ComponentResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/GetComponentConfiguration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) GetComponentConfigurationWithLastIndex(ctx context.Context, in *ComponentRequest, opts ...grpc.CallOption) (*ComponentResponseWithLastIndex, error) {
	out := new(ComponentResponseWithLastIndex)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/GetComponentConfigurationWithLastIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apricotClient) ImportComponentConfiguration(ctx context.Context, in *ImportComponentConfigurationRequest, opts ...grpc.CallOption) (*ImportComponentConfigurationResponse, error) {
	out := new(ImportComponentConfigurationResponse)
	err := c.cc.Invoke(ctx, "/apricot.Apricot/ImportComponentConfiguration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApricotServer is the server API for Apricot service.
// All implementations should embed UnimplementedApricotServer
// for forward compatibility
type ApricotServer interface {
	NewRunNumber(context.Context, *Empty) (*RunNumberResponse, error)
	GetDefaults(context.Context, *Empty) (*StringMap, error)
	GetVars(context.Context, *Empty) (*StringMap, error)
	RawGetRecursive(context.Context, *RawGetRecursiveRequest) (*ComponentResponse, error)
	// Detectors and host inventories
	ListDetectors(context.Context, *Empty) (*DetectorsResponse, error)
	GetHostInventory(context.Context, *HostGetRequest) (*HostEntriesResponse, error)
	GetDetectorForHost(context.Context, *HostRequest) (*DetectorResponse, error)
	GetDetectorsForHosts(context.Context, *HostsRequest) (*DetectorsResponse, error)
	GetCRUCardsForHost(context.Context, *HostRequest) (*CRUCardsResponse, error)
	GetEndpointsForCRUCard(context.Context, *CardRequest) (*CRUCardEndpointResponse, error)
	// Runtime KV calls
	GetRuntimeEntry(context.Context, *GetRuntimeEntryRequest) (*ComponentResponse, error)
	SetRuntimeEntry(context.Context, *SetRuntimeEntryRequest) (*Empty, error)
	ListRuntimeEntries(context.Context, *ListRuntimeEntriesRequest) (*ComponentEntriesResponse, error)
	// Component configuration calls
	ListComponents(context.Context, *Empty) (*ComponentEntriesResponse, error)
	ListComponentEntries(context.Context, *ListComponentEntriesRequest) (*ComponentEntriesResponse, error)
	ListComponentEntryHistory(context.Context, *ComponentQuery) (*ComponentEntriesResponse, error)
	GetComponentConfiguration(context.Context, *ComponentRequest) (*ComponentResponse, error)
	GetComponentConfigurationWithLastIndex(context.Context, *ComponentRequest) (*ComponentResponseWithLastIndex, error)
	ImportComponentConfiguration(context.Context, *ImportComponentConfigurationRequest) (*ImportComponentConfigurationResponse, error)
}

// UnimplementedApricotServer should be embedded to have forward compatible implementations.
type UnimplementedApricotServer struct {
}

func (UnimplementedApricotServer) NewRunNumber(context.Context, *Empty) (*RunNumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewRunNumber not implemented")
}
func (UnimplementedApricotServer) GetDefaults(context.Context, *Empty) (*StringMap, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefaults not implemented")
}
func (UnimplementedApricotServer) GetVars(context.Context, *Empty) (*StringMap, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVars not implemented")
}
func (UnimplementedApricotServer) RawGetRecursive(context.Context, *RawGetRecursiveRequest) (*ComponentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawGetRecursive not implemented")
}
func (UnimplementedApricotServer) ListDetectors(context.Context, *Empty) (*DetectorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDetectors not implemented")
}
func (UnimplementedApricotServer) GetHostInventory(context.Context, *HostGetRequest) (*HostEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHostInventory not implemented")
}
func (UnimplementedApricotServer) GetDetectorForHost(context.Context, *HostRequest) (*DetectorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDetectorForHost not implemented")
}
func (UnimplementedApricotServer) GetDetectorsForHosts(context.Context, *HostsRequest) (*DetectorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDetectorsForHosts not implemented")
}
func (UnimplementedApricotServer) GetCRUCardsForHost(context.Context, *HostRequest) (*CRUCardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCRUCardsForHost not implemented")
}
func (UnimplementedApricotServer) GetEndpointsForCRUCard(context.Context, *CardRequest) (*CRUCardEndpointResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEndpointsForCRUCard not implemented")
}
func (UnimplementedApricotServer) GetRuntimeEntry(context.Context, *GetRuntimeEntryRequest) (*ComponentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRuntimeEntry not implemented")
}
func (UnimplementedApricotServer) SetRuntimeEntry(context.Context, *SetRuntimeEntryRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetRuntimeEntry not implemented")
}
func (UnimplementedApricotServer) ListRuntimeEntries(context.Context, *ListRuntimeEntriesRequest) (*ComponentEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRuntimeEntries not implemented")
}
func (UnimplementedApricotServer) ListComponents(context.Context, *Empty) (*ComponentEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListComponents not implemented")
}
func (UnimplementedApricotServer) ListComponentEntries(context.Context, *ListComponentEntriesRequest) (*ComponentEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListComponentEntries not implemented")
}
func (UnimplementedApricotServer) ListComponentEntryHistory(context.Context, *ComponentQuery) (*ComponentEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListComponentEntryHistory not implemented")
}
func (UnimplementedApricotServer) GetComponentConfiguration(context.Context, *ComponentRequest) (*ComponentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetComponentConfiguration not implemented")
}
func (UnimplementedApricotServer) GetComponentConfigurationWithLastIndex(context.Context, *ComponentRequest) (*ComponentResponseWithLastIndex, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetComponentConfigurationWithLastIndex not implemented")
}
func (UnimplementedApricotServer) ImportComponentConfiguration(context.Context, *ImportComponentConfigurationRequest) (*ImportComponentConfigurationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportComponentConfiguration not implemented")
}

// UnsafeApricotServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ApricotServer will
// result in compilation errors.
type UnsafeApricotServer interface {
	mustEmbedUnimplementedApricotServer()
}

func RegisterApricotServer(s grpc.ServiceRegistrar, srv ApricotServer) {
	s.RegisterService(&Apricot_ServiceDesc, srv)
}

func _Apricot_NewRunNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).NewRunNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/NewRunNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).NewRunNumber(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_GetDefaults_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).GetDefaults(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/GetDefaults",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).GetDefaults(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_GetVars_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).GetVars(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/GetVars",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).GetVars(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_RawGetRecursive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawGetRecursiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).RawGetRecursive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/RawGetRecursive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).RawGetRecursive(ctx, req.(*RawGetRecursiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_ListDetectors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).ListDetectors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/ListDetectors",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).ListDetectors(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_GetHostInventory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).GetHostInventory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/GetHostInventory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).GetHostInventory(ctx, req.(*HostGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_GetDetectorForHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).GetDetectorForHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/GetDetectorForHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).GetDetectorForHost(ctx, req.(*HostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_GetDetectorsForHosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).GetDetectorsForHosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/GetDetectorsForHosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).GetDetectorsForHosts(ctx, req.(*HostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_GetCRUCardsForHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).GetCRUCardsForHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/GetCRUCardsForHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).GetCRUCardsForHost(ctx, req.(*HostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_GetEndpointsForCRUCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).GetEndpointsForCRUCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/GetEndpointsForCRUCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).GetEndpointsForCRUCard(ctx, req.(*CardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_GetRuntimeEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRuntimeEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).GetRuntimeEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/GetRuntimeEntry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).GetRuntimeEntry(ctx, req.(*GetRuntimeEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_SetRuntimeEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRuntimeEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).SetRuntimeEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/SetRuntimeEntry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).SetRuntimeEntry(ctx, req.(*SetRuntimeEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_ListRuntimeEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRuntimeEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).ListRuntimeEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/ListRuntimeEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).ListRuntimeEntries(ctx, req.(*ListRuntimeEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_ListComponents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).ListComponents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/ListComponents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).ListComponents(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_ListComponentEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListComponentEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).ListComponentEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/ListComponentEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).ListComponentEntries(ctx, req.(*ListComponentEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_ListComponentEntryHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComponentQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).ListComponentEntryHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/ListComponentEntryHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).ListComponentEntryHistory(ctx, req.(*ComponentQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_GetComponentConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComponentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).GetComponentConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/GetComponentConfiguration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).GetComponentConfiguration(ctx, req.(*ComponentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_GetComponentConfigurationWithLastIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComponentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).GetComponentConfigurationWithLastIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/GetComponentConfigurationWithLastIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).GetComponentConfigurationWithLastIndex(ctx, req.(*ComponentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Apricot_ImportComponentConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportComponentConfigurationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApricotServer).ImportComponentConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apricot.Apricot/ImportComponentConfiguration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApricotServer).ImportComponentConfiguration(ctx, req.(*ImportComponentConfigurationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Apricot_ServiceDesc is the grpc.ServiceDesc for Apricot service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Apricot_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "apricot.Apricot",
	HandlerType: (*ApricotServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewRunNumber",
			Handler:    _Apricot_NewRunNumber_Handler,
		},
		{
			MethodName: "GetDefaults",
			Handler:    _Apricot_GetDefaults_Handler,
		},
		{
			MethodName: "GetVars",
			Handler:    _Apricot_GetVars_Handler,
		},
		{
			MethodName: "RawGetRecursive",
			Handler:    _Apricot_RawGetRecursive_Handler,
		},
		{
			MethodName: "ListDetectors",
			Handler:    _Apricot_ListDetectors_Handler,
		},
		{
			MethodName: "GetHostInventory",
			Handler:    _Apricot_GetHostInventory_Handler,
		},
		{
			MethodName: "GetDetectorForHost",
			Handler:    _Apricot_GetDetectorForHost_Handler,
		},
		{
			MethodName: "GetDetectorsForHosts",
			Handler:    _Apricot_GetDetectorsForHosts_Handler,
		},
		{
			MethodName: "GetCRUCardsForHost",
			Handler:    _Apricot_GetCRUCardsForHost_Handler,
		},
		{
			MethodName: "GetEndpointsForCRUCard",
			Handler:    _Apricot_GetEndpointsForCRUCard_Handler,
		},
		{
			MethodName: "GetRuntimeEntry",
			Handler:    _Apricot_GetRuntimeEntry_Handler,
		},
		{
			MethodName: "SetRuntimeEntry",
			Handler:    _Apricot_SetRuntimeEntry_Handler,
		},
		{
			MethodName: "ListRuntimeEntries",
			Handler:    _Apricot_ListRuntimeEntries_Handler,
		},
		{
			MethodName: "ListComponents",
			Handler:    _Apricot_ListComponents_Handler,
		},
		{
			MethodName: "ListComponentEntries",
			Handler:    _Apricot_ListComponentEntries_Handler,
		},
		{
			MethodName: "ListComponentEntryHistory",
			Handler:    _Apricot_ListComponentEntryHistory_Handler,
		},
		{
			MethodName: "GetComponentConfiguration",
			Handler:    _Apricot_GetComponentConfiguration_Handler,
		},
		{
			MethodName: "GetComponentConfigurationWithLastIndex",
			Handler:    _Apricot_GetComponentConfigurationWithLastIndex_Handler,
		},
		{
			MethodName: "ImportComponentConfiguration",
			Handler:    _Apricot_ImportComponentConfiguration_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/apricot.proto",
}
