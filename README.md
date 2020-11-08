# distributed-lock
The distributed-lock is a simple distributed lock.This version was developed by redigo.
You can find locks developed with go-redis in my other projects.

## download and install
```
$ go get github.com/newneod/distributed-lock-redigo
```

## demo
```
package main

import "time"

// demo
func main() {
	Init("127.0.0.1:6379")

	strUUID, err := Lock("a")
	if err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)
	err = Unlock("a", string(strUUID))
	if err != nil {
		panic(err)
	}
}
```