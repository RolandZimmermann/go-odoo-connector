# Go Odoo Connector

A simple and efficient Go package for interacting with Odoo's XML-RPC API.

## Installation

```bash
go get github.com/RolandZimmermann/go-odoo-connector
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/RolandZimmermann/go-odoo-connector"
)

func main() {
    // Initialize the connector
    connector, err := odoo.NewConnector(
        "https://your-odoo-instance.com",
        "your-username",
        "your-api-key",
        "your-database",
    )
    if err != nil {
        log.Fatal(err)
    }

    // Search for CRM leads
    leads, err := connector.SearchReadRecords("crm.lead", odoo.SearchReadOptions{
        Fields: []string{"id", "name", "email_from"},
        Limit:  10,
        Order:  "create_date desc",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Process results
    for _, lead := range leads {
        fmt.Printf("Lead: %v\n", lead["name"])
    }
}
```

### Using Configuration File

Create a `config.json` file:

```json
{
    "url": "https://your-odoo-instance.com",
    "username": "your-username",
    "api_key": "your-api-key",
    "db": "your-database-name"
}
```

Then use it in your code:

```go
connector, err := odoo.NewConnectorFromConfig("config.json")
if err != nil {
    log.Fatal(err)
}
```

### Domain Filters

The package supports Odoo's domain filters for searching records:

```go
// Basic filter
opts := odoo.SearchReadOptions{
    Fields: []string{"id", "name"},
    Domain: []interface{}{
        []interface{}{"type", "=", "lead"},
    },
}

// AND condition (implicit)
opts := odoo.SearchReadOptions{
    Fields: []string{"id", "name"},
    Domain: []interface{}{
        []interface{}{"type", "=", "lead"},
        []interface{}{"stage_id.name", "=", "New"},
    },
}

// OR condition
opts := odoo.SearchReadOptions{
    Fields: []string{"id", "name"},
    Domain: []interface{}{"|",
        []interface{}{"type", "=", "lead"},
        []interface{}{"type", "=", "opportunity"},
    },
}
```

## Features

- Simple and intuitive API
- Configuration file support
- Comprehensive domain filter support
- Error handling with Go's error wrapping
- Thread-safe

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 