package lib

import (
	"bytes"
	"html/template"
	"log"

	"github.com/go-ldap/ldap/v3"
)

func connect() *ldap.Conn {
	l, err := ldap.DialURL(EnvConfig.LdapServerUrl)
	if err != nil {
		log.Fatal(err)
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
	tmpl := template.Must(template.New("template").Parse(EnvConfig.LdapBindTemplate))

	// Execute the template with the data and store the result in a buffer
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	// Retrieve the result from the buffer as a string
	return buf.String()
}

func init() {
	// Test the conection
	connect().Close()
	// Test template parse
	parseUserBindTemplate("Test")
}
