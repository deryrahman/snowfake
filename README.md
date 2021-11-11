# Snowfake
[![Build Status](https://travis-ci.org/deryrahman/snowfake.svg?branch=master)](https://travis-ci.org/deryrahman/snowfake) [![GoDoc](https://godoc.org/github.com/deryrahman/snowfake?status.svg)](https://godoc.org/github.com/deryrahman/snowfake) [![Coverage](http://gocover.io/_badge/github.com/deryrahman/snowfake)](https://gocover.io/github.com/deryrahman/snowfake) [![Go Report Card](https://goreportcard.com/badge/github.com/deryrahman/snowfake)](https://goreportcard.com/report/github.com/deryrahman/snowfake)

Snowfake is just a Twitter Snowflake IDs alternative for generating unique short ID at high scale.

## Overview
The objective of Snowfake is to generate smaller ID compare with the original Snowflake. It's suitable for generating short IDs at scale. Case study: short link eg. bit.ly  

Snowfake also guarantees lifetime 136 years (original Snowflake 69 years). By default all IDs will be ran out in 2159. (Epoch is configurable).  

By default, a Snowfake ID is composed of

```markdown
32 bits for time in units of second
5 bits for node/machine numbers (configurable)
27 bits for sequence numbers (configurable)
```

With this configuration, Snowfake allows you to generate `2^27` unique IDs per second per machine. It also provides `2^5` distributed machines works together without additional config.  

Both machine numbers and sequence numbers are configurable. To achieve smaller ID you can minimize bits allocation by reducing node bits and sequence bits.

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
	sf, _ := snowfake.New(nodeID)

	sfID := sf.GenerateID()
	println(snowfake.EncodeBase58(sfID)) // encoding base58
}
```

To generate smaller ID, tune the node bits and seq bits numbers as small as possible

```go
package main

import (
	"github.com/deryrahman/snowfake"
)

func main() {
	snowfake.SetEpoch(1577836800) // timestamp start from 01/01/2020 @ 12:00am (UTC)
	snowfake.SetNodeBits(1)       // reserve 2^1 machine numbers
	snowfake.SetSeqBits(4)        // ~2^4 unique ID per second
	_ = snowfake.Init()           // must be called to instantiate new config

	nodeID := uint64(1)
	sf, _ := snowfake.New(nodeID)

	sfID := sf.GenerateID()
	println(snowfake.EncodeBase58(sfID)) // encoding base58
}
```

## Contribute

Just make a pull request with a clear description ☕️
