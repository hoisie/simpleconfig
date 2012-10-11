## *Changes*

I use .properties file format to replace hoisie's config format. This .properties format is more common, especially to Java developer. 

now:

    host = 192.168.0.1
    port = 80

origin:

    host 192.168.0.1
    port 80

I change public methods, so that is more simple to use

    func Read(dst interface{}, src io.Reader) error
	func ReadString(dst interface{}, src string) error
	func ReadFile(dst interface{}, src string) error
	
I update these codes to go1 compatible 

## Description 

Simpleconfig is a library that gives your Go project a flexible configuration mechanism. It is inspired by the configuration system for properties. It can read config from different sources (a file, string, or io.Reader), and write it to an arbitrary struct or map.

The simpleconfig format is primitive -- it ignores lines that are blank or start with a pound sign `#`. Otherwise, it expects lines to have the format `key value`, where key and value are strings separated by an equal mark (=), it's similar with .properties file' . 


## Installation

Simpleconfig is a go package, so it can be installed with:

 * `go get github.com/jijinggang/simpleconfig`
 * Clone the repo and run 'make install'
 * Or, just copy and paste simpleconfig.go into a new file in your project to avoid the dependency (but you won't get updates)

## Usage

Simpleconfig has three methods:

    func Read(dst interface{}, src io.Reader) error
	func ReadString(dst interface{}, src string) error
	func ReadFile(dst interface{}, src string) error

The three methods are similar, only use different input source. it takes the configuration contained in `src` and tries to write it to the value represented by `dst`. For example, if you have a struct holding configuration variables, and a string with the values:

    type Config struct{
        Option1 string
        Option2 bool
        Option3 int64
    }

    var configString = `
    #Sample configuration
    
    option1 = hello
    option2 = true
    option3 = 101
    `

You can write:

    var config Config
    err := simpleconfig.ReadString(&config, configString)
