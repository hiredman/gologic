# Gologic

A port of core minikanren with inequality to Go based on doctor Byrd's
thesis.

Provides reasoning over structs via Go's `reflect` package.

Provides a "database" in to which you can assert facts and then reason
over the facts in the database

## Usage

```go
a,b,c := Fresh3()
ch := Run(a,Or(And(Unify(a,b),Unify(b,c),Unify(c,1)),
               And(Unify(a,b),Unify(b,c),Unify(c,3))))
if 1 == <- ch {t.Fatal("not a 1")}
if 3 == <- ch {t.Fatal("not a 3")}
```

`gologic.Fresh()` returns a fresh unbound logic var. Fresh2-6 return
more logic vars at once.

`gologic.Run(...)` takes a logic var and a goal, and returns a channel
(type `chan interface{}`) in to which results will be placed.

provided goals constructors are:
  * `gologic.Unify(a,b)`
  * `gologic.Neq(a,b)`
  * `gologic.And(a,...)`
  * `gologic.Or(a,...)`

a goal is a `func (s S) R {...}`

a goal constructor is a function that returns a goal

See `gologic_test.go`

`zebra.go` has an example of solving the zebra puzzle ported from
Clojure's `core.logic`. You can run the example like:

```sh
GOPATH=$PWD go build examples/zebra.go && time ./zebra 
```

`gologic.Call` is a goal constructor that is a useful helper for
constructing recursive goals, by delaying the recursive call to the
goal constructor until the logic program is being run. See `ancestoro`
in `gologic_test.go`

## Benchmarks

```sh
 GOPATH=$PWD go test -bench ".*" gologic 
```

## License

The MIT License (MIT)

Copyright © 2013 Kevin Downey

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

