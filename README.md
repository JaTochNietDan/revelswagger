Revel Swagger
============

## What is it?
Revel swagger is a plug and play filter designed for the full stack web framework called [Revel](http://github.com/revel/revel). The purpose of the filter is to parse a JSON file which follows the [JSON Schema](http://json-schema.org/) standard, in particular following the [Swagger](http://swagger.io) documented standard.

The idea is that you define your routes, the route parameters and some validation (if necessary) for the parameters and it will all be pre-validated. Therefore your code is based off of the spec. From this, there are some major advantages:

- API documentation can be generated from the spec and is guaranteed to be 100% accurate and up-to-date
- Simple validation is all handled in the same way and is consistent throughout your application
- Thanks to basic validation being done for you, you don't have to worry about it once you write the spec, you can get straight to more complex code
- Tests can be automated based on your spec file

## What works so far?

Currently only GET calls are working (need to find a better way of doing this potentially) and simple validation.

## Installation
All you need to do to install this is add revelswagger.Filter to your Revel filters (in place of the normal Revel filter). It is however important that you order it like so:

```
revel.Filters = []revel.Filter{
  filters.PanicFilter,           // Recover from panics and display an error page instead.
  revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
  revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
  revelswagger.Filter,           // Use the revelswagger filter to set up routes
  // any other filters
```

The reason for this is that we need the first three filters to be initialized before revelswagger in order to make it work and make use of Revel features within the filter.

Once you've done this, you need to add a spec.json file to your Revel's app/conf folder. You can see an example specification on [Swagger's editor](http://editor.swagger.wordnik.com/#/edit). Remember, currently only JSON is supported so you must format it as json.

## Configuration
Currently only one configuration option is available which is `swagger.strict` which defaults to true. If set to false, then only routes that are defined in spec.json will be validated. Any other routes which happen to be defined in your normal Revel routes file will still go through.

## Current major issue
The idea is to completely replace Revel's routes file, however we still rely on Revel's route file to tell us what controller and method to route our action to. If someone has an eloquent solution to this, I would like to hear it!

It would be wonderful to only define the routes in one, centralized place, the spec file.
