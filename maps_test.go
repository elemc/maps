package maps_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/elemc/maps"
)

var (
	gmap *maps.Map
)

func testCreateMap() (err error) {
	gmap = maps.New()

	someData := make(map[string]interface{})
	someData["test_bool"] = "true"
	someData["test_float"] = 1.2345
	someData["test_int"] = "12345"
	someData["test_str"] = "some string"
	someData["test_time"] = "1979-09-24T05:35:00Z"

	gmap.AddMap(someData)

	if gmap.GetMap().Length() != 5 {
		err = fmt.Errorf("unexpected map length: %d", gmap.GetMap().Length())
		return
	}

	return
}

func testGetValues() (err error) {
	if v := gmap.GetBool("test_bool"); !v {
		err = fmt.Errorf("unexpected boolean value: %t", v)
	}
	if v := gmap.GetFloat64("test_float"); v != 1.2345 {
		err = fmt.Errorf("unexpected float value: %f", v)
	}
	if v := gmap.GetInt64("test_float"); v != 1 {
		err = fmt.Errorf("unexpected int64 value: %d", v)
	}
	if v := gmap.GetInt("test_int"); v != 12345 {
		err = fmt.Errorf("unexpected int value: %d", v)
	}
	if v := gmap.GetString("test_str"); v != "some string" {
		err = fmt.Errorf("unexpected string value: %s", v)
	}
	if v := gmap.GetTime("test_time"); v.Unix() != 306999300 {
		err = fmt.Errorf("unexpected time value: %s (%d)", v, v.Unix())
	}
	if v := gmap.GetJSTime("test_time"); v != 306999300000 {
		err = fmt.Errorf("unexpected JS time value: %d", v)
	}

	return
}

func TestCreateMap(t *testing.T) {
	if err := testCreateMap(); err != nil {
		t.Fatal(err)
	}
}

func TestGetValues(t *testing.T) {
	if err := testGetValues(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkMap(b *testing.B) {
	delta := time.Now().Unix() * 1000
	gmap = maps.New()
	var gwg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		gwg.Add(1)
		go func(i int) {
			defer gwg.Done()
			var wg sync.WaitGroup
			wg.Add(1)
			go func(m *maps.Map, iter int) {
				defer wg.Done()
				gmap.Set(fmt.Sprintf("test_key_%d", iter), iter)
				gmap.SetJSTime(fmt.Sprintf("test_key_time_%d", iter), delta+int64(iter*1000))
			}(gmap, i)
			wg.Wait()

			wg.Add(1)
			go func(m *maps.Map, iter int) {
				defer wg.Done()
				testI := gmap.GetInt(fmt.Sprintf("test_key_%d", iter))
				if testI != iter {
					b.Fatalf("Unexpected value int: %d", testI)
				}
			}(gmap, i)

			wg.Add(1)
			go func(m *maps.Map, iter int) {
				defer wg.Done()
				testI := gmap.GetJSTime(fmt.Sprintf("test_key_time_%d", iter))
				if testI != delta+int64(iter*1000) {
					b.Fatalf("Unexpected value JS time: %d", testI)
				}
			}(gmap, i)
		}(i)
		gwg.Wait()
	}
}
