# go-dumper

Recursive dump routine based on ancient dumper at https://code.google.com/p/golang/ .

## SYNOPSIS

Usage:

```
package main

import (
    "time"
    "github.com/acidlemon/go-dumper"
)

func main() {
    now := time.Now().Local()
    dump.Dump(now)
}
```

Output:

```
Time {
  sec: 63539975841 (int64),
  nsec: 0x6200c19 (uintptr),
  loc: (0x0014c3a0) &Location {
    name: "",
    zone: nil ([]time.zone),
    tx: nil ([]time.zoneTrans),
    cacheStart: 0 (int64),
    cacheEnd: 0 (int64),
    cacheZone: nil (*time.zone)
  }
}
```

## LICENSE

[Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0) (same with original code)

