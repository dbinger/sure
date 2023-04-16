# The sure package

The sure package provides a struct that wraps `*testing.T` and 
provides only 2 basic assertion methods:
* `Same` 
* `Diff`

Comparisons are performed using `cmp.Equal`, and error messages
include appropriate diffs.  If the arguments can't be compared using
`cmp.Equal`, the test fails.

### **Objectives**

 1. Make go tests easier to read.
 2. Make go tests easier to write.
 3. Make test error messages easier to interpret.

### **Example Usage**

```
import "github.com/dbinger/sure"

func TestFun(t *testing.T) {
    b := sure.Be(t)
    result, err := Fun("1")
    b.Same(err, nil)
    b.Diff(result, 1)
}
```

### Error Messages

All error messages include file name, line number, and test name.

Error messages are constructed on the assumption that the first argument
is the "got".  

When one of the arguments contains a struct with an unexported field,
the error message will contain the error string produced by the
resulting `cmp.Equal` panic call.  See the documentation in the `cmp`
package for comparison options.  Comparison options can be provided
as additional arguments to the `sure.Be` constructor.

By default, the `cmpoption.EquateErrors()` is included in the `Be` constructor.  This means that comparisons will treat errors as the same whenever `errors.Is` would return true.

If you want error messsages to include additional 
notes about the error, you can add them as an additional 
arguments to the assertion method call like this:

```
b.Diff(err, nil, "could not find user", username)
```

Note that the `cmp.Equal` function returns false when comparing nil to a nil pointer.  
If you want to use these methods to check if a pointer like *X is nil, 
you need to compare it to (*X)(nil), not plain nil.

***


### Example Error Messages


***

### Example 1

Failing Test
```
be := sure.Be(t)

be.Same(nil, "hello world")
```

Error Message
```
example_test.go:23: FAIL in TestExample 
GOT:  nil
WANT: string("hello world")
```

***

### Example 2

Failing Test
```
be := sure.Be(t)

be.Same(42, nil, "the answer")
```

Error Message
```
example_test.go:23: FAIL in TestExample 
the answer
GOT:  int(42)
WANT: nil
```

***

### Example 3

Failing Test
```
be := sure.Be(t)
err := errors.New("ouch")

be.Same(err, nil)
```

Error Message
```
example_test.go:23: FAIL in TestExample 
GOT:  e"ouch"
WANT: nil
```

***

### Example 4

Failing Test
```
be := sure.Be(t)
tls1 := tls.Config{InsecureSkipVerify: true, ServerName: "alpha"}
tls2 := tls.Config{InsecureSkipVerify: true, ServerName: "beta", MaxVersion: 3}

be.Same(tls1, tls2)
```

Error Message
```
example_test.go:23: FAIL in TestExample 
error: compare: cannot handle unexported field at {tls.Config}.mutex:
	"crypto/tls".Config
consider using a custom Comparer; if you control the implementation of type, you can also consider using an Exporter, AllowUnexported, or cmpopts.IgnoreUnexported
```


This comparison didn't work because `cmp.Equal` does not have a default behavior for comparing structs with unexported fields. 

The next example shows how this can be fixed by using `cmpoption.IgnoreUnexported`.

***

### Example 5

Failing Test
```
be := sure.Be(t, cmpopts.IgnoreUnexported(tls.Config{}))
tls1 := tls.Config{InsecureSkipVerify: true, ServerName: "alpha"}
tls2 := tls.Config{InsecureSkipVerify: true, ServerName: "beta", MaxVersion: 3}

be.Same(tls1, tls2)
```

Error Message
```
example_test.go:23: FAIL in TestExample 
tls.Config{
      ... // 9 identical fields
      RootCAs:    nil,
      NextProtos: nil,
GOT:  ServerName: "alpha",
WANT: ServerName: "beta",
      ClientAuth: s"NoClientCert",
      ClientCAs:  nil,
      ... // 5 identical fields
      ClientSessionCache:          nil,
      MinVersion:                  0,
GOT:  MaxVersion:                  0,
WANT: MaxVersion:                  3,
      CurvePreferences:            nil,
      DynamicRecordSizingDisabled: false,
      ... // 3 ignored and 2 identical fields
  }
```

***

### Example 6

Failing Test
```
be := sure.Be(t)
map1    = map[string]int{"A": 1, "B": 2}
map2    = map[string]int{"A": 1, "B": 2, "C": 3}

be.Same(map1, map2)
```

Error Message
```
example_test.go:23: FAIL in TestExample 
map[string]int{
      "A": 1,
      "B": 2,
WANT: "C": 3,
  }
```

***

### Example 7

Failing Test
```
be := sure.Be(t)
map1    = map[string]int{"A": 1, "B": 2, "C": 3}
map2    = map[string]int{"A": 1, "B": 2}


be.Same(map1, map2)
```

Error Message
```
example_test.go:23: FAIL in TestExample 
map[string]int{
      "A": 1,
      "B": 2,
GOT:  "C": 3,
  }
```

***

### Example 8

Failing Test
```
be := sure.Be(t)
map1    = map[string]int{"A": 1, "B": 2: "D": 4, "E": 5}
map2    = map[string]int{"A": 1, "B": 2, "C": 3}

be.Same(map1, map2)
```

Error Message
```
example_test.go:23: FAIL in TestExample 
map[string]int{
      "A": 1,
      "B": 2,
WANT: "C": 3,
GOT:  "D": 4,
GOT:  "E": 5,
  }
```

***

### Example 9

Failing Test
```
be := sure.Be(t)
err := errors.New("ouch")

be.Diff(err, err)
```

Error Message
```
example_test.go:23: FAIL in TestExample 
GOT:  e"ouch"
WANT: anything else
```

***

### Example 10

Failing Test
```
be := sure.Be(t)

be.Diff(42, 42, "answers")
```

Error Message
```
example_test.go:23: FAIL in TestExample 
answers
GOT:  int(42)
WANT: anything else
```
