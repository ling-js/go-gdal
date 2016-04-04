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
import "unsafe"

type GridAlgorithm int

const (
	InverseDistanceToAPower                = GridAlgorithm(C.GGA_InverseDistanceToAPower)
	MovingAverage                          = GridAlgorithm(C.GGA_MovingAverage)
	NearestNeighbor                        = GridAlgorithm(C.GGA_NearestNeighbor)
	MetricMinimum                          = GridAlgorithm(C.GGA_MetricMinimum)
	MetricMaximum                          = GridAlgorithm(C.GGA_MetricMaximum)
	MetricRange                            = GridAlgorithm(C.GGA_MetricRange)
	MetricCount                            = GridAlgorithm(C.GGA_MetricCount)
	MetricAverageDistance                  = GridAlgorithm(C.GGA_MetricAverageDistance)
	MetricAverageDistancePts               = GridAlgorithm(C.GGA_MetricAverageDistancePts)
	Linear                                 = GridAlgorithm(C.GGA_Linear)
	InverseDistanceToAPowerNearestNeighbor = GridAlgorithm(C.GGA_InverseDistanceToAPowerNearestNeighbor)
)

/* --------------------------------------------- */
/* Misc functions                                */
/* --------------------------------------------- */

// Compute optimal PCT for RGB image
func ComputeMedianCutPCT(
	red, green, blue RasterBand,
	colors int,
	ct ColorTable,
	progress ProgressFunc,
	data interface{},
) int {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	err := C.GDALComputeMedianCutPCT(
		red.cval,
		green.cval,
		blue.cval,
		nil,
		C.int(colors),
		ct.cval,
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	)
	return int(err)
}

// 24bit to 8bit conversion with dithering
func DitherRGB2PCT(
	red, green, blue, target RasterBand,
	ct ColorTable,
	progress ProgressFunc,
	data interface{},
) int {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	err := C.GDALDitherRGB2PCT(
		red.cval,
		green.cval,
		blue.cval,
		target.cval,
		ct.cval,
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	)
	return int(err)
}

// Compute checksum for image region
func (rb RasterBand) Checksum(xOff, yOff, xSize, ySize int) int {
	sum := C.GDALChecksumImage(rb.cval, C.int(xOff), C.int(yOff), C.int(xSize), C.int(ySize))
	return int(sum)
}

// Compute the proximity of all pixels in the image to a set of pixels in the source image
func (src RasterBand) ComputeProximity(
	dest RasterBand,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALComputeProximity(
		src.cval,
		dest.cval,
		(**C.char)(unsafe.Pointer(&opts[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

// Fill selected raster regions by interpolation from the edges
func (src RasterBand) FillNoData(
	mask RasterBand,
	distance float64,
	iterations int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALFillNodata(
		src.cval,
		mask.cval,
		C.double(distance),
		0,
		C.int(iterations),
		(**C.char)(unsafe.Pointer(&opts[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

// Create polygon coverage from raster data using an integer buffer
func (src RasterBand) Polygonize(
	mask RasterBand,
	layer Layer,
	fieldIndex int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALPolygonize(
		src.cval,
		mask.cval,
		layer.cval,
		C.int(fieldIndex),
		(**C.char)(unsafe.Pointer(&opts[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

// Create polygon coverage from raster data using a floating point buffer
func (src RasterBand) FPolygonize(
	mask RasterBand,
	layer Layer,
	fieldIndex int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALFPolygonize(
		src.cval,
		mask.cval,
		layer.cval,
		C.int(fieldIndex),
		(**C.char)(unsafe.Pointer(&opts[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

// Removes small raster polygons
func (src RasterBand) SieveFilter(
	mask, dest RasterBand,
	threshold, connectedness int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	arg := &goGDALProgressFuncProxyArgs{
		progress, data,
	}

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return C.GDALSieveFilter(
		src.cval,
		mask.cval,
		dest.cval,
		C.int(threshold),
		C.int(connectedness),
		(**C.char)(unsafe.Pointer(&opts[0])),
		C.goGDALProgressFuncProxyB(),
		unsafe.Pointer(arg),
	).Err()
}

/* --------------------------------------------- */
/* Warp functions                                */
/* --------------------------------------------- */

//Unimplemented: CreateGenImgProjTransformer
//Unimplemented: CreateGenImgProjTransformer2
//Unimplemented: CreateGenImgProjTransformer3
//Unimplemented: SetGenImgProjTransformerDstGeoTransform
//Unimplemented: DestroyGenImgProjTransformer
//Unimplemented: GenImgProjTransform

//Unimplemented: CreateReprojectionTransformer
//Unimplemented: DestroyReprojection
//Unimplemented: ReprojectionTransform
//Unimplemented: CreateGCPTransformer
//Unimplemented: CreateGCPRefineTransformer
//Unimplemented: DestroyGCPTransformer
//Unimplemented: GCPTransform

//Unimplemented: CreateTPSTransformer
//Unimplemented: DestroyTPSTransformer
//Unimplemented: TPSTransform

//Unimplemented: CreateRPCTransformer
//Unimplemented: DestroyRPCTransformer
//Unimplemented: RPCTransform

//Unimplemented: CreateGeoLocTransformer
//Unimplemented: DestroyGeoLocTransformer
//Unimplemented: GeoLocTransform

//Unimplemented: CreateApproxTransformer
//Unimplemented: DestroyApproxTransformer
//Unimplemented: ApproxTransform

//Unimplemented: SimpleImageWarp
//Unimplemented: SuggestedWarpOutput
//Unimplemented: SuggsetedWarpOutput2
//Unimplemented: SerializeTransformer
//Unimplemented: DeserializeTransformer

//Unimplemented: TransformGeolocations

/* --------------------------------------------- */
/* Contour line functions                        */
/* --------------------------------------------- */

//Unimplemented: CreateContourGenerator
//Unimplemented: FeedLine
//Unimplemented: Destroy
//Unimplemented: ContourWriter
//Unimplemented: ContourGenerate

/* --------------------------------------------- */
/* Rasterizer functions                          */
/* --------------------------------------------- */

// Burn geometries into raster
//Unimplmemented: RasterizeGeometries

// Burn geometries from the specified list of layers into the raster
//Unimplemented: RasterizeLayers

// Burn geometries from the specified list of layers into the raster
//Unimplemented: RasterizeLayersBuf

/* --------------------------------------------- */
/* Gridding functions                            */
/* --------------------------------------------- */
// The option structs are fully public.

type InverseDistanceToAPowerOptions struct {
	Power           float64
	Smoothing       float64
	AnisotropyRatio float64
	AnisotropyAngle float64
	Radius1         float64
	Radius2         float64
	Angle           float64
	MaxPoints       uint32
	MinPoints       uint32
	NoDataValue     float64
	//C.GDALGridInverseDistanceToAPowerOptions
}

func (o InverseDistanceToAPowerOptions) gridOption() unsafe.Pointer {
	c := C.GDALGridInverseDistanceToAPowerOptions{
		dfPower:           C.double(o.Power),
		dfSmoothing:       C.double(o.Smoothing),
		dfAnisotropyRatio: C.double(o.AnisotropyRatio),
		dfAnisotropyAngle: C.double(o.AnisotropyAngle),
		dfRadius1:         C.double(o.Radius1),
		dfRadius2:         C.double(o.Radius2),
		dfAngle:           C.double(o.Angle),
		nMaxPoints:        C.GUInt32(o.MaxPoints),
		nMinPoints:        C.GUInt32(o.MinPoints),
		dfNoDataValue:     C.double(o.NoDataValue)}

	return unsafe.Pointer(&c)
}

type InverseDistanceToAPowerNearestNeighborOptions struct {
	Power       float64
	Radius      float64
	MaxPoints   uint32
	MinPoints   uint32
	NoDataValue float64
	//C.GDALGridInverseDistanceToAPowerNearestNeighborOptions
}

func (o InverseDistanceToAPowerNearestNeighborOptions) gridOption() unsafe.Pointer {
	c := C.GDALGridInverseDistanceToAPowerNearestNeighborOptions{
		dfPower:       C.double(o.Power),
		dfRadius:      C.double(o.Radius),
		nMaxPoints:    C.GUInt32(o.MaxPoints),
		nMinPoints:    C.GUInt32(o.MinPoints),
		dfNoDataValue: C.double(o.NoDataValue)}

	return unsafe.Pointer(&c)
}

type MovingAverageOptions struct {
	Radius1     float64
	Radius2     float64
	Angle       float64
	MinPoints   uint32
	NoDataValue float64
	//C.GDALGridMovingAverageOptions
}

func (o MovingAverageOptions) gridOption() unsafe.Pointer {
	c := C.GDALGridMovingAverageOptions{
		dfRadius1:     C.double(o.Radius1),
		dfRadius2:     C.double(o.Radius2),
		dfAngle:       C.double(o.Angle),
		nMinPoints:    C.GUInt32(o.MinPoints),
		dfNoDataValue: C.double(o.NoDataValue)}

	return unsafe.Pointer(&c)
}

type NearestNeighborOptions struct {
	Radius1     float64
	Radius2     float64
	Angle       float64
	NoDataValue float64
	//C.GDALGridNearestNeighborOptions
}

func (o NearestNeighborOptions) gridOption() unsafe.Pointer {
	c := C.GDALGridNearestNeighborOptions{
		dfRadius1:     C.double(o.Radius1),
		dfRadius2:     C.double(o.Radius2),
		dfAngle:       C.double(o.Angle),
		dfNoDataValue: C.double(o.NoDataValue)}

	return unsafe.Pointer(&c)
}

type DataMetricsOptions struct {
	Radius1     float64
	Radius2     float64
	Angle       float64
	MinPoints   uint32
	NoDataValue float64
	//C.GDALGridDataMetricsOptions
}

func (o DataMetricsOptions) gridOption() unsafe.Pointer {
	c := C.GDALGridDataMetricsOptions{
		dfRadius1:     C.double(o.Radius1),
		dfRadius2:     C.double(o.Radius2),
		dfAngle:       C.double(o.Angle),
		nMinPoints:    C.GUInt32(o.MinPoints),
		dfNoDataValue: C.double(o.NoDataValue)}

	return unsafe.Pointer(&c)
}

type LinearOptions struct {
	Radius      float64
	NoDataValue float64
	//C.GDALGridLinearOptions
}

func (o LinearOptions) gridOption() unsafe.Pointer {
	c := C.GDALGridLinearOptions{
		dfRadius:      C.double(o.Radius),
		dfNoDataValue: C.double(o.NoDataValue)}

	return unsafe.Pointer(&c)
}

type gridOption interface {
	gridOption() unsafe.Pointer
}

func GridCreate(
	algorithm GridAlgorithm,
	options gridOption,
	x, y, z []float64,
	xMin, xMax, yMin, yMax float64,
	xSize, ySize uint32,
	buffer interface{},
	progress ProgressFunc,
	pData interface{}) int {

	var dataType DataType
	var dataPtr unsafe.Pointer
	var length int
	switch data := buffer.(type) {
	case []int8:
		dataType = Byte
		dataPtr = unsafe.Pointer(&data[0])
		length = len(data)
	case []uint8:
		dataType = Byte
		dataPtr = unsafe.Pointer(&data[0])
		length = len(data)
	case []int16:
		dataType = Int16
		dataPtr = unsafe.Pointer(&data[0])
		length = len(data)
	case []uint16:
		dataType = UInt16
		dataPtr = unsafe.Pointer(&data[0])
		length = len(data)
	case []int32:
		dataType = Int32
		dataPtr = unsafe.Pointer(&data[0])
		length = len(data)
	case []uint32:
		dataType = UInt32
		dataPtr = unsafe.Pointer(&data[0])
		length = len(data)
	case []float32:
		dataType = Float32
		dataPtr = unsafe.Pointer(&data[0])
		length = len(data)
	case []float64:
		dataType = Float64
		dataPtr = unsafe.Pointer(&data[0])
		length = len(data)
	}

	err := C.GDALGridCreate(
		C.GDALGridAlgorithm(algorithm),
		options.gridOption(),
		C.GUInt32(length),
		(*C.double)(unsafe.Pointer(&x[0])),
		(*C.double)(unsafe.Pointer(&y[0])),
		(*C.double)(unsafe.Pointer(&z[0])),
		C.double(xMin), C.double(xMax),
		C.double(yMin), C.double(yMax),
		C.GUInt32(xSize), C.GUInt32(ySize),
		C.GDALDataType(dataType), dataPtr,
		nil, nil)

	return int(err)
}

//Unimplemented: ComputeMatchingPoints
