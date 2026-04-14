**Countries Service**

Backend web application written in Go for fetching, storing, and displaying country data.

**Features**
  Fetch data from REST Countries API
  Store data in ClickHouse
  Search country by name
  View all countries
  Population statistics (total, average, most populated)
  Server-side rendered UI (HTML templates)

**Tech Stack**
  Go, fasthttp
  ClickHouse
  HTML, CSS

**Architecture**
  Client → Router → Handlers → (API / ClickHouse) → HTML response
  handlers.go — request handling
  fetch.go — external API
  clickhouse.go — database
  router.go — routes

**Notes**
  ClickHouse used for learning purposes
  Simple architecture, no separate service layer

**Improvements**
  Batch inserts
  Better separation of layers
  Validation and error handling
