const { Pool, Client } = require('pg');
const path = require("path")
const express = require('express');

const pool = new Pool({
    connectionString: 'postgres://postgres:password1@localhost:5433/postgres?sslmode=disable',
    max: 90, // maximum number of connections
    idleTimeoutMillis: 2 * 60 * 1000, // maximum amount of time a connection can be idle
});
const conn = new Client({
    connectionString: 'postgres://postgres:password1@localhost:5433/postgres?sslmode=disable',
    max: 1,
});
const query = 'SELECT id, name, price, description FROM products limit 1000';

const app = express();

app.use((req, res, next) => {
    console.log(`[${new Date().toISOString()}] ${req.method} ${req.url}`);
    next();
});

//assuming app is express Object.
app.get('/',function(req,res) {
    res.sendFile('index.html', {root: path.resolve(__dirname, "../") });
});

let allTime = 0;
let allCount = 0;

let newTime = 0;
let newCount = 0;

let poolTime = 0;
let poolCount = 0;
app.get('/products/normal', async (req, res) => {
    const start = Date.now();
    try {
        console.log('Query started');
        const { rows } = await conn.query(query);
        console.log('Query finished');
        const products = rows;
        const elapsed = Date.now() - start;
        allCount++;
        allTime += elapsed;
        res.json({
            elapsed,
            average: (allTime + elapsed) / (allCount + 1),
            products,
        });
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: err.message });
    }
    finally {
        conn.release();
    }
});

app.get('/products/pooled', async (req, res) => {
    const start = Date.now();
    const client = await pool.connect();
    try {
        const { rows } = await client.query(query);
        const products = rows;
        const elapsed = Date.now() - start;
        poolCount++;
        poolTime += elapsed;
        res.json({
            elapsed,
            average: (poolTime + elapsed) / (poolCount + 1),
            products,
        });
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: err.message });
    }
    finally {
        client.release();
    }
});

app.get('/products/new', async (req, res) => {
    const start = Date.now();
    const conn = new Client({
        connectionString: 'postgres://postgres:password1@localhost:5433/postgres?sslmode=disable',
        max: 1,
    });
    try {
        const { rows } = await conn.query(query);
        const products = rows;
        const elapsed = Date.now() - start;
        newCount++;
        newTime += elapsed;
        res.json({
            elapsed,
            average: (newTime + elapsed) / (newCount + 1),
            products,
        });
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: err.message });
    } finally {
        conn.end();
    }
});

app.listen(8080, () => {
    console.log('Server started on port 8080');
});