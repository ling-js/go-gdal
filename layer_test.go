package gdal

import "testing"

func getLayer(t *testing.T) *Layer {
	fname := "/vsizip/test/gdalautotest-2.1.1.zip/gdalautotest-2.1.1/ogr/data/testshp/poly.shp"
	ds, err := OpenEx(fname, VectorDrivers|ReadOnly, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	lyr, err := ds.LayerByIndex(0)
	if err != nil {
		t.Fatal(err)
	}
	return &lyr
}

func TestLayerName(t *testing.T) {
	lyr := getLayer(t)
	if lyr.Name() != "poly" {
		t.Errorf("invalid layer name: %s", lyr.Name())
	}
}

func TestLayerGeomType(t *testing.T) {
	lyr := getLayer(t)
	if lyr.Type() != GT_Polygon {
		t.Errorf("invalid geometry type: %+v", lyr.Type())
	}
}

func TestNextFeature(t *testing.T) {
	lyr := getLayer(t)
	feat := lyr.NextFeature()
	if feat == nil {
		t.Fatal("nil feature handle")
	}
}

func TestFeatureCount(t *testing.T) {
	lyr := getLayer(t)
	if n, ok := lyr.FeatureCount(true); !ok || n != 10 {
		t.Errorf("invalid feature count: %d", n)
	}
}

func TestSpatialFilter(t *testing.T) {
	// Extent: (478315.531250, 4762880.500000) - (481645.312500, 4765610.500000)
	lyr := getLayer(t)
	if lyr.SpatialFilter() != nil {
		t.Errorf("spatial filter not nil")
	}
	// Create a spatial filter that has no features
	geom, err := CreateFromWKT("POLYGON((0 0, 1 0, 1 1, 0 1, 0 0))", lyr.SpatialReference())
	if err != nil {
		t.Fatal("could not create filter geometry")
	}
	lyr.SetSpatialFilter(&geom)
	lyr.ResetReading()
	if lyr.SpatialFilter() == nil {
		t.Error("spatial filter not set")
	}
	if n, ok := lyr.FeatureCount(true); ok != true || n != 0 {
		t.Error("spatial filter not set")
	}
	lyr.SetSpatialFilter(nil)
	lyr.ResetReading()
	if lyr.SpatialFilter() != nil {
		t.Error("spatial filter set")
	}
	if n, ok := lyr.FeatureCount(true); ok != true || n != 10 {
		t.Error("spatial filter set")
	}
	// TODO(kyle): add test for filter inclusion, use SetSpatialFilterRect
}
