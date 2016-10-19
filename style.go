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

type StyleMgr struct {
	cval C.OGRStyleMgrH
}

type StyleTool struct {
	cval C.OGRStyleToolH
}

type StyleTable struct {
	cval C.OGRStyleTableH
}

// Unimplemented: CreateStyleManager

// Unimplemented: Destroy

// Unimplemented: InitFromFeature

// Unimplemented: InitStyleString

// Unimplemented: PartCount

// Unimplemented: PartCount

// Unimplemented: AddPart

// Unimplemented: AddStyle

// Unimplemented: CreateStyleTool

// Unimplemented: Destroy

// Unimplemented: Type

// Unimplemented: Unit

// Unimplemented: SetUnit

// Unimplemented: ParamStr

// Unimplemented: ParamNum

// Unimplemented: ParamDbl

// Unimplemented: SetParamStr

// Unimplemented: SetParamNum

// Unimplemented: SetParamDbl

// Unimplemented: StyleString

// Unimplemented: RGBFromString

// Unimplemented: CreateStyleTable

// Unimplemented: Destroy

// Unimplemented: Save

// Unimplemented: Load

// Unimplemented: Find

// Unimplemented: ResetStyleStringReading

// Unimplemented: NextStyle

// Unimplemented: LastStyleName
