package kons

import (
	"reflect"
	"testing"
)

func TestKone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		data interface{}
	}{
		{"jess"},
		{21},
	}
	for _, test := range tests {
		kone := NewKone(nil)
		kone.SetData(test.data)
		if !reflect.DeepEqual(kone.GetData(), test.data) {
			t.Errorf("Data mismatch: want %v, have %v", test.data, kone.GetData())
		}
	}
}

func TestKoneUpsert(t *testing.T) {
	t.Parallel()

	kone := NewKone(nil)
	tests := []struct {
		settings *UpsertSettings
		val      interface{}
		path     []string
		wantn    int
	}{
		{nil, "john", []string{"host", "owner"}, 1},
		{&UpsertSettings{MakePath: false}, "localhost", []string{"host", "hostname"}, 0},
	}

	for _, test := range tests {
		n, err := kone.Upsert(test.settings, test.val, test.path...)
		if err != nil {
			t.Errorf("Upsert error: %+v", err)
			continue
		}
		if n != test.wantn {
			t.Errorf("Invalid upserts count: want %d, have %d", test.wantn, n)
		}
	}
}

func TestKoneFind(t *testing.T) {
	t.Parallel()

	kone := NewKone(nil)
	kone.Upsert(nil, 12, "usa", "node1", "clients", "john", "balance")
	kone.Upsert(nil, 21, "usa", "node1", "clients", "jessy", "balance")
	kone.Upsert(nil, 32, "usa", "node1", "clients", "bob", "balance")
	tests := []struct {
		path  []string
		wantn int
	}{
		{[]string{"uk", "node1"}, 0},
		{[]string{"usa", "node1", "clients", "[j].*", "balance"}, 2},
		{[]string{"usa", "node1", "clients", "[j].*n$", "balance"}, 1},
	}

	for _, test := range tests {
		data, err := kone.Find(nil, test.path...)
		if err != nil {
			t.Errorf("find error: %+v", err)
			continue
		}
		if len(data) != test.wantn {
			t.Errorf("invalid found count: want %d, have %d", test.wantn, len(data))
		}
	}
}
