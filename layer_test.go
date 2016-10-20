package gdal

import "testing"

func getLayer(t *testing.T) *Layer {
	fname := "test/poly.shp"
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

func TestSpatialFilters(t *testing.T) {
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
	lyr.SetSpatialFilterRect(478903, 4764757, 478904, 4764756)
	lyr.ResetReading()
	if lyr.SpatialFilter() == nil {
		t.Error("spatial filter not set")
	}
	n, ok := lyr.FeatureCount(true)
	if ok != true {
		t.Error("failed to get feature count")
	}
	if n != 1 {
		t.Error("spatial filter does not include expected feature")
	}
	feat := lyr.NextFeature()
	if feat == nil || feat.FieldAsInteger(feat.FieldIndex("EAS_ID")) != 173 {
		t.Error("failed to get feature with proper fid")
	}
}
