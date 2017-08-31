package device

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/AutogrowSystems/go-intelli/hid"
	"github.com/AutogrowSystems/go-intelli/util/tell"
)

// NewManager will return a new device manager with the given intervals
func NewManager(enumerateInterval, updateInterval int) *Manager {
	mgr := &Manager{
		mutex:             new(sync.RWMutex),
		enumerateInterval: time.Duration(enumerateInterval) * time.Second,
		updateInterval:    time.Duration(updateInterval) * time.Second,
		devices:           []*Device{},
		deviceUpdatedFunc: func(d Device) {},
	}

	return mgr
}

// Manager represents a manager of attached Intelli devices
type Manager struct {
	devices           []*Device
	enumerateInterval time.Duration
	updateInterval    time.Duration
	mutex             *sync.RWMutex
	deviceUpdatedFunc func(Device)
}

// OnDeviceUpdated allows a callback to be fired whenever a device is updated
func (mgr *Manager) OnDeviceUpdated(callback func(Device)) {
	mgr.deviceUpdatedFunc = callback
}

// AttachAPI attaches a GIN API engine to the manager so it's internals can be inspected
// via HTTP REST endpoints
func (mgr *Manager) AttachAPI(r *gin.Engine) {

	r.GET("/devices/count", func(c *gin.Context) {
		count := len(mgr.devices)
		c.JSON(200, struct {
			Count int `json:"count"`
		}{count})
	})

	r.GET("/devices", func(c *gin.Context) {
		mgr.mutex.RLock()
		defer mgr.mutex.RUnlock()

		if len(mgr.devices) == 0 {
			c.AbortWithStatus(404)
			return
		}

		c.JSON(200, mgr.devices)
	})

}

// Interrogate will interrogate discovered devices for their readings and
// update their local shadow.
func (mgr *Manager) Interrogate() {
	for {
		for _, device := range mgr.devices {
			if !device.IsOpen {
				if err := device.open(); err != nil {
					tell.Errorf("%s", err)
					continue
				}
			}

			go device.updateShadow()
		}

		time.Sleep(mgr.updateInterval)
	}
}

// Discover will continuously try to discover devices attached via USB and add
// them to the internal slice of devices.  It will rediscover every time the
// enumerateInterval setting on the manager is passed.
func (mgr *Manager) Discover() {
	for {
		devicesInfo, err := hid.Devices()
		if err != nil {
			tell.IfErrorf(err, "failed to enumerate devices")
			time.Sleep(5 * time.Second)
			continue
		}

		mgr.addDevices(devicesInfo)

		if len(mgr.devices) == 0 {
			tell.Warnf("No Autogrow device is connected")
		}

		mgr.purgeDevices(devicesInfo)

		time.Sleep(mgr.enumerateInterval)
	}
}

func (mgr *Manager) addDevices(devicesInfo []*hid.DeviceInfo) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for _, info := range devicesInfo {
		name := info.Product
		sn := info.SerialNumber
		deviceType := "Unknown"

		if !isValidDevice(name) {
			continue
		}

		tell.Debugf("found device %s %s", name, sn)

		if _, found := mgr.FindDevice(info.SerialNumber); found {
			// 	d.hidDevice = &HIDDevice{
			// 		hidDevice: *info,
			// 	}
			//
			continue
		}

		switch name {
		case IntelliDoseDeviceName, IntelliDoseDeviceNameLinux:
			deviceType = IntelliDoseDeviceType
		case IntelliClimateDeviceName, IntelliClimateDeviceNameLinux:
			deviceType = IntelliClimateDeviceType
		default:
		}

		newdev := NewDevice(sn, deviceType, name, *info)
		newdev.OnUpdate(mgr.deviceUpdatedFunc)

		mgr.devices = append(mgr.devices, newdev)

		tell.Infof("connected device %s", sn)
	}
}

func (mgr *Manager) purgeDevices(devicesInfo []*hid.DeviceInfo) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if len(mgr.devices) > 0 {
		return
	}

	for i, d := range mgr.devices {
		found := false
		for _, info := range devicesInfo {
			if d.SerialNumber == info.SerialNumber {
				found = true
			}
		}

		if !found {
			d.close()
			mgr.devices = append(mgr.devices[:i], mgr.devices[i+1:]...)
		}
	}
}

// FindDevice returns the device by the given serial number and true, or else it will
// return nil device and false.
func (mgr *Manager) FindDevice(sn string) (*Device, bool) {
	if len(mgr.devices) == 0 {
		return nil, false
	}

	for _, d := range mgr.devices {
		if d.SerialNumber == sn {
			return d, true
		}
	}

	return nil, false
}

// HasDevice returns true if the manager contains the device with the given serial number
func (mgr *Manager) HasDevice(serialNumber string) bool {
	_, found := mgr.FindDevice(serialNumber)
	return found
}
