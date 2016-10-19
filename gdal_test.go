// Copyright 2011 go-gdal. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdal

import (
	"testing"
)

func TestTiffDriver(t *testing.T) {
	_, err := GetDriverByName("GTiff")
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestShapeDriver(t *testing.T) {
	drv := OGRDriverByName("ESRI Shapefile")
	drv.TestCapability("NA")
}

func TestOpen(t *testing.T) {
	if !HTTPEnabled() {
		t.Skip()
	}
	ds, err := Open("/vsicurl/http://download.osgeo.org/gdal/data/gtiff/small_world.tif", ReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	ds.Close()
}

func TestOpenEx(t *testing.T) {
	if !HTTPEnabled() {
		t.Skip()
	}
	ds, err := OpenEx("/vsicurl/http://download.osgeo.org/gdal/data/gtiff/small_world.tif",
		ReadOnly|RasterDrivers, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	ds.Close()
}

// Test open with flags we know should fail
func TestOpenExMismatchFlags(t *testing.T) {
	if !HTTPEnabled() {
		t.Skip()
	}
	tests := []struct {
		flags   Access
		drivers []string
	}{
		{
			ReadOnly | RasterDrivers,
			[]string{"JPEG", "PNG"},
		},
		{
			Update | RasterDrivers,
			nil,
		},
		{
			ReadOnly | VectorDrivers,
			nil,
		},
	}

	for _, test := range tests {
		_, err := OpenEx("/vsicurl/http://download.osgeo.org/gdal/data/gtiff/small_world.tif",
			test.flags, test.drivers, nil, nil)
		if err == nil {
			t.Fatal("driver sublist failed")
		}
	}
}

func TestHistogram(t *testing.T) {
	drv, err := GetDriverByName("MEM")
	if err != nil {
		t.Errorf("%+v", err)
	}
	ds := drv.Create("/vsimem/tmp", 10, 10, 1, Byte, nil)
	defer ds.Close()
	band := ds.RasterBand(1)
	data := make([]uint8, 100)
	for i := uint8(0); i < 100; i++ {
		data[i] = i
	}
	err = band.IO(Write, 0, 0, 10, 10, data, 10, 10, 0, 0)
	if err != nil {
		t.Errorf("%+v", err)
	}
	band.FlushCache()
	hist, err := band.Histogram(0, 100, 10, 1, 0, DummyProgress, nil)
	if err != nil {
		t.Errorf("%+v", err)
	}
	for i := 0; i < 10; i++ {
		if hist[i] != 10 {
			t.Errorf("failed to compute histogram. got: %+v, expected: 10\n", hist[i])
		}
	}
}

func TestGetLayer(t *testing.T) {
	if !HTTPEnabled() {
		t.Skip()
	}
	fname := "/vsizip/vsicurl/"
	fname += "http://download.osgeo.org"
	fname += "/gdal/CURRENT/gdalautotest-2.1.1.zip"
	fname += "/gdalautotest-2.1.1/ogr/data/testshp/poly.shp"
	ds, err := OpenEx(fname, ReadOnly|VectorDrivers, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ds.LayerByIndex(0)
	if err != nil {
		t.Error(err)
	}
	_, err = ds.LayerByName("poly")
	if err != nil {
		t.Error(err)
	}
}

func TestExecuteSQL(t *testing.T) {
	if !HTTPEnabled() {
		t.Skip()
	}
	fname := "/vsizip/vsicurl/"
	fname += "http://download.osgeo.org"
	fname += "/gdal/CURRENT/gdalautotest-2.1.1.zip"
	fname += "/gdalautotest-2.1.1/ogr/data/testshp/poly.shp"
	ds, err := OpenEx(fname, ReadOnly|VectorDrivers, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	lyr, err := ds.ExecuteSQL("SELECT * FROM poly", Geometry{}, "")
	if err != nil {
		t.Error(err)
	}
	if n, ok := lyr.FeatureCount(true); !ok || n < 1 {
		t.Error("failed to get a valid layer")
	}
}

func TestConfigOption(t *testing.T) {
	k, v := "GDAL_GO_TEST", "ON"
	SetConfigOption(k, v)
	value := GetConfigOption(k, "")
	if value != v {
		t.Errorf("Invalid value: %s\n", value)
	}
}
