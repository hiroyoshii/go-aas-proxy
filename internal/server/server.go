package server

import (
	"context"
	"fmt"
	basyxpb "hiroyoshii/go-aas-proxy/gen/proto"
	"log"
	"net"

	"github.com/caarlos0/env"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type config struct {
	Port string `env:"GRPC_PORT" envDefault:":50051"`
}

type server struct {
	cfg        *config
	grpcServer *grpc.Server
	listener   net.Listener
	isAlive    bool
	basyxpb.UnimplementedBasyxOpenapiServer
	basyxpb.UnsafeBasyxOpenapiServer
}

// Server is the interface for API Server
type Server interface {
	GetAssetAdministrationShell(context.Context, *emptypb.Empty) (*basyxpb.AssetAdministrationShellDescriptor, error)
	GetSubmodelsFromShell(context.Context, *emptypb.Empty) (*basyxpb.Submodel, error)
	GETAasSubmodelsSubmodelIdShort(context.Context, *basyxpb.GETAasSubmodelsSubmodelIdShortParameters) (*basyxpb.Submodel, error)
	PutSubmodelToShell(context.Context, *basyxpb.PutSubmodelToShellParameters) (*basyxpb.Submodel, error)
	DeleteSubmodelFromShellByIdShort(context.Context, *basyxpb.DeleteSubmodelFromShellByIdShortParameters) (*basyxpb.Result, error)
	GetSubmodelFromShellByIdShort(context.Context, *basyxpb.GetSubmodelFromShellByIdShortParameters) (*basyxpb.Submodel, error)
	ShellGetSubmodelValues(context.Context, *basyxpb.ShellGetSubmodelValuesParameters) (*basyxpb.Result, error)
	ShellGetSubmodelElements(context.Context, *basyxpb.ShellGetSubmodelElementsParameters) (*basyxpb.SubmodelElement, error)
	ShellGetSubmodelElementByIdShort(context.Context, *basyxpb.ShellGetSubmodelElementByIdShortParameters) (*basyxpb.SubmodelElement, error)
	ShellPutSubmodelElement(context.Context, *basyxpb.ShellPutSubmodelElementParameters) (*basyxpb.SubmodelElement, error)
	ShellDeleteSubmodelElementByIdShort(context.Context, *basyxpb.ShellDeleteSubmodelElementByIdShortParameters) (*basyxpb.Result, error)
	ShellGetSubmodelElementValueByIdShort(context.Context, *basyxpb.ShellGetSubmodelElementValueByIdShortParameters) (*basyxpb.ShellGetSubmodelElementValueByIdShortOK, error)
	ShellPutSubmodelElementValueByIdShort(context.Context, *basyxpb.ShellPutSubmodelElementValueByIdShortParameters) (*basyxpb.ElementValue, error)
	ShellInvokeOperationByIdShort(context.Context, *basyxpb.ShellInvokeOperationByIdShortParameters) (*basyxpb.Result, error)
	ShellGetInvocationResultByIdShort(context.Context, *basyxpb.ShellGetInvocationResultByIdShortParameters) (*basyxpb.InvocationResponse, error)
	Serve() error
	GracefulStop()
	IsAlive() bool
}

// healthServer implements health.HealthServer
type healthServer struct {
}

// Check is a function to return health status
func (h *healthServer) Check(context.Context, *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is not used but needs to be implemented
func (h *healthServer) Watch(*health.HealthCheckRequest, health.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "watch is not implemented.")
}

// NewServer get the new Server struct
func NewServer(base context.Context) (Server, error) {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	lis, err := net.Listen("tcp", cfg.Port)
	defer log.Printf("Server is listening on port %s\n", cfg.Port)
	if err != nil {
		return nil, err
	}

	// Setup metrics.
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)

	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(recoveryFunc)),
			otelgrpc.UnaryServerInterceptor(),
		),
	)

	gs := &server{
		cfg:        cfg,
		grpcServer: s,
		listener:   lis,
		isAlive:    false,
	}
	hs := &healthServer{}
	health.RegisterHealthServer(s, hs)
	basyxpb.RegisterBasyxOpenapiServer(s, gs)

	reflection.Register(s)
	return gs, nil
}

func recoveryFunc(p interface{}) error {
	log.Printf("unexpected panic occured: %v\n", p)
	return status.Errorf(codes.Internal, "Unexpected error")
}

func (s *server) GetAssetAdministrationShell(context.Context, *emptypb.Empty) (*basyxpb.AssetAdministrationShellDescriptor, error) {
	return nil, status.Error(codes.Unimplemented, "method GetAssetAdministrationShell not implemented")
}
func (s *server) GetSubmodelsFromShell(context.Context, *emptypb.Empty) (*basyxpb.Submodel, error) {
	return nil, status.Error(codes.Unimplemented, "method GetSubmodelsFromShell not implemented")
}
func (s *server) GETAasSubmodelsSubmodelIdShort(context.Context, *basyxpb.GETAasSubmodelsSubmodelIdShortParameters) (*basyxpb.Submodel, error) {
	return nil, status.Error(codes.Unimplemented, "method GETAasSubmodelsSubmodelIdShort not implemented")
}
func (s *server) PutSubmodelToShell(context.Context, *basyxpb.PutSubmodelToShellParameters) (*basyxpb.Submodel, error) {
	return nil, status.Error(codes.Unimplemented, "method PutSubmodelToShell not implemented")
}
func (s *server) DeleteSubmodelFromShellByIdShort(context.Context, *basyxpb.DeleteSubmodelFromShellByIdShortParameters) (*basyxpb.Result, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteSubmodelFromShellByIdShort not implemented")
}
func (s *server) GetSubmodelFromShellByIdShort(context.Context, *basyxpb.GetSubmodelFromShellByIdShortParameters) (*basyxpb.Submodel, error) {
	return nil, status.Error(codes.Unimplemented, "method GetSubmodelFromShellByIdShort not implemented")
}
func (s *server) ShellGetSubmodelValues(context.Context, *basyxpb.ShellGetSubmodelValuesParameters) (*basyxpb.Result, error) {
	return nil, status.Error(codes.Unimplemented, "method ShellGetSubmodelValues not implemented")
}
func (s *server) ShellGetSubmodelElements(context.Context, *basyxpb.ShellGetSubmodelElementsParameters) (*basyxpb.SubmodelElement, error) {
	return nil, status.Error(codes.Unimplemented, "method ShellGetSubmodelElements not implemented")
}
func (s *server) ShellGetSubmodelElementByIdShort(context.Context, *basyxpb.ShellGetSubmodelElementByIdShortParameters) (*basyxpb.SubmodelElement, error) {
	return nil, status.Error(codes.Unimplemented, "method ShellGetSubmodelElementByIdShort not implemented")
}
func (s *server) ShellPutSubmodelElement(context.Context, *basyxpb.ShellPutSubmodelElementParameters) (*basyxpb.SubmodelElement, error) {
	return nil, status.Error(codes.Unimplemented, "method ShellPutSubmodelElement not implemented")
}
func (s *server) ShellDeleteSubmodelElementByIdShort(context.Context, *basyxpb.ShellDeleteSubmodelElementByIdShortParameters) (*basyxpb.Result, error) {
	return nil, status.Error(codes.Unimplemented, "method ShellDeleteSubmodelElementByIdShort not implemented")
}
func (s *server) ShellGetSubmodelElementValueByIdShort(context.Context, *basyxpb.ShellGetSubmodelElementValueByIdShortParameters) (*basyxpb.ShellGetSubmodelElementValueByIdShortOK, error) {
	return nil, status.Error(codes.Unimplemented, "method ShellGetSubmodelElementValueByIdShort not implemented")
}
func (s *server) ShellPutSubmodelElementValueByIdShort(context.Context, *basyxpb.ShellPutSubmodelElementValueByIdShortParameters) (*basyxpb.ElementValue, error) {
	return nil, status.Error(codes.Unimplemented, "method ShellPutSubmodelElementValueByIdShort not implemented")
}
func (s *server) ShellInvokeOperationByIdShort(context.Context, *basyxpb.ShellInvokeOperationByIdShortParameters) (*basyxpb.Result, error) {
	return nil, status.Error(codes.Unimplemented, "method ShellInvokeOperationByIdShort not implemented")
}
func (s *server) ShellGetInvocationResultByIdShort(context.Context, *basyxpb.ShellGetInvocationResultByIdShortParameters) (*basyxpb.InvocationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ShellGetInvocationResultByIdShort not implemented")
}

// Serve starts API Server
func (s *server) Serve() error {
	if s.isAlive {
		return fmt.Errorf("server is already running")
	}

	go func() {
		if err := s.grpcServer.Serve(s.listener); err != nil {
			log.Printf("failed to serve grpcServer, detail %v\n", err)
		}
	}()
	s.isAlive = true

	return nil
}

// GracefulStop stops API Server
func (s *server) GracefulStop() {
	log.Println("GracefulStop started")
	defer log.Println("GracefulStop finished")

	s.grpcServer.GracefulStop()
	s.isAlive = false
}

func (s *server) IsAlive() bool {
	return s.isAlive
}
