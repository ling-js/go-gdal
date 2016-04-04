package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"

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

const (
	DMD_LONGNAME           = string(C.GDAL_DMD_LONGNAME)
	DMD_HELPTOPIC          = string(C.GDAL_DMD_HELPTOPIC)
	DMD_MIMETYPE           = string(C.GDAL_DMD_MIMETYPE)
	DMD_EXTENSION          = string(C.GDAL_DMD_EXTENSION)
	DMD_CREATIONOPTIONLIST = string(C.GDAL_DMD_CREATIONOPTIONLIST)
	DMD_CREATIONDATATYPES  = string(C.GDAL_DMD_CREATIONDATATYPES)
	DCAP_CREATE            = string(C.GDAL_DCAP_CREATE)
	DCAP_CREATECOPY        = string(C.GDAL_DCAP_CREATECOPY)
	DCAP_VIRTUALIO         = string(C.GDAL_DCAP_VIRTUALIO)
)

type Driver struct {
	cval C.GDALDriverH
}

// Create a new dataset with this driver.
func (driver Driver) Create(
	filename string,
	xSize, ySize, bands int,
	dataType DataType,
	options []string,
) Dataset {
	name := C.CString(filename)
	defer C.free(unsafe.Pointer(name))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	h := C.GDALCreate(
		driver.cval,
		name,
		C.int(xSize), C.int(ySize), C.int(bands),
		C.GDALDataType(dataType),
		(**C.char)(unsafe.Pointer(&opts[0])),
	)
	return Dataset{h}
}

// Create a copy of a dataset
func (driver Driver) CreateCopy(
	filename string,
	sourceDataset Dataset,
	strict int,
	options []string,
	progress ProgressFunc,
	data interface{},
) Dataset {
	name := C.CString(filename)
	defer C.free(unsafe.Pointer(name))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	var h C.GDALDatasetH

	if progress == nil {
		h = C.GDALCreateCopy(
			driver.cval, name,
			sourceDataset.cval,
			C.int(strict),
			(**C.char)(unsafe.Pointer(&opts[0])),
			nil,
			nil,
		)
	} else {
		arg := &goGDALProgressFuncProxyArgs{
			progress, data,
		}
		h = C.GDALCreateCopy(
			driver.cval, name,
			sourceDataset.cval,
			C.int(strict), (**C.char)(unsafe.Pointer(&opts[0])),
			C.goGDALProgressFuncProxyB(),
			unsafe.Pointer(arg),
		)
	}

	return Dataset{h}
}

// Return the driver needed to access the provided dataset name.
func IdentifyDriver(filename string, filenameList []string) Driver {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	length := len(filenameList)
	cFilenameList := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cFilenameList[i] = C.CString(filenameList[i])
		defer C.free(unsafe.Pointer(cFilenameList[i]))
	}
	cFilenameList[length] = (*C.char)(unsafe.Pointer(nil))

	driver := C.GDALIdentifyDriver(cFilename, (**C.char)(unsafe.Pointer(&cFilenameList[0])))
	return Driver{driver}
}

// Open opens a raster or vector file as a GDALDataset.
//
// This function will try to open the passed file, or virtual dataset name by
// invoking the Open method of each registered GDALDriver in turn. The first
// successful open will result in a returned dataset. If all drivers fail then
// NULL is returned and an error is issued.
//
// Several recommendations:
//
// If you open a dataset object with GDAL_OF_UPDATE access, it is not
// recommended to open a new dataset on the same underlying file.  The
// returned dataset should only be accessed by one thread at a time. If you
// want to use it from different threads, you must add all necessary code
// (mutexes, etc.) to avoid concurrent use of the object. (Some drivers, such
// as GeoTIFF, maintain internal state variables that are updated each time a
// new block is read, thus preventing concurrent use.)
//
// For drivers supporting the VSI virtual file API, it is possible to open a file
// in a .zip archive (see VSIInstallZipFileHandler()), in a .tar/.tar.gz/.tgz
// archive (see VSIInstallTarFileHandler()) or on a HTTP / FTP server (see
// VSIInstallCurlFileHandler())
//
// In some situations (dealing with unverified data), the datasets can be opened
// in another process through the GDAL API Proxy mechanism.
//
// In order to reduce the need for searches through the operating system file
// system machinery, it is possible to give an optional list of files with the
// papszSiblingFiles parameter. This is the list of all files at the same level
// in the file system as the target file, including the target file. The
// filenames must not include any path components, are an essentially just the
// output of CPLReadDir() on the parent directory. If the target object does
// not have filesystem semantics then the file list should be NULL.
func OpenEx(filename string, flags uint, drivers, options, siblings []string) (Dataset, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	length := len(drivers)
	drvs := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		drvs[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(drvs[i]))
	}
	drvs[length] = (*C.char)(unsafe.Pointer(nil))

	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	sibs := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		sibs[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(sibs[i]))
	}
	sibs[length] = (*C.char)(unsafe.Pointer(nil))

	dataset := C.GDALOpenEx(
		cFilename,
		C.uint(flags),
		(**C.char)(unsafe.Pointer(&drvs[0])),
		(**C.char)(unsafe.Pointer(&opts[0])),
		(**C.char)(unsafe.Pointer(&sibs[0])))
	if dataset == nil {
		return Dataset{nil}, fmt.Errorf("Error: dataset '%s' open error", filename)
	}
	return Dataset{dataset}, nil
}

func Open(filename string, access Access) (Dataset, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	dataset := C.GDALOpen(cFilename, C.GDALAccess(access))
	if dataset == nil {
		return Dataset{nil}, fmt.Errorf("Error: dataset '%s' open error", filename)
	}
	return Dataset{dataset}, nil
}

// Open a shared existing dataset
func OpenShared(filename string, access Access) Dataset {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	dataset := C.GDALOpenShared(cFilename, C.GDALAccess(access))
	return Dataset{dataset}
}

// Unimplemented: DumpOpenDatasets

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
