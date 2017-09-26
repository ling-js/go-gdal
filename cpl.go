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
	"unsafe"
)

// ConfigOption reads an internal configuration option.
//
// The value is the value of a (key, value) option set with
// CPLSetConfigOption().  If the given option was no defined with
// CPLSetConfigOption(), it tries to find it in environment variables.
//
// Note: the string returned by CPLGetConfigOption() might be short-lived, and
// in particular it will become invalid after a call to CPLSetConfigOption()
// with the same key.
//
// To override temporary a potentially existing option with a new value, you
// can use the following snippet :
//
//     // backup old value
//     const char* pszOldValTmp = CPLGetConfigOption(pszKey, NULL);
//     char* pszOldVal = pszOldValTmp ? CPLStrdup(pszOldValTmp) : NULL;
//     // override with new value
//     CPLSetConfigOption(pszKey, pszNewVal);
//     // do something useful
//     // restore old value
//     CPLSetConfigOption(pszKey, pszOldVal);
//     CPLFree(pszOldVal);
//
//
// see SetConfigOption()
//
// http://trac.osgeo.org/gdal/wiki/ConfigOptions
func ConfigOption(key, def string) string {
	k := C.CString(key)
	d := C.CString(def)
	defer C.free(unsafe.Pointer(k))
	defer C.free(unsafe.Pointer(d))
	opt := C.CPLGetConfigOption(k, d)
	return C.GoString(opt)
}

// SetConfigOption
//
// Set a configuration option for GDAL/OGR use.
//
// Those options are defined as a (key, value) couple. The value corresponding
// to a key can be got later with the CPLGetConfigOption() method.
//
// This mechanism is similar to environment variables, but options set with
// CPLSetConfigOption() overrides, for CPLGetConfigOption() point of view,
// values defined in the environment.
//
// If CPLSetConfigOption() is called several times with the same key, the
// value provided during the last call will be used.
//
// Options can also be passed on the command line of most GDAL utilities
// with the with '--config KEY VALUE'. For example,
// ogrinfo --config CPL_DEBUG ON ~/data/test/point.shp
//
// This function can also be used to clear a setting by passing NULL as the
// value (note: passing NULL will not unset an existing environment variable;
// it will just unset a value previously set by CPLSetConfigOption()).
//
// see GetConfigOption
//
// http://trac.osgeo.org/gdal/wiki/ConfigOptions
func SetConfigOption(key, value string) {
	k := C.CString(key)
	v := C.CString(value)
	defer C.free(unsafe.Pointer(k))
	defer C.free(unsafe.Pointer(v))
	C.CPLSetConfigOption(k, v)
}

// HTTPEnabled returns if CPLHTTP services can be useful
//
// Those services depend on GDAL being build with libcurl support.
func HTTPEnabled() bool {
	rc := C.CPLHTTPEnabled()
	if int(rc) == 0 {
		return false
	}
	return true
}

// PushQuietHandler installs the no-op error handling mechanism in GDAL.  No
// output is printed out for warnings or errors.  Use PopHandler() to uninstall
// the quiet handler.
func PushQuietHandler() {
	C.CPLPushErrorHandler(C.CPLErrorHandler(C.CPLQuietErrorHandler))
}

// PopHandler pops the current error handler off of the error handling function
// stack.
func PopHandler() {
	C.CPLPopErrorHandler()
}
