package ldap

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"smatflow/platform-installer/pkg/env"

	"github.com/go-ldap/ldap/v3"
)

func connect() *ldap.Conn {
	l, err := ldap.DialURL(env.Config.LDAP_SERVER_URL)
	if err != nil {
		log.Fatalln(err)
	}

	return l
}

func LDAPExistBindUser(username string, password string) bool {
	connection := connect()
	defer connection.Close()

	user := parseUserBindTemplate(username)

	err := connection.Bind(user, password)

	return err == nil
}

func parseUserBindTemplate(username string) string {
	// Define the data to be used for replacing the template
	data := struct {
		Username string
	}{
		Username: username,
	}

	// Create a new template with the input string
	tmpl := template.Must(template.New("template").Parse(env.Config.LDAP_BIND_TEMPLATE))

	// Execute the template with the data and store the result in a buffer
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		log.Panicln(err)
	}

	// Retrieve the result from the buffer as a string
	return buf.String()
}

func init() {
	if env.Config.LDAP_AUTH {
		fmt.Print("LDAP Enabled!")
		// Test the conection
		connect().Close()
		// Test template parse
		parseUserBindTemplate("Test")
	}
}
