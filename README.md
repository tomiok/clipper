# Clipper - circuit breaker

### Lightweight circuit breaker tool, inspired in hystrix from Netflix.


----

This is a simple, dependency free, circuit breaker library. Super easy to use.
You could find also some basic statistics and a handler according to the std lib to expose, given the name of the
command.


```go
clipper.Do(&clipper.Configs{Name: "my_command"}, func() error {
		_, err := http.Get("http://www.google.com/robots.txt")
		return err
	}, nil)
```
Do use 3 args,
1. The Config
2. Operational function
3. Fallback function

And in this case you also have 2 different ways to call the clipper , sync and async, could check out in the 
examples folder.