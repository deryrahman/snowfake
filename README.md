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

WIP

## Contribute

WIP
