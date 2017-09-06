# A Tour of Go

## Packages, variables and functions

- all programs have a package `main` is default?
- *factored* imports are preferred `import ("fmt", "math")`
- package name is the same as last element of import path `math/rand`
- package exports begin with a capital letter `math.Pi`
- type declaration follow name `a int, b int` (shortened to `a, b int`)
- `func` can return multiple values, declared at end of sig
- `func` return values can be named and function as if declared at top of func
- "naked" returns the named return values (best for shor funcs)
- `var` declares a list of variables (type at end), package or function level
- `var` can take assignments after `var i, j int = 1, 2`, can omit type
- inside a function `var` and assignment can be shortened with `:=` eg `k := 3` type is infered
- when doing multiple assignment `a, b = b, a + b` it is if they are done concurrently, not from left to right
- basic types include `bool, string, <u>int<8,16,32,64>, loat<32,64>, complex<32,64>, uintptr`
- and alias types `byte=uint8`, `rune=int32` (unicode cp)
- var can also be factored like import
- non initialized variables get the zero value (`0, false, ""` numer, bool, string)
- no implicit conversion between types, instead *TYPE(v)* eg `i := uint(f)`
- constants are declated with `const` and can't use `:=`
- numeric constants are high-precision, if they are untyped they take the type depending on the context

## Flow control
- `for` is only loop, standard semi-colon syntax
  - no parentheses and braces are required
  - init/post are optional `for ; i < x; { }`, but can shorten `for i < x { }`
  - infinite loop, simply `for { }`
- `if` also have no parentheses and braces are required
  - can also have init part `if i := 20; a > i { }`, only scoped to if block or any associated else blocks
- `switch` has a similar format, `swithch <short statement>; <switch var> { }`
  - case statements break automatically, unless use `fallthrough`
  - case conditions can be statements
  - no condition `switch {}` is same as `switch true {}`, since it flows top to bottom can use for if-else chains
- `defer` statement postpones function execution until surrounding function returns
  - the deferred functions args are evaluated immeadiatly
  - deferred function go onto a LIFO stack

## More Types
- **Pointers** hold the memory address of a value, defined with `var p *int`, has zero value of `nil`
  - `&` creates a pointer to its operand, `p := &i`
  - `*` dereferences the pointer to the actual value, `*p`
- **Structs** are a collection of fields, no functions allowed
  - `type Vertex struct { X, Y int }` create with `v := Vertex{1, 2}`
  - fields are accessed using the dot operator `v.X`
  - a struct pointer can access fields as `(*p).X` or the shorter `p.X`
  - a struct literal list the values of its fields, can be empty or use named fields `Vertex{X: 1}`
  - can also create a pointer to a literal with `&Vertex{1, 2}`
- **Arrays** are denoted `[n]T` *n* values of type *T*
  - declared as `var a [10]int` 10 element int array
  - length is part of type and cannot be changed
  - access with standard index notation, `a[0]`
  - can initialize and declare at once, `a := [3]int{1, 2, 3}`
- A **Slice** is a dynamically sized view into an array `[]T`
  - create a slice by specifying a low and high bound of an array `a[lo:hi]`, *hi* is exclusive
  - in practice more common then arrays
  - slices do not store data and are references or views of a part of an array, modifying a slice modifies the underlying array
  - slice literal can be created like arrays but without the length, internally creates an array and then a slice referencing it
  - can have an inline struct type declaration for arrays and slices:
  ```
  s := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
	}
  ```
  - high and low can be omitted and use defaults, low: `0` high: 'length` eg `a[:4]` is `{0, 1, 2, 3}`
  - slices have *length* `len(s)` and *capcity* `cap(s)`
  - length is the number of elements in a slice
  - capacity is the length of underlying array from the slices first element
  - once the slice goes beyond the first element of the array, capacity is reduced (can't be restored?)
  - the zero value of a slice is `nil`, *cap* and *len* are 0, and has no array
  - the *make* function allocates a zeroed array and returns a slice referring to it (it is how dyanmic sized arrays are created)
    - `make([]int, 4)` specifies a length, and creates a base array of 4 zeros and a slice of len 4 and cap 4
	- `make([]int, 0, 4)` optional third arg specifies a capacity, so base array `[0 0 0 0]` but empty slice `[]`, len 0 cap 4
  - slice can contain any type even other slices
  - *append* allows us to add elements to a slice and to create a new underlying array if it necessary to hold new elements
    - `append(s []T, vs ...T) []T` returns a new slice, taking mulitple params to append
- The for loop has *range* format that iterates slices (and maps), it returns two values the index and a copy of the actual value
```
var a = []int{1, 2, 3, 4}
for i, v := range a { }
```
  - if only want the index value, leave value out `for i := range a {}`
  - to omit the index use an underscore, `for _, v : range a {}`
- **Maps** contain key value pairs
  - zero value is `nil`, a `nil` map has no keys and neither can any be added
  - can use `make` to create a map `m := make(map[int]string)` add with `m[2] = "two"`
  - `map` literals are like structs, but keys are required:
  ```
  var m = map[string]uint8{
	"one": 1,
	"two": 2,
  }
  ```
  - if top-level type is just a type name, it can be omitted from literal elements:
  ```
  type Vertex struct { Lat, Long float64 }
  var m = map[string]Vertex{
	"one": {30.1, 39.1},
	"two": {32.2, 49.1},
  }
  ```
  - get and set elements; `m[key] = elem` and `elem = m[key]`
  - remove an element; `delete(m, key)`
  - test key exists; `elem, ok := m[key]` *ok* boolean if it exists, *elem* is zero value if doesn't exist
- *Functions* are values too, and can be passed around as such e.g. as args and return values
  -  give full function signature
  ```
  func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
  }
  ```
  - functions can be closures, in that they can be bound to a variable outside its body
  - each closure is bound to its own copy of any referenced variables, and can mutate them independantly

##  Methods
- go has no classes as its functional
- *methods* are functions that have a *receiver* argument, allowing them to be defined on types
```
func (v Vertex) Abs() float64 {}
var v Vertex
v.Abs()
```
- this works the same as if *v* was an argumen of `Abs(v Vertex)` and called as `Abs(v)`
- can only declare *methods* on types that are defined in the same package
- method receivers can be pointers, allowing the receiver to be modified as it is a reference to the original
- regular receivers or function arguments are passed by value, a copy of original
- methods with pointer receivers can be called with a value or a pointer (values are converted to pointers)
- for functions with pointer args only pointer will be accepted `&v`
- the reverse is also true, value receivers can be pointer or value `*p`, functions must be as declared
- benefits of a pointer receiver
  1. the method can modify the value
  - can avoid copying the value on each method call
- rule of thumb all methods on a given type should be all pointers or all values, not mixed

## Interfaces
- an `interface` type is a set of method signatures
```
type Abser interface {
	Abs() float64
}
```
- any type that has a method matching all the signatures can be of the interface type
- interfaces are implemented implicitly, a type only needs to implement the required methods
- *pointer types* will implement methods separately to value types of the same type
- an interface value `var i I` is like a tuple `(value, type)`, holding a value of a specific concrete type
- calling a method on an interface type executes the method of the same name on its underlying type
- if the concrete type of an interface value is *nil* than the method receiver is also *nil*, the interface value itself is not *nil* however
- if the interface is *nil* it has no value or type and will cause an error if used to call methods
- the empty interface `interface{}` has no methods, and can hold values of any type
- code that can hanlde values of any type uses the empty interface to specify them `var a interface{}`
- **type assertion** provide access to an interfaces underlying type
  - `t := i.(T)` asserts that interface value *i* has concrete type *T* and assigns the *T* values to *t*
  - it *i* is not a *T* this will caues a panic
  - can be tested with `t, ok := i.(T)` just like acessing a *map*
  - if *ok is false then *t* gets the zero value of *T*
- **type switch** is like a regular *switch* statement, except that it works on type assertions
  - the case statements are actual types which are compared to the underlying type of an interface value
  - the switch value has the same syntax as a type assertion except *type* is used instead of an actual type
  ```
  switch v := i.(type) {
  case int:
	// v is an int
  case string:
    // v is a string
  default:
    // don't know, v is same type as i
  }
  ```
- the **Stringer** interface is defined by the *fmt* package, it is a type that can describe itself as a string
```
type Stringer interface {
	String() string
}
```
- the **error** interface is used to express an error state
```
type error interface {
	Error() string
}
```
  - functions oftern return an *error* value, and callers should test this for *nil*
  - a nil *error* denotes success, and non-nil is failure
  - calling a *print* function (that doesn't have a format arg) on a value implementing *error* from the `Error()` method will cause an infinite loop, because go will call `Error()` on it again to turn it to a string (even if *Stringer* is also implemented)
- the **Reader** interface is defined int the *io* package, and the standard library contains many implementations
  - it has the following read method `func (T) Read(b []byte) (n int, err error)`
  - it fills the byte slice with data and returns number of bytes read and `io.EOF` error if the stream ends
  - it is common for an `io.Reader` to wrap another `io.Reader` to modify the stream

## Concurrency

- **goroutine** is a light thread managed by th Go runtime
  - `go f(a, b)` *f*, *a*, *b* are evaluated in current goroutine, execution happens in a new one
  - they will run in the same shared memory space, use `sync` package to synchronize
- **channels** are a typed conduit which you can send and receive values with the channel operator `<-`
  - `ch <- v` send v to channel ch, `v := <-ch` receive from ch and assign to v
  - *data flows in the direction of the arrow*
  - channels need to be created before first use `ch := make(chan int)`
  - by default sends and receives block until the other side is ready (no need for locks)
  - channels can be *buffered* by adding length to make `ch := make(chan int, 100)`
  - sends to buffered channel only block when full, and receives block when empty
  - a sender can `close` a channel to indicate no more values will be sent
  - receivers can test if a channel is closed `v, ok := <-ch`
  - the `range` operator can loop on a channel, receiving values until it is closed
  - *only the sender should close a channel, not receiver - sending on closed will panic*
  - channels don't ordinarily need to be closed, only in circumstances like *range*
- **select** statement allows a goroutine to wait on multiple communications
  - it will block until one of its cases can run, if multiple chooses one at random
  - the default case is run if no other case is ready (optional)
- mutual exclusion can be implemented using `sync.Mutex` with *Lock* and *Unlock*
  - surround a block of code with calls to `Lock()` and `Unlock()` to allow single access
  - `defer` can be used to ensure a mutex will be unlocked
