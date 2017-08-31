package device

import (
	"strings"
	"sync"
	"time"

	"github.com/AutogrowSystems/go-intelli/hid"
	"github.com/AutogrowSystems/go-intelli/util/tell"
)

const (
	// IntelliDoseDeviceName is the name that appears on Windows
	IntelliDoseDeviceName = "IntelliDose"

	// IntelliDoseDeviceNameLinux is the name that appears on Linux
	IntelliDoseDeviceNameLinux = "ASL IntelliDose"

	// IntelliDoseDeviceType is the device type (short)
	IntelliDoseDeviceType = "idoze"

	// IntelliClimateDeviceName is the name that appears on Windows
	IntelliClimateDeviceName = "IntelliClimate"

	// IntelliClimateDeviceNameLinux is the name that appears on Linux
	IntelliClimateDeviceNameLinux = "ASL IntelliClimate"

	// IntelliClimateDeviceType is the device type (short)
	IntelliClimateDeviceType = "iclimate"
)

var validDevices = []string{
	IntelliDoseDeviceName,
	IntelliClimateDeviceName,
	IntelliDoseDeviceNameLinux,
	IntelliClimateDeviceNameLinux,
}

func isValidDevice(name string) bool {
	return strings.Contains(":"+strings.Join(validDevices, ":")+":", name)
}

// Device represents an Intelli device
type Device struct {
	SerialNumber  string `json:"serial"`
	Name          string `json:"name"`
	DeviceType    string `json:"type"`
	hidDevice     *hidDevice
	HID           hid.DeviceInfo `json:"hid"`
	states        *States
	m             *sync.Mutex
	readWriteLock *sync.Mutex
	updating      *sync.Mutex
	Shadow        interface{} `json:"shadow"`
	IsOpen        bool        `json:"is_open"`
	onUpdateFunc  func(Device)
}

// NewDevice creates a new Intelli device from the given serial, type (dose or climate), name
// and HID device info
func NewDevice(sn, dtype, name string, hiddev hid.DeviceInfo) *Device {
	return &Device{
		SerialNumber:  sn,
		DeviceType:    dtype,
		hidDevice:     &hidDevice{hidDevice: hiddev},
		HID:           hiddev,
		states:        &States{},
		m:             &sync.Mutex{},
		readWriteLock: &sync.Mutex{},
		updating:      &sync.Mutex{},
	}
}

type hidDevice struct {
	hidDevice     hid.DeviceInfo
	hidDeviceImpl hid.Device
}

// States represents the different state packets recieved over USB
type States struct {
	d0State []byte
	d1State []byte
	d2State []byte
	d3State []byte
}

func (d *Device) update(shadow interface{}) {
	d.Shadow = shadow
	d.onUpdated()
}

func (d *Device) onUpdated() {
	go d.onUpdateFunc(*d)
}

// OnUpdate adds a single callback function to be called whenever the devices
// attributes are updated
func (d *Device) OnUpdate(callback func(Device)) {
	d.onUpdateFunc = callback
}

func (d *Device) close() error {
	d.IsOpen = false
	d.hidDevice.hidDeviceImpl.Close()
	return nil
}

func (d *Device) open() error {
	dev, err := d.hidDevice.hidDevice.Open()
	if err != nil {
		d.IsOpen = false
		return err
	}

	d.hidDevice.hidDeviceImpl = dev
	d.IsOpen = true
	return nil
}

func (d *Device) sentRequest(request []byte) ([]byte, error) {
	d.m.Lock()
	defer d.m.Unlock()
	time.Sleep(time.Millisecond * time.Duration(100))
	err := d.hidDevice.hidDeviceImpl.Write(request)
	tell.IfErrorf(err, "Error during sending data to device : ", d.SerialNumber)
	return <-d.hidDevice.hidDeviceImpl.ReadCh(), err
}
