# ipcalc

IPv4 network calculator: netmask conversion, network/broadcast address
derivation, and host count computation.

## Install

```sh
go get github.com/ppreeper/ipcalc
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/ppreeper/ipcalc"
)

func main() {
	cidr, err := ipcalc.CIDRAddressFromString("10.16.1.1/24")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cidr.String())
	// Address: 10.16.1.1/24  Network: 10.16.1.0  Broadcast: 10.16.1.255  Netmask: 255.255.255.0  Wildcard: 0.0.0.255  Hosts: 254

	bits, err := ipcalc.NetmaskStringToBits("255.255.255.0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bits) // 24
}
```

## Notes

- Only IPv4 is supported; pure IPv6 input is rejected. IPv4-mapped IPv6 input
  (e.g. `::ffff:10.0.0.1`) is accepted and normalized to IPv4.
- Netmasks must be contiguous (e.g. `255.0.255.0` is rejected).
- Prefix lengths passed to the low-level helpers (`CIDRNetmask`, `WildcardMask`,
  `MaximumSubnets`, `MaximumAddresses`) are clamped to `[0, 32]`.
