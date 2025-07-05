package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	// "github.com/kedacore/keda/v2/pkg/scalers/externalscaler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	externalscaler "github.com/jxs1211/keda-external-scaler/pkg/api"
	fn "github.com/jxs1211/keda-external-scaler/pkg/path/to/func"
	"github.com/jxs1211/keda-external-scaler/pkg/util/debug"
)

var (
	dlog = debug.NewLogger()
)

type StockPriceScaler struct {
	// externalscaler.ExternalScalerServer
	externalscaler.ExternalScalerServer
	targetPrice float64
	stockSymbol string
}

func (s *StockPriceScaler) IsActive(ctx context.Context, scaledObject *externalscaler.ScaledObjectRef) (*externalscaler.IsActiveResponse, error) {
	log.Println("isActive with: ", scaledObject)
	currentPrice, err := getCurrentStockPrice(s.stockSymbol)
	if err != nil {
		return &externalscaler.IsActiveResponse{Result: false}, nil
	}

	return &externalscaler.IsActiveResponse{Result: currentPrice > s.targetPrice}, nil
}

func (s *StockPriceScaler) GetMetricSpec(context.Context, *externalscaler.ScaledObjectRef) (*externalscaler.GetMetricSpecResponse, error) {
	return &externalscaler.GetMetricSpecResponse{
		MetricSpecs: []*externalscaler.MetricSpec{{
			MetricName: "stockPrice",
			TargetSize: 1,
		}},
	}, nil
}

func (s *StockPriceScaler) GetMetrics(_ context.Context, metricRequest *externalscaler.GetMetricsRequest) (*externalscaler.GetMetricsResponse, error) {
	currentPrice, err := getCurrentStockPrice(s.stockSymbol)
	if err != nil {
		return nil, err
	}

	return &externalscaler.GetMetricsResponse{
		MetricValues: []*externalscaler.MetricValue{{
			MetricName:  "stockPrice",
			MetricValue: int64(currentPrice * 100),
		}},
	}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	externalscaler.RegisterExternalScalerServer(grpcServer, &StockPriceScaler{})
	reflection.Register(grpcServer)
	// Create a new HTTP server for health check
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }),
	}
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		panic(err)
	}
	// Start the gRPC server in a goroutine
	go func() {
		fmt.Println("StockPrice Scaler listening on port 6000")
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// Start the HTTP server
	fmt.Println("Health check server listening on port 8080")
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

// ... existing helper functions ...

func getCurrentStockPrice(symbol string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", symbol))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	log.Printf("response body: %s, err: %v", string(body), err)
	fn.RunDebugLog()
	// debug.Printf("response body: %s, err: %v", string(body), err)
	// dlog.Printf("response body: %s, err: %v", string(body), err)
	// dlog.Println("response body:", string(body))
	if err != nil {
		return 0, err
	}
	// Add actual JSON parsing logic here
	// For example, using a JSON library like "encoding/json"
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}
	// Extract and return the current price
	// price, err := strconv.ParseFloat(result["chart"]["result"][0]["indicators"]["quote"][0]["close"][0].(string), 64)
	// if err != nil {
	// 	return 0, err
	// }
	return 0, nil
}
