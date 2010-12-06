## Description 

Simpleconfig is a library that gives your Go project a flexible configuration mechanism. It can read config from different sources (a file, string, or io.Reader), and write it to an arbitrary struct or map.

The simpleconfig format is primitive -- it ignores lines that are blank or start with a pound sign `#`. Otherwise, it expects lines to have the format `key value`, where key and value are strings separated by a blank space. 

As an example of a configuration file, see [redis.conf](https://github.com/antirez/redis/blob/master/redis.conf)

## Installation

Simpleconfig is a go package, so it can be installed with:

 * `goinstall github.com/hoisie/simpleconfig`
 * Clone the repo and run 'make install'
 * Or, just copy and paste simpleconfig.go into a new file in your project to avoid the dependency (but you won't get updates)

## Usage

Simpleconfig has two methods:

    func Read(source interface{}) (map[string]string, os.Error)
    func Unmarshal(dst interface{}, source interface{}) os.Error

The first method takes a `source` (either a string, or an io.Reader), and returns a `map[string]string` holding the configuration. If the source argument is a string, it tries to open a file by the name, and if that fails, it treats the string itself as the configuration data.

The second method is more useful, it takes the configuration contained in `source` and tries to write it to the value represented by `dst`. For example, if you have a struct holding configuration variables, and a string with the values:

    type Config struct{
        Option1 string
        Option2 bool
        Option3 int64
    }

    var configString = `
    #Sample configuration
    option1 hello
    option2 true
    option3 101
    `

You can write:

    var config Config
    err := simpleconfig.Unmarshal(&config, configString)
