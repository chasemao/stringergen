package example

// run below command in bash
// stringergen -source=example.go -destination=example_stringer.go
type someStruct struct{}

var _ = &someStruct{}
