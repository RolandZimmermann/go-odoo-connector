/*
Package odoo provides a connector for interacting with Odoo's XML-RPC API.

Example usage:

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

	// Search for recent CRM leads
	leads, err := connector.SearchReadRecords("crm.lead", odoo.SearchReadOptions{
		// Specify fields to retrieve
		Fields: []string{"id", "name", "email_from", "description"},
		// Filter by domain (optional)
		Domain: []interface{}{
			[]interface{}{"type", "=", "lead"},
			[]interface{}{"stage_id.name", "=", "New"},
		},
		// Limit results
		Limit: 10,
		// Order by creation date
		Order: "create_date desc",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Process the results
	for _, lead := range leads {
		fmt.Printf("Lead: %v (Email: %v)\n",
			lead["name"],
			lead["email_from"])
	}

	// Using configuration file
	connector, err := odoo.NewConnectorFromConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

Domain Examples:

	// Basic comparison
	[]interface{}{{"field", "=", value}}

	// AND condition (implicit)
	[]interface{}{
		[]interface{}{"field1", "=", value1},
		[]interface{}{"field2", "!=", value2},
	}

	// OR condition
	[]interface{}{"|",
		[]interface{}{"field1", "=", value1},
		[]interface{}{"field2", "=", value2},
	}

	// Complex condition (A AND (B OR C))
	[]interface{}{
		[]interface{}{"field1", "=", value1},
		"|",
		[]interface{}{"field2", "=", value2},
		[]interface{}{"field3", "=", value3},
	}

Common Operators:
	=, !=, >, >=, <, <=, like, ilike, in, not in, child_of, parent_of
*/
package odoo

import (
	"fmt"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// Connector represents an Odoo API connection
type Connector struct {
	URL      string
	Username string
	APIKey   string
	DB       string
	UID      int
	common   *rpc.Client
	models   *rpc.Client
}

// SearchReadOptions contains options for searching and reading records
type SearchReadOptions struct {
	Fields []string
	Domain []interface{}
	Offset int
	Limit  int
	Order  string
}

// NewConnector creates and initializes a new Odoo connector
func NewConnector(url, username, apiKey, db string) (*Connector, error) {
	c := &Connector{
		URL:      url,
		Username: username,
		APIKey:   apiKey,
		DB:       db,
	}

	// Initialize RPC clients
	var err error
	c.common, err = jsonrpc.Dial("tcp", fmt.Sprintf("%s/xmlrpc/2/common", url))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to common endpoint: %w", err)
	}

	c.models, err = jsonrpc.Dial("tcp", fmt.Sprintf("%s/xmlrpc/2/object", url))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to models endpoint: %w", err)
	}

	// Authenticate and get user ID
	var uid int
	err = c.common.Call("authenticate", []interface{}{db, username, apiKey, map[string]interface{}{}}, &uid)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	if uid == 0 {
		return nil, fmt.Errorf("authentication failed: invalid credentials")
	}

	c.UID = uid
	log.Printf("Successfully initialized Odoo connector with UID: %d", uid)
	return c, nil
}

// SearchReadRecords searches and reads records from Odoo
func (c *Connector) SearchReadRecords(model string, opts SearchReadOptions) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	if opts.Domain == nil {
		opts.Domain = []interface{}{}
	}

	params := map[string]interface{}{
		"fields": opts.Fields,
		"offset": opts.Offset,
		"limit":  opts.Limit,
		"order":  opts.Order,
	}

	err := c.models.Call("execute_kw", []interface{}{
		c.DB, c.UID, c.APIKey,
		model, "search_read",
		[]interface{}{opts.Domain},
		params,
	}, &result)

	if err != nil {
		return nil, fmt.Errorf("search_read failed for model %s: %w", model, err)
	}

	return result, nil
}
