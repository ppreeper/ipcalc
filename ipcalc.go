package ipcalc

import (
	"encoding/binary"
	"fmt"
	"net/netip"
)

type CIDR struct {
	Address          netip.Addr
	Bits             int
	Netmask          netip.Addr
	WildcardMask     netip.Addr
	MaximumSubnets   int
	MaximumAddresses int
	NetworkAddress   netip.Addr
	BroadcastAddress netip.Addr
}

func NetmaskStringToBits(mask string) int {
	maskAddr, err := netip.ParseAddr(mask)
	if err != nil {
		fmt.Println(err)
	}
	return NetmaskToBits(maskAddr)
}

func NetmaskToBits(mask netip.Addr) int {
	maskBytes := mask.As4()
	var bits int
	for _, maskByte := range maskBytes {
		for maskByte > 0 {
			if maskByte&1 > 0 {
				bits++
			}
			maskByte >>= 1
		}
	}
	return bits
}

func CIDRNetmask(maskBits int) netip.Addr {
	mask := 0xffffffff << (32 - maskBits)
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(mask))
	bmask := [4]byte{bytes[0], bytes[1], bytes[2], bytes[3]}
	return netip.AddrFrom4(bmask)
}

func WildcardMask(maskBits int) netip.Addr {
	mask := 0xffffffff >> maskBits
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(mask))
	bmask := [4]byte{bytes[0], bytes[1], bytes[2], bytes[3]}
	return netip.AddrFrom4(bmask)
}

func MaximumSubnets(maskBits int) int {
	return 1 << (32 - maskBits)
}

func MaximumAddresses(maskBits int) int {
	if maskBits >= 31 {
		return MaximumSubnets(maskBits)
	} else {
		return MaximumSubnets(maskBits) - 2
	}
}

func AddrToBinary(addr [4]byte) uint32 {
	return binary.BigEndian.Uint32(addr[:])
}

func BinaryToAddr(addr uint32) netip.Addr {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(addr))
	bmask := [4]byte{bytes[0], bytes[1], bytes[2], bytes[3]}
	return netip.AddrFrom4(bmask)
}

func CIDRAddress(addr string, bits int) CIDR {
	cidrAddr, _ := netip.ParseAddr(addr)
	cidrBits := bits
	netMask := CIDRNetmask(cidrBits)

	// cidr calculation
	cidrAddrBinary := AddrToBinary(cidrAddr.As4())
	netMaskBinary := AddrToBinary(netMask.As4())
	cidrNetworkBinary := cidrAddrBinary & netMaskBinary
	cidrBroadcastBinary := cidrNetworkBinary | ^netMaskBinary

	return CIDR{
		Address:          cidrAddr,
		Bits:             cidrBits,
		Netmask:          netMask,
		WildcardMask:     WildcardMask(cidrBits),
		MaximumSubnets:   MaximumSubnets(cidrBits),
		MaximumAddresses: MaximumAddresses(cidrBits),
		NetworkAddress:   BinaryToAddr(cidrNetworkBinary),
		BroadcastAddress: BinaryToAddr(cidrBroadcastBinary),
	}
}

func CIDRAddressFromString(cidr string) CIDR {
	pfx, _ := netip.ParsePrefix(cidr)

	cidrAddr := pfx.Addr()
	cidrBits := pfx.Bits()
	netMask := CIDRNetmask(cidrBits)

	// cidr calculation
	cidrAddrBinary := AddrToBinary(cidrAddr.As4())
	netMaskBinary := AddrToBinary(netMask.As4())
	cidrNetworkBinary := cidrAddrBinary & netMaskBinary
	cidrBroadcastBinary := cidrNetworkBinary | ^netMaskBinary

	return CIDR{
		Address:          cidrAddr,
		Bits:             cidrBits,
		Netmask:          netMask,
		WildcardMask:     WildcardMask(cidrBits),
		MaximumSubnets:   MaximumSubnets(cidrBits),
		MaximumAddresses: MaximumAddresses(cidrBits),
		NetworkAddress:   BinaryToAddr(cidrNetworkBinary),
		BroadcastAddress: BinaryToAddr(cidrBroadcastBinary),
	}
}
