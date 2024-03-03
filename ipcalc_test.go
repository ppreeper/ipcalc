package ipcalc_test

import (
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
		bits := ipcalc.NetmaskStringToBits(mt.NetmaskString)
		if bits != mt.NetmaskBits {
			t.Errorf("expected %d, got %d", mt.NetmaskBits, bits)
		}
	}
}

func TestNetmaskToBits(t *testing.T) {
	for _, mt := range NetmaskTable {
		bits := ipcalc.NetmaskToBits(mt.NetmaskIP)
		if bits != mt.NetmaskBits {
			t.Errorf("expected %d, got %d", mt.NetmaskBits, bits)
		}
	}
}

func TestCIDRNetmask(t *testing.T) {
	for _, mt := range NetmaskTable {
		mask := ipcalc.CIDRNetmask(mt.NetmaskBits)
		if mask != mt.NetmaskIP {
			t.Errorf("expected %v, got %v", mt.NetmaskIP, mask)
		}
	}
}

func TestWildcardMask(t *testing.T) {
	for _, mt := range NetmaskTable {
		mask := ipcalc.WildcardMask(mt.NetmaskBits)
		if mask != mt.WildcardMask {
			t.Errorf("expected %v, got %v", mt.WildcardMask, mask)
		}
	}
}

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
		nets := ipcalc.MaximumSubnets(mt.NetmaskBits)
		if nets != mt.SubnetCount {
			t.Errorf("expected %v, got %v", mt.SubnetCount, nets)
		}
	}
}

func TestMaximumAddresses(t *testing.T) {
	for _, mt := range SubnetCountTable {
		nets := ipcalc.MaximumAddresses(mt.NetmaskBits)
		if nets != mt.AddrCount {
			t.Errorf("expected %v, got %v", mt.AddrCount, nets)
		}
	}
}

func TestAddrToBinary(t *testing.T) {
	ipcalc.AddrToBinary(IPAddressByte)
}

func TestBinaryToAddr(t *testing.T) {
	ipcalc.BinaryToAddr(IPAddressUint32)
}

func TestCIDRAddress(t *testing.T) {
	ipcalc.CIDRAddress(IPAddress, NetmaskBits)
}

func TestCIDRAddressFromString(t *testing.T) {
	ipcalc.CIDRAddressFromString(CIDRAddress)
}

// Benchmarks

var (
	IPAddress           = "10.16.1.1"
	IPAddressUint32     = uint32(168427521)
	IPAddressByte       = [4]byte{10, 16, 1, 1}
	NetmaskString       = "255.255.255.0"
	NetmaskIP, _        = netip.ParseAddr("255.255.255.0")
	NetmaskBits     int = 24
	CIDRAddress         = "10.16.1.1/24"
)

func BenchmarkNetmaskStringToBits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipcalc.NetmaskStringToBits(NetmaskString)
	}
}

func BenchmarkNetmaskToBits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipcalc.NetmaskToBits(NetmaskIP)
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

// MaximumSubnets
func BenchmarkMaximumSubnets(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipcalc.MaximumSubnets(NetmaskBits)
	}
}

// MaximumAddresses
func BenchmarkMaximumAddresses(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipcalc.MaximumAddresses(NetmaskBits)
	}
}

// AddrToBinary
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
		ipcalc.CIDRAddress(IPAddress, NetmaskBits)
	}
}

func BenchmarkCIDRAddressFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipcalc.CIDRAddressFromString(CIDRAddress)
	}
}
