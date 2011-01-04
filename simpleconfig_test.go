package simpleconfig

import (
    "reflect"
    "testing"
)

type TestConfig struct {
    Option1 string
    Option2 bool
    Option3 int
    Option4 int64
    Option5 bool
    Option6 string
    Option7 float64
}

var expected = TestConfig{
    "/path/to/arbitrary/file",
    true,
    100,
    -1,
    false,
    "this has multiple values",
    19.01,
}

var testString = `
#This is an example configuration file

#Blank lines or lines that begin with a pound sign are ignored

option1 /path/to/arbitrary/file
option2 true
Option3 100
Option4 -1
Option5 false
option6 this has multiple values
option7 19.01
`

func TestBasic(t *testing.T) {
    var actual TestConfig
    config := ReadString(testString)
    config.Unmarshal(&actual)
    if !reflect.DeepEqual(expected, actual) {
        t.Fatalf("Actual %v expected %v", actual, expected)
    }
}
