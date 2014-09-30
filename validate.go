package revelswagger

import (
	"strconv"

	"github.com/revel/revel"
)

func validateParameters(params []parameter, c *revel.Controller) {
	for _, param := range params {
		validateType(param, c)

		if param.Required {
			c.Validation.Required(c.Params.Get(param.Name)).Message("%s is required.", param.Name)
		}
	}
}

func required(data string, c *revel.Controller) {

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
