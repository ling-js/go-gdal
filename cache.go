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

/* ==================================================================== */
/*      GDAL Cache Management                                           */
/* ==================================================================== */

// Set maximum cache memory
func SetCacheMax(bytes int) {
	C.GDALSetCacheMax64(C.GIntBig(bytes))
}

// Get maximum cache memory
func GetCacheMax() int {
	bytes := C.GDALGetCacheMax64()
	return int(bytes)
}

// Get cache memory used
func GetCacheUsed() int {
	bytes := C.GDALGetCacheUsed64()
	return int(bytes)
}

// Try to flush one cached raster block
func FlushCacheBlock() bool {
	flushed := C.GDALFlushCacheBlock()
	return flushed != 0
}
