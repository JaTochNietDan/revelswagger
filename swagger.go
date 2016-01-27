package revelswagger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path"
	"runtime"
	"strings"
	"text/template"

	"github.com/howeyc/fsnotify"
	"github.com/revel/revel"
)

var spec Specification

func init() {
	revel.OnAppStart(func() {
		fmt.Println("[SWAGGER]: Loading schema...")

		loadSpecFile()

		go watchSpecFile()
	})
}

func loadSpecFile() {
	spec = Specification{}
    specPath := "/conf/spec.json"


	if runtime.GOOS == "windows" {
        specPath = "\\conf\\spec.json"
    }

	content, err := ioutil.ReadFile(revel.BasePath + specPath)

	if err != nil {
		fmt.Println("[SWAGGER]: Couldn't load spec.json.", err)
		return
	}

	err = json.Unmarshal(content, &spec)

	if err != nil {
		fmt.Println("[SWAGGER]: Error parsing schema file.", err)
	}
}

func watchSpecFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case <-watcher.Event:
				loadSpecFile()
			case err := <-watcher.Error:
				fmt.Println("[SWAGGER]: Watcher error:", err)
			}
		}
	}()

    specPath := "/conf/spec.json"
    if runtime.GOOS == "windows" {
        specPath = "\\conf\\spec.json"
    }

	err = watcher.Watch(revel.BasePath + specPath)

	if err != nil {
		fmt.Println("[SWAGGER]: Error watching spec file:", err)
	} else {
		fmt.Println("[SWAGGER]: Spec watcher initialized")
	}

	<-done

	/* ... do stuff ... */
	watcher.Close()
}

func Filter(c *revel.Controller, fc []revel.Filter) {
	c.Request.URL.Path = strings.ToLower(c.Request.URL.Path)

	var route *revel.RouteMatch = revel.MainRouter.Route(c.Request.Request)

	if route == nil {
		c.Result = c.NotFound("No matching route found: " + c.Request.RequestURI)
		return
	}

	if len(route.Params) == 0 {
		c.Params.Route = map[string][]string{}
	} else {
		c.Params.Route = route.Params
	}

	if err := c.SetAction(route.ControllerName, route.MethodName); err != nil {
		c.Result = c.NotFound(err.Error())
		return
	}

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

	leaf, _ := revel.MainRouter.Tree.Find(treePath(c.Request.Method, c.Request.URL.Path))

	r := leaf.Value.(*revel.Route)

	method := spec.Paths[r.Path].Get

	if method == nil {
		// Check if strict mode is enabled and throw an error, otherwise
		// just move onto the next filter like revel normally would
		if revel.Config.BoolDefault("swagger.strict", true) {
			_, filename, _, _ := runtime.Caller(0)

			t, err := template.ParseFiles(path.Dir(filename) + "/views/notfound.html")

			if err != nil {
				panic(err)
			}

			t.Execute(c.Response.Out, map[string]interface{}{
				"routes": spec.Paths,
				"path":   c.Request.RequestURI,
			})
			return
		} else {
			// Move onto the next filter
			fc[0](c, fc[1:])
			return
		}
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

	// Move onto the next filter
	fc[0](c, fc[1:])
}

func treePath(method, path string) string {
	if method == "*" {
		method = ":METHOD"
	}
	return "/" + method + path
}
