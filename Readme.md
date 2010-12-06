== Description ==

Simpleconfig is an simple configuration system you can add to any Go project. It will automatically try to read configuration from a source -- a file, a string, or an io.Reader. 

The simpleconfig format is very, well, simple. It ignores lines that are blank or start with the pound sign `#`. Otherwise, it expects lines to have the format `key value`, where key and value are strings separated by a blank space. 

simpleconfig has two methods: 

As an example, see [redis.conf](https://github.com/antirez/redis/blob/master/redis.conf)

== Installation == 

The simpleconfig source code is arranged as a go package, so you can install it in the following fashion:

    * `goinstall github.com/hoisie/simpleconfig`
    * Cloning the repository and running 'make install'

Because it's such a small package, it might be easier to copy and paste simpleconfig.go to your project, and simply change the package name. However, if there are updates to the package you might miss them. 

== Usage == 

Simpleconfig has two methods

    Read(source interface{}) (map[string]string, os.Error)
    Unmarshal(dst interface{}, source interface{}) os.Error

The first method takes a source (either a string, or an io.Reader), and returns a map[string]string with the configuration. If the source argument is a string, it tries to open a file with the name, and if that fails, it treats the string itself as the configuration. 

The second method is more useful, it takes the configuration source and tries to write it to a struct. For example, if you have a struct holding configuration variables, and a string with the values:

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

