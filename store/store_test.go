package store

import "testing"
import "reflect"

type inOutTest struct {
	in  []string
	out []TotalAndFavorite
}

var toTest = []inOutTest {
	{ []string{"a", "b", "3"}, []TotalAndFavorite {{"a", "b", 3, 3}} },
	{ []string{"a", "b", "3", "a", "c", "4"}, []TotalAndFavorite {{"a", "c", 4, 7}} },
	{ []string{"a", "b", "3", "c", "b", "4"}, []TotalAndFavorite {{"c", "b", 4, 4}, {"a", "b", 3, 3}} },
}

func TestGetTotalAndFavourite(t *testing.T) {
	for _, v := range toTest {
		result := getTotalAndFavourite(v.in)
		if !reflect.DeepEqual(result, v.out) {
			t.Errorf("getTotalAndFavourite expected: %v, actual: %v", v.out, result)
		}
	}
}

type inOutTestTd struct {
	in  []byte
	out []string
	err bool
}

var toTestTds = []inOutTestTd {
	{ []byte("a"), []string{}, true },
	{ []byte(`<table id="top.customers" class="top.customers details"><tbody></tbody></table>`), []string{}, false },
	{ []byte(`<table id="top.customers" class="top.customers details"><tbody><td>td</td></tbody></table>`), []string{"td"}, false },
}

func TestGetCustomersTds(t *testing.T) {
	for _, v := range toTestTds {
		result, err := getCustomersTds(v.in)
		if v.err {
			if err == nil {
				t.Error("getCustomersTds should return an error")
			}
		} else {
			if err != nil {
				t.Error("getCustomersTds returned an error")
			}
			if !reflect.DeepEqual(result, v.out) {
				t.Errorf("getCustomersTds expected: %v, actual: %v", v.out, result)
			}
		}
	}
}
