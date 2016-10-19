package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"
#include "cpl_http.h"

#cgo linux  pkg-config: gdal
#cgo darwin pkg-config: gdal
#cgo windows LDFLAGS: -Lc:/gdal/release-1600-x64/lib -lgdal_i
#cgo windows CFLAGS: -IC:/gdal/release-1600-x64/include
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Return the driver by short name
func GetDriverByName(driverName string) (Driver, error) {
	cName := C.CString(driverName)
	defer C.free(unsafe.Pointer(cName))

	driver := C.GDALGetDriverByName(cName)
	if driver == nil {
		return Driver{driver}, fmt.Errorf("Error: driver '%s' not found", driverName)
	}
	return Driver{driver}, nil
}

// Fetch the number of registered drivers.
func GetDriverCount() int {
	nDrivers := C.GDALGetDriverCount()
	return int(nDrivers)
}

// Fetch driver by index
func GetDriver(index int) Driver {
	driver := C.GDALGetDriver(C.int(index))
	return Driver{driver}
}

// Destroy a GDAL driver
func (driver Driver) Destroy() {
	C.GDALDestroyDriver(driver.cval)
}

// Registers a driver for use
func (driver Driver) Register() int {
	index := C.GDALRegisterDriver(driver.cval)
	return int(index)
}

// Reregister the driver
func (driver Driver) Deregister() {
	C.GDALDeregisterDriver(driver.cval)
}

// Destroy the driver manager
func DestroyDriverManager() {
	C.GDALDestroyDriverManager()
}

// Delete named dataset
func (driver Driver) DeleteDataset(name string) error {
	cDriver := driver.cval
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.GDALDeleteDataset(cDriver, cName).Err()
}

// Rename named dataset
func (driver Driver) RenameDataset(newName, oldName string) error {
	cDriver := driver.cval
	cNewName := C.CString(newName)
	defer C.free(unsafe.Pointer(cNewName))
	cOldName := C.CString(oldName)
	defer C.free(unsafe.Pointer(cOldName))
	return C.GDALRenameDataset(cDriver, cNewName, cOldName).Err()
}

// Copy all files associated with the named dataset
func (driver Driver) CopyDatasetFiles(newName, oldName string) error {
	cDriver := driver.cval
	cNewName := C.CString(newName)
	defer C.free(unsafe.Pointer(cNewName))
	cOldName := C.CString(oldName)
	defer C.free(unsafe.Pointer(cOldName))
	return C.GDALCopyDatasetFiles(cDriver, cNewName, cOldName).Err()
}

// Get the short name associated with this driver
func (driver Driver) ShortName() string {
	cDriver := driver.cval
	return C.GoString(C.GDALGetDriverShortName(cDriver))
}

// Get the long name associated with this driver
func (driver Driver) LongName() string {
	cDriver := driver.cval
	return C.GoString(C.GDALGetDriverLongName(cDriver))
}
