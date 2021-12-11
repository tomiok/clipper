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

The config should have the name and the max duration expressed in seconds.
The Function to operate,.
The fallback function, is not necessary at all.

### Wait the response
```go
out := make(chan bool, 1)
	var res *http.Response
	valChan := clipper.Do(&clipper.Configs{Name: "my_command"}, func() error {
		r, err := http.Get("http://www.google.com/robots.txt")
		res = r
		out <- true
		return err
	}, nil)

	select {
	case <-out:
		value := <-valChan
		if value == 0 {
			fmt.Println("no errors")
		} else {
			fmt.Println("some errors here")
		}
	}
	fmt.Println(res)
```

### No wait to the response
```go
clipper.Do(&clipper.Configs{Name: "my_command"}, func() error {
		r, err := http.Get("http://www.google.com/robots.txt")
		res = r
		return err
	}, nil)

	clipper.FillStats("my_command", true)
```

### Web
Check the function `ExposeMetrics()` and get all the metrics info. You must provide a valid command name. In the 
URL like: 
`localhost:8080/metrics?c=my_command`

--- 
### Tests

```shell
go test -v .
```
