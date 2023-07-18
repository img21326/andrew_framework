package andrewframework

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/img21326/andrew_framework/helper"
	"github.com/img21326/andrew_framework/middleware"
)

func InitServer() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.WithLoggerMiddleware())
	r.Use(middleware.WithGormMiddleware())
	r.Use(middleware.ReturnErrorMiddleware())

	for _, router := range RouterList {
		router.AddRoute(r)
	}
	return r
}

func Start() {
	srv := &http.Server{
		Addr:    ":8000",
		Handler: InitServer(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer helper.WaitForLoggerComplete()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

func StartAWSLambda() {
	var ginLambda *ginAdapter.GinLambda

	handler := func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return ginLambda.ProxyWithContext(ctx, req)
	}

	r := InitServer()
	ginLambda = ginAdapter.New(r)
	lambda.Start(handler)
}
