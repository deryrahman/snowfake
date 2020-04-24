# Snowfake
Snowfake is just a Twitter Snowflake IDs alternative for generating unique short ID at high scale.
The objective of Snowfake is to generate smaller ID compare with the original Snowflake.
It's suitable for generating short IDs at scale.

## Overview
By default, a Snowfake ID is composed of

```markdown
32 bits for time in units of second
8 bits for node/machine numbers (configurable)
24 bits for sequence numbers (configurable)
```

With this configuration, Snowfake guarantees to generate `16777216` unique IDs per second per machine.
To achieve smaller ID you can minimize bits allocation by reducing node bits and sequence bits. 

## Installation

```shell script
go get github.com/deryrahman/snowfake
```

## Usage

```go
package main

import (
	"github.com/deryrahman/snowfake"
)

func main() {
	nodeID := uint64(1)
	sf := snowfake.New(nodeID) // instantiate with default config

	sfID := sf.GenerateID()
	println(snowfake.Encode(sfID)) // encoding base58
}
```

To generate smaller ID, tune the node bits and step bits numbers as small as possible

```go
package main

import (
	"github.com/deryrahman/snowfake"
)

func main() {
	nodeID := uint64(1)
	epoch := uint64(1577836800) // timestamp start from 01/01/2020 @ 12:00am (UTC)
	nodeBits := uint8(1)        // reserve 2^1 machine numbers
	stepBits := uint8(4)        // ~2^4 unique ID per second

	sf, _ := snowfake.NewWithConfig(nodeID, epoch, nodeBits, stepBits) // instantiate with custom config

	sfID := sf.GenerateID()
	println(snowfake.Encode(sfID)) // encoding base58
}
```

## Contribute

Just make a pull request with a clear description ☕️
