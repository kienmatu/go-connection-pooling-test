package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/kienmatu/go-connection-pooling/model"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var allTime int64 = 0
var allCount int64 = 0

var newTime int64 = 0
var newCount int64 = 0

var poolTime int64 = 0
var poolCount int64 = 0

var dsn = "postgres://postgres:password1@localhost:5433/postgres?sslmode=disable"
var query = "SELECT id, name, price, description FROM products limit 1000"

func scanProducts(rows *sql.Rows) ([]*model.Product, error) {
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
	poolConn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer poolConn.Close()
	poolConn.SetMaxOpenConns(maxConnections)
	poolConn.SetMaxIdleConns(idleConn)
	poolConn.SetConnMaxLifetime(maxConnLifetime)

	// normal connection
	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	// default will be 2 idle connections
	// so set it to 1 to simulate
	conn.SetMaxIdleConns(1)

	// Initialize the HTTP router
	router := gin.Default()
	router.StaticFile("/", "./index.html")

	router.GET("/products/normal", func(c *gin.Context) {
		startTime := time.Now()

		// Query the database for all products
		rows, err := conn.Query(query)
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
		// Query the database for all products
		rows, err := poolConn.Query(query)
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
		conn, err := sqlx.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}

		rows, err := conn.Query(query)
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
