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

`gologic.A` is a value that will unify with anything, but will not
extend the substitution map.

provided goals constructors are:
  * `gologic.Unify(a,b)`
  * `gologic.Neq(a,b)`
  * `gologic.And(a,...)`
  * `gologic.Or(a,...)`

a goal is a `func (s S) R {...}`

a goal constructor is a function that returns a goal

See `gologic_test.go`

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

