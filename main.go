package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mineiros-io/terradoc/config"
)

func main() {

	module := config.ParseVariables("fixtures.tf")

	// 	- **`name`**: *(Required `string`, Forces new resource)*
	// 	The name of the user pool.

	//   - **`advanced_security_mode`**: *(Optional `string`)*
	// 	The mode for advanced security, must be one of `OFF`, `AUDIT` or `ENFORCED`. Additional pricing applies for Amazon Cognito advanced security features. For details see https://aws.amazon.com/cognito/pricing/.
	// 	Default is `OFF`.

	for _, v := range module.Variables {
		// please don't complain about hacky naming yet :D
		re := regexp.MustCompile("^\\((.*?)\\)")
		a := re.FindStringSubmatch(v.Description)

		// Shadowing the same variable inside the same scrope doesn't work in go
		re = regexp.MustCompile("\\n")
		variableType := re.ReplaceAllString(v.Type, "")

		fmt.Printf("**`%s`**: *(%s `%s`, Forces new resource)*\n", v.Name, a[1], variableType)
		fmt.Printf("%s", strings.TrimLeft(v.Description, a[0]))
		fmt.Print("\n\n")

	}
}
