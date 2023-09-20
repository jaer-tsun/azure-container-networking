package netio

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

type getInterfaceValidationFn func(name string) (*net.Interface, error)

type MockNetIO struct {
	fail           bool
	failAttempt    int
	numTimesCalled int
	getInterfaceFn getInterfaceValidationFn
}

// ErrMockNetIOFail - mock netio error
var (
	ErrMockNetIOFail = errors.New("netio fail")
	hwAddr, _        = net.ParseMAC("ab:cd:ef:12:34:56")
)

func NewMockNetIO(fail bool, failAttempt int) *MockNetIO {
	return &MockNetIO{
		fail:        fail,
		failAttempt: failAttempt,
	}
}

func (netshim *MockNetIO) SetGetInterfaceValidatonFn(fn getInterfaceValidationFn) {
	netshim.getInterfaceFn = fn
}

func (netshim *MockNetIO) GetNetworkInterfaceByName(name string) (*net.Interface, error) {
	netshim.numTimesCalled++

	if netshim.fail && netshim.failAttempt == netshim.numTimesCalled {
		return nil, fmt.Errorf("%w:%s", ErrMockNetIOFail, name)
	}

	if netshim.getInterfaceFn != nil {
		return netshim.getInterfaceFn(name)
	}

	return &net.Interface{
		//nolint:gomnd // Dummy MTU
		MTU:          1000,
		Name:         name,
		HardwareAddr: hwAddr,
		//nolint:gomnd // Dummy interface index
		Index: 2,
	}, nil
}

func (netshim *MockNetIO) GetNetworkInterfaceAddrs(iface *net.Interface) ([]net.Addr, error) {
	return []net.Addr{}, nil
}

func (netshim *MockNetIO) GetNetworkInterfaceByMac(mac net.HardwareAddr) (*net.Interface, error) {
	if !bytes.Equal(mac, hwAddr) {
		return nil, fmt.Errorf("%w: %s", ErrMockNetIOFail, mac)
	}

	return &net.Interface{
		//nolint:gomnd // Dummy MTU
		MTU:          1000,
		Name:         "eth1",
		HardwareAddr: mac,
		//nolint:gomnd // Dummy interface index
		Index: 2,
	}, nil
}
