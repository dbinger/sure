# The sure package

The sure package provides a struct that wraps `*testing.T` and 
provides only 2 basic assertion methods:
* `Same` 
* `Diff`

Comparisons are performed using `cmp.Equal`, and error messages
include appropriate diffs.  If the arguments can't be compared using
`cmp.Equal`, the test fails.

### **Objectives**

 1. Make go tests a bit easier to read.
 2. Use minimal set of assertion methods.
 3. Reduce or eliminate the need for format strings in test code.
 4. Make error messages easier to interpret.

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
as additional arguments to the sure.Be constructor.

If you want error messsages to include additional 
notes about the error, you can add them as an additional 
arguments to the assertion method call like this:

```
b.Diff(err, nil, "could not find user", username)
```

# Usage Note

The `cmp.Equal` function returns false when comparing nil to a nil pointer.  
If you want to use these methods to check if a pointer like *X is nil, 
you need to compare it to (*X)(nil), not plain nil.

