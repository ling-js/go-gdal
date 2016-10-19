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

/* ==================================================================== */
/*      Color tables.                                                   */
/* ==================================================================== */

// Construct a new color table
func CreateColorTable(interp PaletteInterp) ColorTable {
	ct := C.GDALCreateColorTable(C.GDALPaletteInterp(interp))
	return ColorTable{ct}
}

// Destroy the color table
func (ct ColorTable) Destroy() {
	C.GDALDestroyColorTable(ct.cval)
}

// Make a copy of the color table
func (ct ColorTable) Clone() ColorTable {
	newCT := C.GDALCloneColorTable(ct.cval)
	return ColorTable{newCT}
}

// Fetch palette interpretation
func (ct ColorTable) PaletteInterpretation() PaletteInterp {
	pi := C.GDALGetPaletteInterpretation(ct.cval)
	return PaletteInterp(pi)
}

// Get number of color entries in table
func (ct ColorTable) EntryCount() int {
	count := C.GDALGetColorEntryCount(ct.cval)
	return int(count)
}

// Fetch a color entry from table
func (ct ColorTable) Entry(index int) ColorEntry {
	entry := C.GDALGetColorEntry(ct.cval, C.int(index))
	return ColorEntry{entry}
}

// Unimplemented: EntryAsRGB

// Set entry in color table
func (ct ColorTable) SetEntry(index int, entry ColorEntry) {
	C.GDALSetColorEntry(ct.cval, C.int(index), entry.cval)
}

// Create color ramp
func (ct ColorTable) CreateColorRamp(start, end int, startColor, endColor ColorEntry) {
	C.GDALCreateColorRamp(ct.cval, C.int(start), startColor.cval, C.int(end), endColor.cval)
}
