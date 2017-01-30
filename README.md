# ERNEST PROVIDER

# Please avoid using this software as it is under development

master : [![CircleCI](https://circleci.com/gh/ernestio/ernestconnector/tree/master.svg?style=svg)](https://circleci.com/gh/ernestio/ernestprovider/tree/master) | develop : [![CircleCI](https://circleci.com/gh/ernestio/ernestprovider/tree/develop.svg?style=svg)](https://circleci.com/gh/ernestio/ernestprovider/tree/develop)

This library aims to be a wrapper on top of provider go sdk, so it concentrates all provider specific logic on ernest.

Example:
```go
package main

import(
  "fmt"

	"github.com/ernestio/ernestprovider"
	"github.com/ernestio/ernestprovider/network"
)

func main() {
	event := network.New("network.create.provider", "{....}")

	subject, data := ernestprovider.Handle(&event)
	fmt.Println("Response: ")
	fmt.Println(subject)
	fmt.Println(data)
}
```

## Using it

You can start by importing


## Contributing

Please read through our
[contributing guidelines](CONTRIBUTING.md).
Included are directions for opening issues, coding standards, and notes on
development.

Moreover, if your pull request contains patches or features, you must include
relevant unit tests.

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2015 r3labs.io authors.

Code released under
[the Mozilla Public License Version 2.0](LICENSE).

