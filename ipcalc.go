// Package ipcalc provides IPv4 CIDR calculation utilities including
// netmask conversion, network/broadcast address derivation, and
// host count computation.
package ipcalc

import (
	"encoding/binary"
	"fmt"
	"math/bits"
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

// NetmaskStringToBits parses an IPv4 netmask string (e.g. "255.255.255.0")
// and returns the number of set bits in the mask.
// Returns an error if the string is not a valid IPv4 address or the mask is
// not contiguous.
func NetmaskStringToBits(mask string) (int, error) {
	maskAddr, err := netip.ParseAddr(mask)
	if err != nil {
		return 0, fmt.Errorf("NetmaskStringToBits: %w", err)
	}
	n, err := NetmaskToBits(maskAddr)
	if err != nil {
		return 0, fmt.Errorf("NetmaskStringToBits: %w", err)
	}
	return n, nil
}

// NetmaskToBits returns the number of set bits in an IPv4 netmask address.
// Returns an error if the address is not an IPv4 address or if the mask is
// not a contiguous sequence of set bits (for example 255.0.255.0 is rejected).
func NetmaskToBits(mask netip.Addr) (int, error) {
	if !mask.Unmap().Is4() {
		return 0, fmt.Errorf("IPv6 not supported")
	}
	b := mask.Unmap().As4()
	n := bits.OnesCount8(b[0]) + bits.OnesCount8(b[1]) +
		bits.OnesCount8(b[2]) + bits.OnesCount8(b[3])
	if AddrToBinary(b) != AddrToBinary(CIDRNetmask(n).As4()) {
		return 0, fmt.Errorf("%s is not a contiguous netmask", mask)
	}
	return n, nil
}

// clampBits constrains a prefix length to the valid IPv4 range [0, 32].
func clampBits(maskBits int) int {
	if maskBits < 0 {
		return 0
	}
	if maskBits > 32 {
		return 32
	}
	return maskBits
}

// CIDRNetmask returns the subnet mask for the given IPv4 prefix length as a
// netip.Addr. Prefix lengths outside [0, 32] are clamped to that range.
func CIDRNetmask(maskBits int) netip.Addr {
	maskBits = clampBits(maskBits)
	// Shifting a uint32 by 32 is well-defined in Go and correctly produces 0 for /0.
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], ^uint32(0)<<(32-uint(maskBits)))
	return netip.AddrFrom4(b)
}

// WildcardMask returns the wildcard (inverse) mask for the given IPv4 prefix
// length as a netip.Addr. Prefix lengths outside [0, 32] are clamped to that range.
func WildcardMask(maskBits int) netip.Addr {
	maskBits = clampBits(maskBits)
	// Shifting a uint32 by 32 is well-defined in Go and correctly produces 0xffffffff for /0.
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], ^uint32(0)>>uint(maskBits))
	return netip.AddrFrom4(b)
}

// MaximumSubnets returns the total number of addresses in a subnet with the
// given prefix length, including network and broadcast addresses. Prefix
// lengths outside [0, 32] are clamped to that range. Note that /0 yields 2^32,
// which does not fit in an int on 32-bit platforms.
func MaximumSubnets(maskBits int) int {
	maskBits = clampBits(maskBits)
	return 1 << (32 - maskBits)
}

// MaximumAddresses returns the number of usable host addresses in a subnet
// with the given prefix length. For /31 and /32 all addresses are usable; for
// shorter prefixes network and broadcast addresses are excluded. Prefix
// lengths outside [0, 32] are clamped to that range.
func MaximumAddresses(maskBits int) int {
	maskBits = clampBits(maskBits)
	if maskBits >= 31 {
		return MaximumSubnets(maskBits)
	}
	return MaximumSubnets(maskBits) - 2
}

// AddrToBinary converts a 4-byte IPv4 address to its big-endian uint32 representation.
func AddrToBinary(addr [4]byte) uint32 {
	return binary.BigEndian.Uint32(addr[:])
}

// BinaryToAddr converts a big-endian uint32 to a netip.Addr IPv4 address.
func BinaryToAddr(addr uint32) netip.Addr {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], addr)
	return netip.AddrFrom4(b)
}

// buildCIDR constructs a CIDR value from a parsed IPv4 address and prefix length.
// Preconditions: addr must be an IPv4 address, bits must be in [0, 32].
func buildCIDR(addr netip.Addr, cidrBits int) CIDR {
	netMask := CIDRNetmask(cidrBits)
	addrBin := AddrToBinary(addr.As4())
	maskBin := AddrToBinary(netMask.As4())
	netBin := addrBin & maskBin
	bcBin := netBin | ^maskBin
	return CIDR{
		Address:          addr,
		Bits:             cidrBits,
		Netmask:          netMask,
		WildcardMask:     WildcardMask(cidrBits),
		MaximumSubnets:   MaximumSubnets(cidrBits),
		MaximumAddresses: MaximumAddresses(cidrBits),
		NetworkAddress:   BinaryToAddr(netBin),
		BroadcastAddress: BinaryToAddr(bcBin),
	}
}

// String returns a human-readable summary of the CIDR block.
func (c CIDR) String() string {
	return fmt.Sprintf(
		"Address: %s/%d  Network: %s  Broadcast: %s  Netmask: %s  Wildcard: %s  Hosts: %d",
		c.Address, c.Bits, c.NetworkAddress, c.BroadcastAddress,
		c.Netmask, c.WildcardMask, c.MaximumAddresses,
	)
}

// validateIPv4Bits returns an error if addr is not an IPv4 address
// or if bits is outside the valid prefix length range [0, 32].
func validateIPv4Bits(addr netip.Addr, bits int) error {
	if !addr.Unmap().Is4() {
		return fmt.Errorf("IPv6 not supported")
	}
	if bits < 0 || bits > 32 {
		return fmt.Errorf("prefix length %d out of range [0, 32]", bits)
	}
	return nil
}

// CIDRAddress parses an IPv4 address string and prefix length, returning a populated CIDR.
// Returns an error if the address is invalid, not IPv4, or bits is outside [0, 32].
func CIDRAddress(addr string, bits int) (CIDR, error) {
	cidrAddr, err := netip.ParseAddr(addr)
	if err != nil {
		return CIDR{}, fmt.Errorf("CIDRAddress: %w", err)
	}
	if err := validateIPv4Bits(cidrAddr, bits); err != nil {
		return CIDR{}, fmt.Errorf("CIDRAddress: %w", err)
	}
	return buildCIDR(cidrAddr, bits), nil
}

// CIDRAddressFromString parses a CIDR notation string (e.g. "192.168.1.0/24")
// and returns a populated CIDR.
// Returns an error if the string is invalid or the address is not IPv4.
func CIDRAddressFromString(cidr string) (CIDR, error) {
	pfx, err := netip.ParsePrefix(cidr)
	if err != nil {
		return CIDR{}, fmt.Errorf("CIDRAddressFromString: %w", err)
	}
	addr := pfx.Addr()
	bits := pfx.Bits()
	if err := validateIPv4Bits(addr, bits); err != nil {
		return CIDR{}, fmt.Errorf("CIDRAddressFromString: %w", err)
	}
	return buildCIDR(addr, bits), nil
}
