package revelswagger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/revel/revel"
)

var spec Specification
var router *revel.Router

func Init(path string, r *revel.Router) {
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

	router = r
}

func Filter(c *revel.Controller, fc []revel.Filter) {
	var route *revel.RouteMatch = router.Route(c.Request.Request)

	if route == nil {
		c.Result = c.NotFound("No matching route found: " + c.Request.RequestURI)
		return
	}

	c.Params.Route = route.Params

	// Add the fixed parameters mapped by name.
	// TODO: Pre-calculate this mapping.
	for i, value := range route.FixedParams {
		if c.Params.Fixed == nil {
			c.Params.Fixed = make(url.Values)
		}
		if i < len(c.MethodType.Args) {
			arg := c.MethodType.Args[i]
			c.Params.Fixed.Set(arg.Name, value)
		} else {
			fmt.Println("Too many parameters to", route.Action, "trying to add", value)
			break
		}
	}

	leaf, _ := router.Tree.Find(treePath(c.Request.Method, c.Request.URL.Path))

	r := leaf.Value.(*revel.Route)

	method := spec.Paths[r.Path].Get

	if method == nil {
		c.Result = c.NotFound("No matching route found: " + c.Request.RequestURI)
		return
	}

	if err := c.SetAction(route.ControllerName, route.MethodName); err != nil {
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

func treePath(method, path string) string {
	if method == "*" {
		method = ":METHOD"
	}
	return "/" + method + path
}
