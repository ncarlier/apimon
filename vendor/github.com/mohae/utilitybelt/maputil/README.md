utilitybelt/maps
====================

This provides helper functions for the [Go maps](http://blog.golang.org/go-maps-in-action) datatype.

## ToSlices Functions
ToSlices functions take an incoming map datatype and returns two slices, keys and values, with their indexes matching, enabling using slices as the equivalent of maps. This is useful where *n* is small. The datatype of the map keys and values are the first two parts of the *ToSlices funcs, e.g. StringInterfaceToSlices creates slices out of a map of type `map[string]interface{}`.

Currently supported maps:
	map[string]string
	map[string]bool
	map[string]int
	map[string]interface{}
