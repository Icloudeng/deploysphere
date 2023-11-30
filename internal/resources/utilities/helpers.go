package utilities

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/icloudeng/platform-installer/internal/filesystem"
	"github.com/icloudeng/platform-installer/internal/structs"
)

type helpers struct{}

var Helpers helpers

func (helpers) ConcatenateSubdomain(subDomains ...string) string {
	// Create a new slice to hold filtered subdomains
	var filteredSubdomains []string

	// Filter out any empty or zero-length subdomains
	for _, subDomain := range subDomains {
		if subDomain != "" {
			filteredSubdomains = append(filteredSubdomains, subDomain)
		}
	}

	// Join the filtered subdomains with a "." separator
	subDomainResult := strings.Join(filteredSubdomains, ".")

	return subDomainResult
}

func (helpers) ConcatenateAndCleanParams(params ...string) string {
	// Join the parameters with "-"
	result := strings.Join(params, "-")

	// Replace non-alphanumeric characters with "-"
	regex := regexp.MustCompile("[^a-zA-Z0-9]+")
	result = regex.ReplaceAllString(result, "-")

	return result
}

func (helpers) IsProdEnv(value string) bool {
	return value == "prod"
}

func (helpers) JoinIntSlice(slice []int, separator string) (int, error) {
	// Convert the slice elements to strings
	stringSlice := make([]string, len(slice))
	for i, num := range slice {
		stringSlice[i] = strconv.Itoa(int(num))
	}

	// Join the string slice with the separator
	joinedStr := strings.Join(stringSlice, separator)

	// Parse the result back into an integer
	result, err := strconv.Atoi(joinedStr)
	return result, err
}

func (h helpers) GenerateVMId(platform string, env string, params ...int) (int, error) {
	environments := structs.Environments{}

	content := filesystem.ReadEnvironmentsFile()
	if err := json.Unmarshal(content, &environments); err != nil {
		return 0, err
	}

	plaform_code, ok := environments.Platforms[platform]
	if !ok {
		return 0, errors.New("cannot found code for the specified Platform")
	}

	environment_code, ok := environments.Environments[env]
	if !ok {
		return 0, errors.New("cannot found code for the specified Environment")
	}

	params = append(params, environment_code, plaform_code)

	// Convert into string
	joinedInt, err := h.JoinIntSlice(params, "")
	if err != nil {
		return 0, err
	}

	return joinedInt, nil
}

func (helpers) ExtractSubdomainAndRootDomain(fqdn string) (subdomain, rootDomain string) {
	parts := strings.Split(fqdn, ".")
	if len(parts) > 1 {
		subdomain = strings.Join(parts[:len(parts)-2], ".")
		rootDomain = strings.Join(parts[len(parts)-2:], ".")
	} else if len(parts) == 1 {
		rootDomain = parts[0]
	}
	return subdomain, rootDomain
}

func (helpers) ExtractCommandKeyValuePairs(command string) map[string]string {
	keyValueMap := make(map[string]string)

	// Define a regular expression pattern to match key-value pairs.
	pattern := `--(\S+)\s+([^\s]+)`

	// Compile the regular expression pattern.
	regex := regexp.MustCompile(pattern)

	// Find all submatches in the input command string.
	matches := regex.FindAllStringSubmatch(command, -1)

	for _, match := range matches {
		key := match[1]
		value := match[2]
		// Remove any leading and trailing quotes if present.
		value = strings.Trim(value, `"'`)
		keyValueMap[key] = value
	}

	return keyValueMap
}

func RemoveFirstSegment(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) > 1 {
		return strings.Join(parts[1:], ".")
	}
	return ""
}
