
# UIXDB

> A lightweight visual web workspace for PostgreSQL database administration and design, combining the power of an SQL console, a table editor, and an interactive ERD schema designer.

<img width="2560" height="1468" alt="image" src="https://github.com/user-attachments/assets/6cafd9e0-56a6-4d31-b4d5-3c1709e53c18" />

---

## Overview

**UIXDB** is built for engineers and database designers who want a streamlined, developer-first tool for inspecting and managing PostgreSQL databases. Instead of switching between heavy database clients and external diagramming tools, UIXDB brings relationship visualization, data browsing, and transactional SQL execution into a single, unified interface.

---

## Key Features

* **Interactive ER Diagram:** Automatic visualization of tables, fields, Primary Keys, and Foreign Keys. Parses foreign key constraints into a directed graph for clear relationship modeling.
* **Transactional SQL Console:** Runs batch statements inside an isolated transaction (`BEGIN ... COMMIT / ROLLBACK`), dynamically parses result sets (`SELECT`), and returns structured records.
* **Deep System Catalog Inspection:** Direct queries against PostgreSQL internals (`pg_constraint`, `pg_attribute`, `information_schema`) for fast, precise extraction of schema metadata.
* **Data Grid Viewer:** Inspect records in a flexible tabular view with robust handling of `NULL` values, types, and raw query output.

---

## Tech Stack

* **Backend:** Go 1.26.4, [Gin Web Framework](https://github.com/gin-gonic/gin)
* **Database Driver:** `database/sql`, [`github.com/lib/pq`](https://github.com/lib/pq)
* **Target Database:** PostgreSQL 13+
* **Frontend:** Interactive canvas & graph engine for ER diagrams and responsive data grids

---

## Project Structure

```text
├── assets/
│   └── img/
│       └── logo.png         # Project branding and UI logo
├── databases/
│   └── postgres/
│       ├── connect.go       # PostgreSQL connection pool lifecycle (*sql.DB)
│       ├── models.go        # Request/response structs and database schema types
│       ├── send_console.go  # Batch runner with transaction rollback support
│       └── tables.go        # Catalog inspection (tables, columns, foreign keys)
├── templates/
│   └── index.html           # Web UI workspace
├── go.mod                   # Module definition
├── go.sum                   # Dependency checksums
└── main.go                  # Routing and HTTP server entry point
```

---

## Getting Started

### Prerequisites

* **Go** 
* **PostgreSQL** 

### Installation & Run

1. **Clone the repository:**
```bash
git clone https://github.com/DipasDi/UIXDB.git
cd UIXDB
```


2. **Download dependencies:**
```bash
go mod download
```


3. **Start the server:**
```bash
go run main.go
```



By default, the API will be available at `http://localhost:8082`.

---

## API Overview

### 1. Execute SQL Statements

Executes a list of SQL queries within a single transaction. Returns result rows if a `SELECT` query is provided.

* **Endpoint:** `POST /api/console`
* **Request Payload:**
```json
{
  "statements": [
    "CREATE TABLE users (id SERIAL PRIMARY KEY, username VARCHAR(50));",
    "INSERT INTO users (username) VALUES ('dipas');",
    "SELECT id, username FROM users;"
  ]
}

```


* **Response:**
```json
{
  "success": true,
  "rows": [
    ["1", "dipas"]
  ]
}
```



---

### 2. Retrieve Foreign Key Relationships

Extracts foreign key constraints from PostgreSQL catalog tables to power the ERD viewer.

* **Endpoint:** `GET /api/relations`
* **Response:**
```json
{
  "connected": true,
  "Key": [
    {
      "id": "fk_posts_user_id",
      "from": "posts.user_id",
      "to": "users.id"
    }
  ]
}

```




---
