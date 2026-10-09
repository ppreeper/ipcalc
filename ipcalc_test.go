package ipcalc_test

import (
	"fmt"
	"net/netip"
	"testing"

	"github.com/ppreeper/ipcalc"
)

// Tests

var NetmaskTable = []struct {
	NetmaskString string
	NetmaskIP     netip.Addr
	WildcardMask  netip.Addr
	NetmaskBits   int
}{
	{"255.255.255.255", netip.AddrFrom4([4]byte{255, 255, 255, 255}), netip.AddrFrom4([4]byte{0, 0, 0, 0}), 32},
	{"255.255.255.254", netip.AddrFrom4([4]byte{255, 255, 255, 254}), netip.AddrFrom4([4]byte{0, 0, 0, 1}), 31},
	{"255.255.255.252", netip.AddrFrom4([4]byte{255, 255, 255, 252}), netip.AddrFrom4([4]byte{0, 0, 0, 3}), 30},
	{"255.255.255.248", netip.AddrFrom4([4]byte{255, 255, 255, 248}), netip.AddrFrom4([4]byte{0, 0, 0, 7}), 29},
	{"255.255.255.240", netip.AddrFrom4([4]byte{255, 255, 255, 240}), netip.AddrFrom4([4]byte{0, 0, 0, 15}), 28},
	{"255.255.255.224", netip.AddrFrom4([4]byte{255, 255, 255, 224}), netip.AddrFrom4([4]byte{0, 0, 0, 31}), 27},
	{"255.255.255.192", netip.AddrFrom4([4]byte{255, 255, 255, 192}), netip.AddrFrom4([4]byte{0, 0, 0, 63}), 26},
	{"255.255.255.128", netip.AddrFrom4([4]byte{255, 255, 255, 128}), netip.AddrFrom4([4]byte{0, 0, 0, 127}), 25},
	{"255.255.255.0", netip.AddrFrom4([4]byte{255, 255, 255, 0}), netip.AddrFrom4([4]byte{0, 0, 0, 255}), 24},
	{"255.255.254.0", netip.AddrFrom4([4]byte{255, 255, 254, 0}), netip.AddrFrom4([4]byte{0, 0, 1, 255}), 23},
	{"255.255.252.0", netip.AddrFrom4([4]byte{255, 255, 252, 0}), netip.AddrFrom4([4]byte{0, 0, 3, 255}), 22},
	{"255.255.248.0", netip.AddrFrom4([4]byte{255, 255, 248, 0}), netip.AddrFrom4([4]byte{0, 0, 7, 255}), 21},
	{"255.255.240.0", netip.AddrFrom4([4]byte{255, 255, 240, 0}), netip.AddrFrom4([4]byte{0, 0, 15, 255}), 20},
	{"255.255.224.0", netip.AddrFrom4([4]byte{255, 255, 224, 0}), netip.AddrFrom4([4]byte{0, 0, 31, 255}), 19},
	{"255.255.192.0", netip.AddrFrom4([4]byte{255, 255, 192, 0}), netip.AddrFrom4([4]byte{0, 0, 63, 255}), 18},
	{"255.255.128.0", netip.AddrFrom4([4]byte{255, 255, 128, 0}), netip.AddrFrom4([4]byte{0, 0, 127, 255}), 17},
	{"255.255.0.0", netip.AddrFrom4([4]byte{255, 255, 0, 0}), netip.AddrFrom4([4]byte{0, 0, 255, 255}), 16},
	{"255.254.0.0", netip.AddrFrom4([4]byte{255, 254, 0, 0}), netip.AddrFrom4([4]byte{0, 1, 255, 255}), 15},
	{"255.252.0.0", netip.AddrFrom4([4]byte{255, 252, 0, 0}), netip.AddrFrom4([4]byte{0, 3, 255, 255}), 14},
	{"255.248.0.0", netip.AddrFrom4([4]byte{255, 248, 0, 0}), netip.AddrFrom4([4]byte{0, 7, 255, 255}), 13},
	{"255.240.0.0", netip.AddrFrom4([4]byte{255, 240, 0, 0}), netip.AddrFrom4([4]byte{0, 15, 255, 255}), 12},
	{"255.224.0.0", netip.AddrFrom4([4]byte{255, 224, 0, 0}), netip.AddrFrom4([4]byte{0, 31, 255, 255}), 11},
	{"255.192.0.0", netip.AddrFrom4([4]byte{255, 192, 0, 0}), netip.AddrFrom4([4]byte{0, 63, 255, 255}), 10},
	{"255.128.0.0", netip.AddrFrom4([4]byte{255, 128, 0, 0}), netip.AddrFrom4([4]byte{0, 127, 255, 255}), 9},
	{"255.0.0.0", netip.AddrFrom4([4]byte{255, 0, 0, 0}), netip.AddrFrom4([4]byte{0, 255, 255, 255}), 8},
	{"254.0.0.0", netip.AddrFrom4([4]byte{254, 0, 0, 0}), netip.AddrFrom4([4]byte{1, 255, 255, 255}), 7},
	{"252.0.0.0", netip.AddrFrom4([4]byte{252, 0, 0, 0}), netip.AddrFrom4([4]byte{3, 255, 255, 255}), 6},
	{"248.0.0.0", netip.AddrFrom4([4]byte{248, 0, 0, 0}), netip.AddrFrom4([4]byte{7, 255, 255, 255}), 5},
	{"240.0.0.0", netip.AddrFrom4([4]byte{240, 0, 0, 0}), netip.AddrFrom4([4]byte{15, 255, 255, 255}), 4},
	{"224.0.0.0", netip.AddrFrom4([4]byte{224, 0, 0, 0}), netip.AddrFrom4([4]byte{31, 255, 255, 255}), 3},
	{"192.0.0.0", netip.AddrFrom4([4]byte{192, 0, 0, 0}), netip.AddrFrom4([4]byte{63, 255, 255, 255}), 2},
	{"128.0.0.0", netip.AddrFrom4([4]byte{128, 0, 0, 0}), netip.AddrFrom4([4]byte{127, 255, 255, 255}), 1},
	{"0.0.0.0", netip.AddrFrom4([4]byte{0, 0, 0, 0}), netip.AddrFrom4([4]byte{255, 255, 255, 255}), 0},
}

func TestNetmaskStringToBits(t *testing.T) {
	for _, mt := range NetmaskTable {
		t.Run(mt.NetmaskString, func(t *testing.T) {
			bits, err := ipcalc.NetmaskStringToBits(mt.NetmaskString)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if bits != mt.NetmaskBits {
				t.Errorf("expected %d, got %d", mt.NetmaskBits, bits)
			}
		})
	}
}

func TestNetmaskStringToBitsError(t *testing.T) {
	t.Run("invalid address", func(t *testing.T) {
		_, err := ipcalc.NetmaskStringToBits("not-an-ip")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
	t.Run("IPv6 rejected", func(t *testing.T) {
		_, err := ipcalc.NetmaskStringToBits("::1")
		if err == nil {
			t.Error("expected error for IPv6, got nil")
		}
	})
	t.Run("non-contiguous netmask rejected", func(t *testing.T) {
		_, err := ipcalc.NetmaskStringToBits("255.0.255.0")
		if err == nil {
			t.Error("expected error for non-contiguous netmask, got nil")
		}
	})
}

func TestNetmaskToBits(t *testing.T) {
	for _, mt := range NetmaskTable {
		t.Run(mt.NetmaskString, func(t *testing.T) {
			bits, err := ipcalc.NetmaskToBits(mt.NetmaskIP)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if bits != mt.NetmaskBits {
				t.Errorf("expected %d, got %d", mt.NetmaskBits, bits)
			}
		})
	}
}

func TestNetmaskToBitsError(t *testing.T) {
	t.Run("IPv6 rejected", func(t *testing.T) {
		_, err := ipcalc.NetmaskToBits(netip.MustParseAddr("::1"))
		if err == nil {
			t.Error("expected error for IPv6, got nil")
		}
	})
	t.Run("non-contiguous mask rejected", func(t *testing.T) {
		for _, mask := range []string{"255.0.255.0", "255.255.255.1", "0.255.0.0"} {
			t.Run(mask, func(t *testing.T) {
				_, err := ipcalc.NetmaskToBits(netip.MustParseAddr(mask))
				if err == nil {
					t.Errorf("expected error for non-contiguous mask %s, got nil", mask)
				}
			})
		}
	})
}

func TestCIDRNetmask(t *testing.T) {
	for _, mt := range NetmaskTable {
		t.Run(mt.NetmaskString, func(t *testing.T) {
			mask := ipcalc.CIDRNetmask(mt.NetmaskBits)
			if mask != mt.NetmaskIP {
				t.Errorf("expected %v, got %v", mt.NetmaskIP, mask)
			}
		})
	}
}

func TestWildcardMask(t *testing.T) {
	for _, mt := range NetmaskTable {
		t.Run(mt.NetmaskString, func(t *testing.T) {
			mask := ipcalc.WildcardMask(mt.NetmaskBits)
			if mask != mt.WildcardMask {
				t.Errorf("expected %v, got %v", mt.WildcardMask, mask)
			}
		})
	}
}

// SubnetCountTable covers prefix lengths /1 through /32.
// /0 is intentionally omitted: MaximumSubnets(0) = 1<<32 = 4294967296,
// which overflows int on 32-bit platforms.
var SubnetCountTable = []struct {
	NetmaskBits int
	SubnetCount int
	AddrCount   int
}{
	{32, 1, 1},
	{31, 2, 2},
	{30, 4, 2},
	{29, 8, 6},
	{28, 16, 14},
	{27, 32, 30},
	{26, 64, 62},
	{25, 128, 126},
	{24, 256, 254},
	{23, 512, 510},
	{22, 1024, 1022},
	{21, 2048, 2046},
	{20, 4096, 4094},
	{19, 8192, 8190},
	{18, 16384, 16382},
	{17, 32768, 32766},
	{16, 65536, 65534},
	{15, 131072, 131070},
	{14, 262144, 262142},
	{13, 524288, 524286},
	{12, 1048576, 1048574},
	{11, 2097152, 2097150},
	{10, 4194304, 4194302},
	{9, 8388608, 8388606},
	{8, 16777216, 16777214},
	{7, 33554432, 33554430},
	{6, 67108864, 67108862},
	{5, 134217728, 134217726},
	{4, 268435456, 268435454},
	{3, 536870912, 536870910},
	{2, 1073741824, 1073741822},
	{1, 2147483648, 2147483646},
}

func TestMaximumSubnets(t *testing.T) {
	for _, mt := range SubnetCountTable {
		t.Run(fmt.Sprintf("/%d", mt.NetmaskBits), func(t *testing.T) {
			nets := ipcalc.MaximumSubnets(mt.NetmaskBits)
			if nets != mt.SubnetCount {
				t.Errorf("expected %v, got %v", mt.SubnetCount, nets)
			}
		})
	}
}

func TestMaximumAddresses(t *testing.T) {
	for _, mt := range SubnetCountTable {
		t.Run(fmt.Sprintf("/%d", mt.NetmaskBits), func(t *testing.T) {
			nets := ipcalc.MaximumAddresses(mt.NetmaskBits)
			if nets != mt.AddrCount {
				t.Errorf("expected %v, got %v", mt.AddrCount, nets)
			}
		})
	}
}

// TestBitRangeClamping verifies that out-of-range prefix lengths are clamped to
// [0, 32] instead of panicking or returning undefined results.
func TestBitRangeClamping(t *testing.T) {
	addrTests := []struct {
		name string
		got  netip.Addr
		want netip.Addr
	}{
		{"CIDRNetmask(-1)", ipcalc.CIDRNetmask(-1), ipcalc.CIDRNetmask(0)},
		{"CIDRNetmask(33)", ipcalc.CIDRNetmask(33), ipcalc.CIDRNetmask(32)},
		{"WildcardMask(-1)", ipcalc.WildcardMask(-1), ipcalc.WildcardMask(0)},
		{"WildcardMask(33)", ipcalc.WildcardMask(33), ipcalc.WildcardMask(32)},
	}
	for _, tt := range addrTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, tt.got)
			}
		})
	}

	countTests := []struct {
		name string
		got  int
		want int
	}{
		{"MaximumSubnets(-1)", ipcalc.MaximumSubnets(-1), ipcalc.MaximumSubnets(0)},
		{"MaximumSubnets(33)", ipcalc.MaximumSubnets(33), ipcalc.MaximumSubnets(32)},
		{"MaximumAddresses(-1)", ipcalc.MaximumAddresses(-1), ipcalc.MaximumAddresses(0)},
		{"MaximumAddresses(33)", ipcalc.MaximumAddresses(33), ipcalc.MaximumAddresses(32)},
	}
	for _, tt := range countTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, tt.got)
			}
		})
	}
}

func TestAddrToBinary(t *testing.T) {
	got := ipcalc.AddrToBinary(IPAddressByte)
	if got != IPAddressUint32 {
		t.Errorf("expected %d, got %d", IPAddressUint32, got)
	}
}

func TestBinaryToAddr(t *testing.T) {
	got := ipcalc.BinaryToAddr(IPAddressUint32)
	want := netip.AddrFrom4(IPAddressByte)
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestCIDRString(t *testing.T) {
	tests := []struct {
		cidr string
		want string
	}{
		{
			"10.16.1.1/24",
			"Address: 10.16.1.1/24  Network: 10.16.1.0  Broadcast: 10.16.1.255  Netmask: 255.255.255.0  Wildcard: 0.0.0.255  Hosts: 254",
		},
		{
			"10.16.1.1/31",
			"Address: 10.16.1.1/31  Network: 10.16.1.0  Broadcast: 10.16.1.1  Netmask: 255.255.255.254  Wildcard: 0.0.0.1  Hosts: 2",
		},
		{
			"10.16.1.1/32",
			"Address: 10.16.1.1/32  Network: 10.16.1.1  Broadcast: 10.16.1.1  Netmask: 255.255.255.255  Wildcard: 0.0.0.0  Hosts: 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.cidr, func(t *testing.T) {
			cidr, err := ipcalc.CIDRAddressFromString(tt.cidr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := cidr.String(); got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestCIDRAddress(t *testing.T) {
	cidr, err := ipcalc.CIDRAddress(IPAddress, NetmaskBits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cidr.Bits != NetmaskBits {
		t.Errorf("expected bits %d, got %d", NetmaskBits, cidr.Bits)
	}
	if cidr.NetworkAddress != netip.AddrFrom4([4]byte{10, 16, 1, 0}) {
		t.Errorf("unexpected network address: %v", cidr.NetworkAddress)
	}
	if cidr.BroadcastAddress != netip.AddrFrom4([4]byte{10, 16, 1, 255}) {
		t.Errorf("unexpected broadcast address: %v", cidr.BroadcastAddress)
	}
}

func TestCIDRAddressErrors(t *testing.T) {
	t.Run("invalid address", func(t *testing.T) {
		_, err := ipcalc.CIDRAddress("not-an-ip", 24)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
	t.Run("IPv6 rejected", func(t *testing.T) {
		_, err := ipcalc.CIDRAddress("::1", 24)
		if err == nil {
			t.Error("expected error for IPv6, got nil")
		}
	})
	t.Run("bits too high", func(t *testing.T) {
		_, err := ipcalc.CIDRAddress("10.0.0.1", 33)
		if err == nil {
			t.Error("expected error for bits=33, got nil")
		}
	})
	t.Run("bits negative", func(t *testing.T) {
		_, err := ipcalc.CIDRAddress("10.0.0.1", -1)
		if err == nil {
			t.Error("expected error for bits=-1, got nil")
		}
	})
}

// TestCIDRAddressIPv4Mapped documents that IPv4-mapped IPv6 input is accepted
// and that the derived network address is a plain IPv4 address.
func TestCIDRAddressIPv4Mapped(t *testing.T) {
	cidr, err := ipcalc.CIDRAddress("::ffff:10.0.0.1", 24)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cidr.NetworkAddress != netip.AddrFrom4([4]byte{10, 0, 0, 0}) {
		t.Errorf("unexpected network address: %v", cidr.NetworkAddress)
	}
	if cidr.BroadcastAddress != netip.AddrFrom4([4]byte{10, 0, 0, 255}) {
		t.Errorf("unexpected broadcast address: %v", cidr.BroadcastAddress)
	}
}

func TestCIDRAddressFromString(t *testing.T) {
	cidr, err := ipcalc.CIDRAddressFromString(CIDRAddressStr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cidr.Bits != NetmaskBits {
		t.Errorf("expected bits %d, got %d", NetmaskBits, cidr.Bits)
	}
	if cidr.NetworkAddress != netip.AddrFrom4([4]byte{10, 16, 1, 0}) {
		t.Errorf("unexpected network address: %v", cidr.NetworkAddress)
	}
	if cidr.BroadcastAddress != netip.AddrFrom4([4]byte{10, 16, 1, 255}) {
		t.Errorf("unexpected broadcast address: %v", cidr.BroadcastAddress)
	}
}

func TestCIDRAddressFromStringErrors(t *testing.T) {
	t.Run("invalid CIDR", func(t *testing.T) {
		_, err := ipcalc.CIDRAddressFromString("not-a-cidr")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
	t.Run("IPv6 rejected", func(t *testing.T) {
		_, err := ipcalc.CIDRAddressFromString("::1/64")
		if err == nil {
			t.Error("expected error for IPv6, got nil")
		}
	})
}

// Benchmarks

var (
	IPAddress       = "10.16.1.1"
	IPAddressUint32 = uint32(168820993)
	IPAddressByte   = [4]byte{10, 16, 1, 1}
	NetmaskString   = "255.255.255.0"
	NetmaskIP, _    = netip.ParseAddr("255.255.255.0")
	NetmaskBits     = 24
	CIDRAddressStr  = "10.16.1.1/24"
)

func BenchmarkNetmaskStringToBits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := ipcalc.NetmaskStringToBits(NetmaskString); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNetmaskToBits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := ipcalc.NetmaskToBits(NetmaskIP); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCIDRNetmask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipcalc.CIDRNetmask(NetmaskBits)
	}
}

func BenchmarkWildcardMask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipcalc.WildcardMask(NetmaskBits)
	}
}

func BenchmarkMaximumSubnets(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipcalc.MaximumSubnets(NetmaskBits)
	}
}

func BenchmarkMaximumAddresses(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipcalc.MaximumAddresses(NetmaskBits)
	}
}

func BenchmarkAddrToBinary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipcalc.AddrToBinary(IPAddressByte)
	}
}

func BenchmarkBinaryToAddr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipcalc.BinaryToAddr(IPAddressUint32)
	}
}

func BenchmarkCIDRAddress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := ipcalc.CIDRAddress(IPAddress, NetmaskBits); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCIDRAddressFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := ipcalc.CIDRAddressFromString(CIDRAddressStr); err != nil {
			b.Fatal(err)
		}
	}
}
