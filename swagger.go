package revelswagger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/revel/revel"
)

var spec Specification

func Init(path string) {
	// We need to load the JSON schema now
	fmt.Println("[SWAGGER]: Loading schema...")

	content, err := ioutil.ReadFile(path + "\\conf\\spec.json")

	if err != nil {
		fmt.Println("[SWAGGER]: Couldn't load spec.json.", err)
		return
	}

	err = json.Unmarshal(content, &spec)

	if err != nil {
		fmt.Println("[SWAGGER]: Error parsing schema file.", err)
	}
}

func Filter(c *revel.Controller, fc []revel.Filter) {
	method := spec.Paths[c.Request.URL.Path].Get

	if method == nil {
		c.Result = c.NotFound("No matching route found: " + c.Request.RequestURI)
		return
	}

	if err := c.SetAction("AppController", "Index"); err != nil {
		c.Result = c.NotFound(err.Error())
		return
	}

	// Action has been found & set, let's validate the parameters
	validateParameters(method.Parameters, c)

	if c.Validation.HasErrors() {
		var errors []string

		for _, e := range c.Validation.Errors {
			errors = append(errors, e.Message)
		}

		c.Result = c.RenderJson(map[string]interface{}{"errors": errors})
		return
	}

	fc[0](c, fc[1:])
}
