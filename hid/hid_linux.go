package hid

// #include <linux/hidraw.h>
import "C"

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"

	"github.com/AutogrowSystems/go-intelli/util/tell"
)

var (
	ioctlHIDIOCGRDESCSIZE = ioR('H', 0x01, C.sizeof_int)
	ioctlHIDIOCGRDESC     = ioR('H', 0x02, C.sizeof_struct_hidraw_report_descriptor)
	ioctlHIDIOCGRAWINFO   = ioR('H', 0x03, C.sizeof_struct_hidraw_devinfo)
	hidUniq              = "HID_UNIQ"
)

func ioctlHIDIOCGRAWNAME(size int) uintptr {
	return ioR('H', 0x04, uintptr(size))
}

func ioctlHIDIOCGRAWPHYS(size int) uintptr {
	return ioR('H', 0x05, uintptr(size))
}

func ioctlHIDIOCSFEATURE(size int) uintptr {
	return ioRW('H', 0x06, uintptr(size))
}

func ioctlHIDIOCGFEATURE(size int) uintptr {
	return ioRW('H', 0x07, uintptr(size))
}

type linuxDevice struct {
	f    *os.File
	info *DeviceInfo

	writeLock *sync.Mutex
	readSetup sync.Once
	readErr   error
	readCh    chan []byte
}

// Devices enumerates the attached USB HID devices as DeviceInfo objects
func Devices() ([]*DeviceInfo, error) {
	sys, err := os.Open("/sys/class/hidraw")
	if err != nil {
		return nil, err
	}

	names, err := sys.Readdirnames(0)
	sys.Close()
	if err != nil {
		return nil, err
	}

	tell.Debugf("found %d HID devices", len(names))

	var res []*DeviceInfo
	for _, dir := range names {
		path := filepath.Join("/dev", filepath.Base(dir))

		tell.Debugf("reading info for %s", path)
		info, err := getDeviceInfo(path)

		if os.IsPermission(err) {
			tell.Errorf("permissions error reading %s: %s", path, err)
			continue
		} else if err != nil {
			tell.Errorf("error reading %s: %s", path, err)
			return nil, err
		}

		tell.Debugf("got info for %s: %+v", path, info)

		res = append(res, info)
	}

	tell.Debugf("got info for %d devices", len(res))

	return res, nil
}

func getDeviceInfo(path string) (*DeviceInfo, error) {
	d := &DeviceInfo{
		Path: path,
	}

	dev, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	defer dev.Close()
	fd := uintptr(dev.Fd())

	var descSize C.int
	if err := ioctl(fd, ioctlHIDIOCGRDESCSIZE, uintptr(unsafe.Pointer(&descSize))); err != nil {
		return nil, err
	}

	rawDescriptor := C.struct_hidraw_report_descriptor{
		size: C.__u32(descSize),
	}
	if err := ioctl(fd, ioctlHIDIOCGRDESC, uintptr(unsafe.Pointer(&rawDescriptor))); err != nil {
		return nil, err
	}
	d.parseReport(C.GoBytes(unsafe.Pointer(&rawDescriptor.value), descSize))

	var rawInfo C.struct_hidraw_devinfo
	if err := ioctl(fd, ioctlHIDIOCGRAWINFO, uintptr(unsafe.Pointer(&rawInfo))); err != nil {
		return nil, err
	}
	d.VendorID = uint16(rawInfo.vendor)
	d.ProductID = uint16(rawInfo.product)

	rawName := make([]byte, 256)
	if err := ioctl(fd, ioctlHIDIOCGRAWNAME(len(rawName)), uintptr(unsafe.Pointer(&rawName[0]))); err != nil {
		return nil, err
	}
	d.Product = string(rawName[:bytes.IndexByte(rawName, 0)])

	if p, err := filepath.EvalSymlinks(filepath.Join("/sys/class/hidraw", filepath.Base(path), "device")); err == nil {
		if rawManufacturer, err := ioutil.ReadFile(filepath.Join(p, "/../../manufacturer")); err == nil {
			d.Manufacturer = string(bytes.TrimRight(rawManufacturer, "\n"))
		}
		config := make(map[string]string)
		if ueventFile, err := os.Open(filepath.Join(p, "uevent")); err == nil {
			defer ueventFile.Close()
			scanner := bufio.NewScanner(ueventFile)
			for scanner.Scan() {
				line := scanner.Text()
				if equal := strings.Index(line, "="); equal >= 0 {
					if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
						value := ""
						if len(line) > equal {
							value = strings.TrimSpace(line[equal+1:])
						}
						config[key] = value
					}
				}
			}
		}
		if len(config[hidUniq]) > 0 {
			d.SerialNumber = config[hidUniq]
		}
	}

	return d, nil
}

// very basic report parser that will pull out the usage page, usage, and the
// sizes of the first input and output reports
func (d *DeviceInfo) parseReport(b []byte) {
	var reportSize uint16

	for len(b) > 0 {
		// read item size, type, and tag
		size := int(b[0] & 0x03)
		if size == 3 {
			size = 4
		}
		typ := (b[0] >> 2) & 0x03
		tag := (b[0] >> 4) & 0x0f
		b = b[1:]

		if len(b) < size {
			return
		}

		// read item value
		var v uint64
		for i := 0; i < size; i++ {
			v += uint64(b[i]) << (8 * uint(i))
		}
		b = b[size:]

		switch {
		case typ == 0 && tag == 8 && d.InputReportLength == 0 && reportSize > 0: // input report type
			d.InputReportLength = reportSize
			reportSize = 0
		case typ == 0 && tag == 9 && d.OutputReportLength == 0 && reportSize > 0: // output report type
			d.OutputReportLength = reportSize
			reportSize = 0
		case typ == 1 && tag == 0: // usage page
			d.UsagePage = uint16(v)
		case typ == 1 && tag == 9: // report count
			reportSize = uint16(v)
		case typ == 2 && tag == 0 && d.Usage == 0: // usage
			d.Usage = uint16(v)
		}

		if d.UsagePage > 0 && d.Usage > 0 && d.InputReportLength > 0 && d.OutputReportLength > 0 {
			return
		}
	}
}

// ByPath returns device info via it's file path
func ByPath(path string) (*DeviceInfo, error) {
	return getDeviceInfo(path)
}

// Open will satisfy the hid.Device interface, opening the device for reading writing
func (d *DeviceInfo) Open() (Device, error) {
	f, err := os.OpenFile(d.Path, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	return &linuxDevice{f: f, info: d, writeLock: new(sync.Mutex)}, nil
}

func (d *linuxDevice) Close() {
	d.f.Close()
}

func (d *linuxDevice) Write(data []byte) error {
	d.writeLock.Lock()
	defer d.writeLock.Unlock()
	_, err := d.f.Write(data)
	return err
}

func (d *linuxDevice) ReadCh() <-chan []byte {
	d.readSetup.Do(func() {
		d.readCh = make(chan []byte, 30)
		go d.readThread()
	})
	return d.readCh
}

func (d *linuxDevice) ReadError() error {
	return d.readErr
}

func (d *linuxDevice) readThread() {
	defer close(d.readCh)
	for {
		buf := make([]byte, d.info.InputReportLength)
		n, err := d.f.Read(buf)
		if err != nil {
			d.readErr = err
			return
		}
		select {
		case d.readCh <- buf[:n]:
		default:
		}
	}
}
