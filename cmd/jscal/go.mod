module github.com/airtrafik/jscal/cmd/jscal

go 1.23

require (
	github.com/airtrafik/jscal v0.1.0
	github.com/airtrafik/jscal/convert/ical v0.1.0
)

require github.com/arran4/golang-ical v0.3.2 // indirect

replace github.com/airtrafik/jscal => ../../

replace github.com/airtrafik/jscal/convert/ical => ../../convert/ical
