# sunnygo
A go pkg to read from SMA Sunnyboy inverters

## Example
```go
package main

import (
	"io/ioutil"

	"github.com/jonkerj/sunnygo/pkg"
)

func main() {
	meta := &pkg.Meta{}

	b, err := ioutil.ReadFile("/tmp/ObjectMetadata_Istl.json")
	if err != nil {
		panic(err)
	}

	err = meta.LoadModel(b)
	if err != nil {
		panic(err)
	}

	b, err = ioutil.ReadFile("/tmp/en-US.json")
	if err != nil {
		panic(err)
	}

	err = meta.LoadLanguage(b)
	if err != nil {
		panic(err)
	}

	b, err = ioutil.ReadFile("/tmp/getAllOnValues.json")
	if err != nil {
		panic(err)
	}

	vr, err := pkg.LoadValues(b)
	if err != nil {
		panic(err)
	}

	root, err := meta.NodifyAllValues("0199-12345678", vr)
	if err != nil {
		panic(err)
	}

	pkg.PrintTree(root, 0)
}
```
Result:
```
root
  AC Side
    Grid measurements
      Phase total Current: 9.248000 A
      EEI displacement power factor: 0.000000
      Grid frequency: 49.990000 Hz
[..]
```
