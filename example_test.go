package odoo_test

import (
	"fmt"
	"log"

	"github.com/RolandZimmermann/go-odoo-connector"
)

func Example() {
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
		Domain: []interface{}{
			[]interface{}{"type", "=", "lead"},
		},
		Limit: 10,
		Order: "create_date desc",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Process results
	for _, lead := range leads {
		fmt.Printf("Lead: %v\n", lead["name"])
	}
}

func ExampleNewConnectorFromConfig() {
	// Create a connector from config file
	connector, err := odoo.NewConnectorFromConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Use the connector
	opts := odoo.SearchReadOptions{
		Fields: []string{"id", "name"},
		Domain: []interface{}{
			[]interface{}{"stage_id.name", "=", "New"},
		},
	}

	leads, err := connector.SearchReadRecords("crm.lead", opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, lead := range leads {
		fmt.Printf("Lead: %v\n", lead["name"])
	}
}
