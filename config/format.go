package config

import (
	"fmt"
	"regexp"
	"strings"
)

// Format takes a Module struct and turns it into a human-readable output
// TODO: Should be based on an interface and expect structured information
func Format(module *Module) strings.Builder {
	outputBuffer := strings.Builder{}

	descriptionRegEx := regexp.MustCompile("^\\((.*?)\\)")
	variableTypeRegEx := regexp.MustCompile("\\n")

	for _, v := range module.Variables {
		description := descriptionRegEx.FindStringSubmatch(v.Description)
		variableType := variableTypeRegEx.ReplaceAllString(v.Type, "")

		// Add type and level ( Required vs Optional)
		// TODO: 'Hardcoded forces new resources' should be generated dynamically
		outputBuffer.WriteString(fmt.Sprintf("**`%s`**: *(%s `%s`, Forces new resource)*\n\n", v.Name, description[1], variableType))

		// Add description to variable block
		outputBuffer.WriteString(fmt.Sprintf("- **`%s`**: *(%s `%s`)*\n", v.Name, description[1], variableType))

		// Replace backticks with ticks
		// TODO: optimize repititions
		outputBuffer.WriteString(fmt.Sprintf(" %s ", strings.Replace(strings.TrimLeft(v.Description, description[0]), "'", "`", 99)))

		// Add default value
		// outputBuffer.WriteString(fmt.Sprintf("Default is %s", ValueStr(v.Default)))

		// Add line breaks to variable block
		outputBuffer.WriteString("\n\n")
	}

	return outputBuffer
}
