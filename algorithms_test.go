package gdal

import "testing"

func TestCreateGrid(t *testing.T) {
	SetConfigOption("GDAL_NUM_THREADS", "1")
	alg := NearestNeighbor
	opt := NearestNeighborOptions{10, 10, 0, -9999}
	x := []float64{1, 2, 3, 1, 2, 3, 1, 2, 3}
	y := []float64{1, 1, 1, 2, 2, 2, 3, 3, 3}
	z := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	out := make([]float64, 9)
	rc := GridCreate(alg, opt, x, y, z, 0., 4., 0., 4., 3, 3, out, nil, nil)
	if rc != 0 {
		t.Fail()
	}
	for i := 0; i < 9; i++ {
		if out[i] < 0 {
			t.Errorf("Invalid pixel output")
		}
	}
}
