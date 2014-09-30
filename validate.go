package revelswagger

import (
	"strconv"
	"strings"

	"github.com/revel/revel"
)

func validateParameters(params []parameter, c *revel.Controller) {
	for _, param := range params {
		validateType(param, c)

		if param.Required {
			c.Validation.Required(c.Params.Get(param.Name)).Message("%s is required.", param.Name)
		}

		// If the parameter is not required and is not set, then don't validate
		if !param.Required && c.Params.Get(param.Name) == "" {
			continue
		}

		if param.Minimum != nil {
			var val int

			c.Params.Bind(&val, param.Name)

			c.Validation.Min(val, *param.Minimum).Message("%s has to be at least %d", param.Name, *param.Minimum)
		}

		if param.Maximum != nil {
			var val int

			c.Params.Bind(&val, param.Name)

			c.Validation.Max(val, *param.Maximum).Message("%s has to be under %d", param.Name, *param.Maximum)
		}

		if len(param.Enum) > 0 {
			val := strings.ToLower(c.Params.Get(param.Name))

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
	val := c.Params.Get(param.Name)

	// Don't bother checking the type if it's empty
	// as we'll let the required deal with that first.
	if val == "" {
		return
	}

	var ok bool

	switch param.Type {
	case "number":
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
