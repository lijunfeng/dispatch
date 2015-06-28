Dispatch [![Build Status](https://drone.io/github.com/tango-contrib/dispatch/status.png)](https://drone.io/github.com/tango-contrib/dispatch/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/dispatch)](http://gocover.io/github.com/tango-contrib/dispatch)
======

Dispacth is a handler for dipatch http request according url's prefix.

# Example

```Go
import (
    "github.com/lunny/tango"
    "github.com/tango-contrib/dispatch"
)

func main() {
    logger := tango.NewLogger(os.Stdout)
    t1 := tango.NewWithLog(logger)
    t2 := tango.NewWithLog(logger)

    dispatch := dispatch.New(map[string]*tango.Tango{
        "/": t1,
        "/api/": t2,
    })

    t3 := tango.NewWithLog(logger)
    t3.Use(dispatch)
    t3.Run(":8000")
}
```


```Go
package main

import (
    "github.com/Unknwon/macaron"
    "github.com/lunny/tango"
    "github.com/tango-contrib/dispatch"
)

func main() {

    t := tango.Classic()

    t.Any("/favicon.ico", func(self *tango.Context) {
        self.Redirect("/static/favicon.ico", 301)
    })
    t.Any("/", func() string {
        return "Hello Tango!"
    })

    m := macaron.Classic()
    m.Get("/", func(ctx *macaron.Context) string {
        return "Macaron on Tango!"
    })

    dispatch := dispatch.Use("/", t)
    dispatch.Add("/m/", m)

    tan := tango.Classic()
    tan.Use(dispatch)
    tan.Run(80)
}
```