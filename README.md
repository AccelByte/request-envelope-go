[![Build Status](https://travis-ci.com/AccelByte/request-envelope-go.svg?branch=master)](https://travis-ci.com/AccelByte/request-envelope-go)

# request-envelope-go
Wrapper of request operations (ie logging, tracing)

# usage
```go
    import (
        "github.com/AccelByte/logger-go"
    )

    logger := loggergo.InitLogger("service-name", "realm")

    scope := NewRootScope(context.Background(), logger, "component-name", "trace-id")
    defer scope.Finish()

    scope.Logger.Errorf("err")
    // output: time="2020-04-03T16:37:45+07:00" level=error msg=err file="scope_test.go:21" func=logger-go.TestServiceNameLogged realm=realm service=service-name caller=component-name trace=xxxx span=yyyy
```