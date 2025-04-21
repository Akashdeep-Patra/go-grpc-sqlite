package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "test-project-grpc/test-project-grpc/user"
	"test-project-grpc/pkg/config"
	"test-project-grpc/pkg/logger"
)

//go:embed swagger-ui
var swaggerUI embed.FS

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
	httpPort           = flag.Int("http-port", 8080, "HTTP server port")
)

func main() {
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(cfg.App.Environment)
	defer logger.Sync()

	logger.Info("Starting gateway server",
		zap.String("grpc_endpoint", *grpcServerEndpoint),
		zap.Int("http_port", *httpPort),
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create gRPC connection to the server
	conn, err := grpc.DialContext(
		ctx,
		*grpcServerEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		logger.Fatal("Failed to dial gRPC server", zap.Error(err))
	}
	defer conn.Close()

	// Create a new ServeMux for the HTTP server
	gwmux := runtime.NewServeMux()

	// Register gRPC service handlers
	err = pb.RegisterUserServiceHandler(ctx, gwmux, conn)
	if err != nil {
		logger.Fatal("Failed to register gateway handler", zap.Error(err))
	}

	// Set up HTTP server
	mux := http.NewServeMux()

	// Register gRPC-Gateway handlers
	mux.Handle("/v1/", gwmux)

	// Register Swagger UI handlers
	mux.HandleFunc("/swagger/", serveSwaggerUI)
	mux.HandleFunc("/swagger.json", serveSwaggerJSON)

	// Start HTTP server
	addr := fmt.Sprintf(":%d", *httpPort)
	logger.Info("HTTP server listening", zap.String("address", addr))
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Fatal("Failed to start HTTP server", zap.Error(err))
	}
}

// serveSwaggerUI serves the Swagger UI files
func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	// Strip the /swagger/ prefix
	uiPath := strings.TrimPrefix(r.URL.Path, "/swagger/")
	if uiPath == "" {
		uiPath = "index.html"
	}

	// Set content type based on file extension
	ext := path.Ext(uiPath)
	if ext != "" {
		contentType := mime.TypeByExtension(ext)
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
	}

	// Open and serve the file from the embedded filesystem
	content, err := fs.ReadFile(swaggerUI, "swagger-ui/"+uiPath)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	_, err = w.Write(content)
	if err != nil {
		logger.Error("Failed to write response", zap.Error(err))
	}
}

// serveSwaggerJSON serves the OpenAPI/Swagger JSON specification
func serveSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Read and serve the swagger JSON file
	json, err := os.ReadFile("api/swagger/api/user.swagger.json")
	if err != nil {
		http.Error(w, "Failed to read Swagger JSON", http.StatusInternalServerError)
		logger.Error("Failed to read Swagger JSON", zap.Error(err))
		return
	}
	
	_, err = w.Write(json)
	if err != nil {
		logger.Error("Failed to write response", zap.Error(err))
	}
} 