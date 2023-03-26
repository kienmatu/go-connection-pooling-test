package main

import (
	"context"
	"github.com/kienmatu/go-connection-pooling/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var allTime int64 = 0
var allCount int64 = 0

var poolTime int64 = 0
var poolCount int64 = 0

var dsn = "postgres://postgres:password1@localhost:5433/postgres"
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
	idleConn := 10
	// Set the maximum number of connections in the pool
	maxConnections := 30
	// Set the maximum amount of time a connection can be reused
	maxConnLifetime := 5 * time.Minute
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
	}
	config.MaxConns = int32(maxConnections)
	config.MinConns = int32(idleConn)
	config.MaxConnLifetime = maxConnLifetime
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()
	// normal connection
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Initialize the HTTP router
	router := gin.Default()
	router.StaticFile("/", "./index.html")

	router.GET("/products/normal", func(c *gin.Context) {
		startTime := time.Now()

		// Query the database for all products
		rows, err := conn.Query(context.Background(), query)
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
		startTime := time.Now()
		// Get a connection from the pool
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to acquire database connection"})
			return
		}
		defer conn.Release()

		// Query the database for all products
		rows, err := conn.Query(context.Background(), query)
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
		// normal connection
		conn, err := pgx.Connect(context.Background(), dsn)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}

		// Query the database for all products
		rows, err := conn.Query(context.Background(), query)
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

	// Start the HTTP server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Unable to start HTTP server: %v\n", err)
	}
}
