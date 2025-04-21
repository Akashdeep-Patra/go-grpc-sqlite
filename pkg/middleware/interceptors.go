package middleware

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Akashdeep-Patra/go-grpc-sqlite/pkg/logger"
)

// LoggingInterceptor returns a gRPC unary server interceptor for logging requests
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		
		// Extract request metadata
		md, _ := metadata.FromIncomingContext(ctx)
		requestID := extractRequestID(md)
		
		// Create a logger for this request
		logFields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("request_id", requestID),
			zap.Any("request", req),
		}
		
		logger.Info("Received gRPC request", logFields...)
		
		// Call the handler
		resp, err := handler(ctx, req)
		
		// Log the response
		duration := time.Since(startTime)
		responseFields := append(logFields,
			zap.Duration("duration", duration),
			zap.Int64("duration_ms", duration.Milliseconds()),
		)
		
		if err != nil {
			st, _ := status.FromError(err)
			responseFields = append(responseFields,
				zap.String("error", err.Error()),
				zap.String("error_code", st.Code().String()),
			)
			logger.Error("gRPC request error", responseFields...)
		} else {
			responseFields = append(responseFields, zap.Any("response", resp))
			logger.Info("gRPC request completed", responseFields...)
		}
		
		return resp, err
	}
}

// LoggingStreamInterceptor returns a gRPC stream server interceptor for logging requests
func LoggingStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		
		// Extract request metadata
		ctx := ss.Context()
		md, _ := metadata.FromIncomingContext(ctx)
		requestID := extractRequestID(md)
		
		// Create a logger for this request
		logFields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("request_id", requestID),
			zap.Bool("is_client_stream", info.IsClientStream),
			zap.Bool("is_server_stream", info.IsServerStream),
		}
		
		logger.Info("Received gRPC stream request", logFields...)
		
		// Call the handler
		err := handler(srv, ss)
		
		// Log the response
		duration := time.Since(startTime)
		responseFields := append(logFields,
			zap.Duration("duration", duration),
			zap.Int64("duration_ms", duration.Milliseconds()),
		)
		
		if err != nil {
			st, _ := status.FromError(err)
			responseFields = append(responseFields,
				zap.String("error", err.Error()),
				zap.String("error_code", st.Code().String()),
			)
			logger.Error("gRPC stream request error", responseFields...)
		} else {
			logger.Info("gRPC stream request completed", responseFields...)
		}
		
		return err
	}
}

// AuthInterceptor returns a gRPC unary server interceptor for authentication
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip authentication for specific methods
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}
		
		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		
		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}
		
		token := authHeader[0]
		
		// Validate token (implement your own validation logic)
		if !validateToken(token) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		
		// Call the handler
		return handler(ctx, req)
	}
}

// AuthStreamInterceptor returns a gRPC stream server interceptor for authentication
func AuthStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Skip authentication for specific methods
		if isPublicMethod(info.FullMethod) {
			return handler(srv, ss)
		}
		
		// Extract token from metadata
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Error(codes.Unauthenticated, "missing metadata")
		}
		
		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			return status.Error(codes.Unauthenticated, "missing authorization header")
		}
		
		token := authHeader[0]
		
		// Validate token (implement your own validation logic)
		if !validateToken(token) {
			return status.Error(codes.Unauthenticated, "invalid token")
		}
		
		// Call the handler
		return handler(srv, ss)
	}
}

// RecoveryInterceptor returns a gRPC unary server interceptor for panic recovery
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Recovered from panic",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
				)
				err = status.Error(codes.Internal, "Internal server error")
			}
		}()
		
		return handler(ctx, req)
	}
}

// RecoveryStreamInterceptor returns a gRPC stream server interceptor for panic recovery
func RecoveryStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Recovered from stream panic",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
				)
				err = status.Error(codes.Internal, "Internal server error")
			}
		}()
		
		return handler(srv, ss)
	}
}

// RateLimitInterceptor returns a gRPC unary server interceptor for rate limiting
func RateLimitInterceptor() grpc.UnaryServerInterceptor {
	// In a production system, you would use a proper rate limiter like token bucket
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Simple rate limiting logic (this is just a placeholder)
		if isRateLimited(info.FullMethod) {
			return nil, status.Error(codes.ResourceExhausted, "Rate limit exceeded")
		}
		
		return handler(ctx, req)
	}
}

// RateLimitStreamInterceptor returns a gRPC stream server interceptor for rate limiting
func RateLimitStreamInterceptor() grpc.StreamServerInterceptor {
	// In a production system, you would use a proper rate limiter like token bucket
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Simple rate limiting logic (this is just a placeholder)
		if isRateLimited(info.FullMethod) {
			return status.Error(codes.ResourceExhausted, "Rate limit exceeded")
		}
		
		return handler(srv, ss)
	}
}

// Helper functions

func extractRequestID(md metadata.MD) string {
	if values := md.Get("x-request-id"); len(values) > 0 {
		return values[0]
	}
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func isPublicMethod(method string) bool {
	// Define which methods are public and don't require authentication
	publicMethods := map[string]bool{
		"/user.UserService/GetUser": false,
	}
	
	isPublic, exists := publicMethods[method]
	return exists && isPublic
}

func validateToken(token string) bool {
	// Implement your token validation logic here
	// This is just a placeholder
	return len(token) > 0
}

func isRateLimited(method string) bool {
	// Implement your rate limiting logic here
	// This is just a placeholder
	return false
} 