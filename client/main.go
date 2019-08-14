package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thuonghidien/grpc-init/proto"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var (
	addr     = flag.String("addr", "127.0.0.1:50051", "Address of grpc server.")
	key      = flag.String("api-key", "", "API key.")
	token    = flag.String("token", "", "Authentication token.")
	keyfile  = flag.String("keyfile", "", "Path to a Google service account key file.")
	audience = flag.String("audience", "", "Audience.")
)

func main() {

	// non HTTPs
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	if *keyfile != "" {
		log.Printf("Authenticating using Google service account key in %s", *keyfile)
		keyBytes, err := ioutil.ReadFile(*keyfile)
		if err != nil {
			log.Fatalf("Unable to read service account key file %s: %v", *keyfile, err)
		}

		tokenSource, err := google.JWTAccessTokenSourceFromJSON(keyBytes, *audience)
		if err != nil {
			log.Fatalf("Error building JWT access token source: %v", err)
		}
		jwt, err := tokenSource.Token()
		if err != nil {
			log.Fatalf("Unable to generate JWT token: %v", err)
		}
		*token = jwt.AccessToken
		// NOTE: the generated JWT token has a 1h TTL.
		// Make sure to refresh the token before it expires by calling TokenSource.Token() for each outgoing requests.
		// Calls to this particular implementation of TokenSource.Token() are cheap.
	}

	g := gin.Default()
	g.GET("/add/:a/:b", func(ctx *gin.Context) {

		a, err := strconv.ParseInt(ctx.Param("a"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}

		b, err := strconv.ParseInt(ctx.Param("b"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}

		req := &proto.Request{A: int64(a), B: int64(b)}
		response, err := client.Add(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	})

	g.GET("/subtract/:a/:b", func(ctx *gin.Context) {

		a, err := strconv.ParseInt(ctx.Param("a"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}

		b, err := strconv.ParseInt(ctx.Param("b"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}

		req := &proto.Request{A: int64(a), B: int64(b)}
		response, err := client.Subtract(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	})

	g.GET("/multiply/:a/:b", func(ctx *gin.Context) {

		a, err := strconv.ParseInt(ctx.Param("a"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}

		b, err := strconv.ParseInt(ctx.Param("b"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}

		req := &proto.Request{A: int64(a), B: int64(b)}
		response, err := client.Multiply(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	})

	g.GET("/divide/:a/:b", func(ctx *gin.Context) {

		a, err := strconv.ParseInt(ctx.Param("a"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}

		b, err := strconv.ParseInt(ctx.Param("b"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}

		req := &proto.Request{A: int64(a), B: int64(b)}
		response, err := client.Divide(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	})

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}