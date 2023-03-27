package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"

	"github.com/kienmatu/go-connection-pooling/model"

	"github.com/gin-gonic/gin"
)

var allTime int64 = 0
var allCount int64 = 0

var newTime int64 = 0
var newCount int64 = 0

var poolTime int64 = 0
var poolCount int64 = 0

var dsn = "postgres://postgres:password1@localhost:5433/postgres?sslmode=disable"
var query = "SELECT id, name, price, description FROM products limit 1000"

func scanProducts(rows pgx.Rows) ([]*model.Product, error) {
	defer rows.Close()

	products := make([]*model.Product, 0)
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

func main() {
	// Postgres allows 100 connections in default
	// Set the maximum number of idle connections in the pool
	idleConn := 50
	// Set the maximum number of connections in the pool
	maxConnections := 90
	// Set the maximum amount of time a connection can be reused
	maxConnLifetime := 2 * time.Minute
	dbConfig, err := pgxpool.ParseConfig(dsn)
	dbConfig.MaxConnIdleTime = maxConnLifetime
	dbConfig.MaxConns = int32(maxConnections)
	dbConfig.MinConns = int32(idleConn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	poolConn, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer poolConn.Close()

	// normal connection
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Initialize the HTTP router
	router := gin.Default()
	router.StaticFile("/", "./index.html")

	// One single connection for all requests
	router.GET("/products/normal", func(c *gin.Context) {
		startTime := time.Now()
		ctx := c.Request.Context()
		// Query the database for all products
		rows, err := conn.Query(ctx, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products, err := scanProducts(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		elapsed := time.Since(startTime).Microseconds()
		allCount++
		allTime += elapsed
		c.JSON(http.StatusOK, model.Response{Elapsed: elapsed, Average: float64(allTime / allCount), Products: products})
	})

	router.GET("/products/pooled", func(c *gin.Context) {
		ctx := c.Request.Context()
		startTime := time.Now()
		conn, err := poolConn.Acquire(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Query the database for all products
		rows, err := conn.Query(ctx, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products, err := scanProducts(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		elapsed := time.Since(startTime).Microseconds()
		poolCount++
		poolTime += elapsed
		c.JSON(http.StatusOK, model.Response{Elapsed: elapsed, Average: float64(poolTime / poolCount), Products: products})
	})

	router.GET("/products/new", func(c *gin.Context) {
		startTime := time.Now()
		ctx := c.Request.Context()
		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}

		rows, err := conn.Query(ctx, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products, err := scanProducts(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		elapsed := time.Since(startTime).Microseconds()
		newCount++
		newTime += elapsed
		c.JSON(http.StatusOK, model.Response{Elapsed: elapsed, Average: float64(newTime / newCount), Products: products})
	})

	// Start the HTTP server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Unable to start HTTP server: %v\n", err)
	}
}
