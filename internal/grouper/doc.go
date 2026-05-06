// Package grouper provides utilities for partitioning diff results
// into named groups using configurable key functions.
//
// Common key functions (ByType, ByStatus) are provided out of the box,
// and callers may supply their own KeyFunc for custom grouping strategies
// such as grouping by annotation label or tag value.
package grouper
