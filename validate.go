package revelswagger

import (
	"strconv"
	"strings"

	"github.com/revel/revel"
)

func validateParameters(params []parameter, c *revel.Controller) {
	for _, param := range params {
		validateType(param, c)

		value := c.Request.URL.Query().Get(param.Name)

		if param.Required {
			c.Validation.Required(value).Message("%s is required.", param.Name)
		}

		// If a default is to be set and param is empty, set it.
		if param.Default != "" && value == "" {
			c.Params.Route.Set(param.Name, param.Default)
		}

		// If the parameter is not required and is not set, then don't validate
		if !param.Required && value == "" {
			continue
		}

		if param.Minimum != nil {
			val, _ := strconv.Atoi(value)

			c.Validation.Min(val, *param.Minimum).Message("%s has to be at least %d", param.Name, *param.Minimum)
		}

		if param.Maximum != nil {
			val, _ := strconv.Atoi(value)

			c.Validation.Max(val, *param.Maximum).Message("%s has to be under %d", param.Name, *param.Maximum)
		}

		if len(param.Enum) > 0 {
			val := strings.ToLower(value)

			valid := false

			for _, e := range param.Enum {
				if val == strings.ToLower(e) {
					valid = true
					break
				}
			}

			if !valid {
				c.Validation.Error("%s needs to one of '%s'", param.Name, strings.Join(param.Enum, ","))
			}
		}
	}
}

func validateType(param parameter, c *revel.Controller) {
	val := c.Request.URL.Query().Get(param.Name)

	// Don't bother checking the type if it's empty
	// as we'll let the required deal with that first.
	if val == "" {
		return
	}

	var ok bool

	switch param.Type {
	case "integer":
		switch param.Format {
		case "int32":
			_, err := strconv.ParseInt(val, 10, 32)
			ok = err == nil
			break
		case "int64":
			_, err := strconv.ParseInt(val, 10, 64)
			ok = err == nil
			break
		}
		break
	case "string":
		ok = true
		break
	case "boolean":
		ok = val == "true" || val == "false"
		break
	}

	if !ok {
		c.Validation.Error("%s needs to be a %s (%s)", param.Name, param.Type, param.Format)
	}
}
