package memcache

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func debugCache() CacheStore {
	return New(50, 10*time.Second, 500*time.Millisecond, 100*time.Millisecond)
}

func Test_cache_haveKey(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		capacity  uint64
		incr      uint64
		items     map[string]Item
		defaultlt time.Duration
		auditor   Auditor
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "case have key",
			fields: fields{
				items: map[string]Item{"test-key": Item{}},
			},
			args: args{
				key: "test-key",
			},
			want: true,
		},
		{
			name: "case key dosent exist",
			args: args{
				key: "test-key",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				mu:        tt.fields.mu,
				capacity:  tt.fields.capacity,
				incr:      tt.fields.incr,
				items:     tt.fields.items,
				defaultlt: tt.fields.defaultlt,
				auditor:   tt.fields.auditor,
			}
			if got := c.haveKey(tt.args.key); got != tt.want {
				t.Errorf("cache.haveKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_PutGet(t *testing.T) {
	c := debugCache()
	type args struct {
		key   string
		value interface{}
		tags  []uint16
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		wantErr bool
	}{
		{
			name:   "put new key",
			fields: c,
			args: args{
				key:   "test-key",
				value: 10,
			},
			wantErr: false,
		},
		{
			name:   "put the same key",
			fields: c,
			args: args{
				key:   "test-key",
				value: 10,
			},
			wantErr: true,
		},
		{
			name: "put another key",
			args: args{
				key:   "test-key-two",
				value: 777,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := c.Put(tt.args.key, tt.args.value, tt.args.tags...); (err != nil) != tt.wantErr {
				t.Errorf("cache.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
			if v, err := c.GetItem(tt.args.key); err != nil || v.Value != tt.args.value || (v.createdAt == time.Time{}) || (v.lifetime == time.Duration(0)) {
				t.Errorf("cache.Put() value = %v, want %v", tt.args.value, v)
			}
		})
	}
}

func Test_cache_Get(t *testing.T) {
	c := debugCache()
	c.Put("10__one", 10, 2, 3)
	c.Put("10__two", true, 2)
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:   "get existing key",
			fields: c,
			args: args{
				key: "10__one",
			},
			want:    10,
			wantErr: false,
		},
		{
			name:   "get existing key",
			fields: c,
			args: args{
				key: "10__two",
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "get non existing key",
			fields: c,
			args: args{
				key: "10__three",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("cache.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_GetItem(t *testing.T) {
	c := debugCache()
	c.Put("10__one", 10, 2, 3)
	c.Put("10__two", true, 2)
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:   "get existing item",
			fields: c,
			args: args{
				key: "10__one",
			},
			want:    10,
			wantErr: false,
		},
		{
			name:   "get existing item",
			fields: c,
			args: args{
				key: "10__two",
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "get non existing item",
			fields: c,
			args: args{
				key: "10__three",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetItem(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("cache.GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(got.Value == tt.want) || (got.lifetime == 0 && !tt.wantErr) {
				fmt.Println(got.createdAt, got.lifetime)
				t.Errorf("cache.GetItem() = %v, want %v", got.Value, tt.want)
			}
		})
	}
}

func Test_cache_Update(t *testing.T) {
	c := debugCache()
	c.Put("10__one", 10, 2, 3)
	c.Put("10__two", true, 2)
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:   "update existing item",
			fields: c,
			args: args{
				key:   "10__one",
				value: 20,
			},
			want:    20,
			wantErr: false,
		},
		{
			name:   "update existing item",
			fields: c,
			args: args{
				key:   "10__two",
				value: 30,
			},
			want:    30,
			wantErr: false,
		},
		{
			name:   "update non existing item",
			fields: c,
			args: args{
				key:   "10__three",
				value: []byte("hello world"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.Update(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("cache.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if v, _ := c.Get(tt.args.key); v != tt.want {
				t.Errorf("cache.Update() value = %v, want %v", v, tt.want)
			}

		})
	}
}

func Test_cache_Patch(t *testing.T) {
	c := debugCache()
	c.Put("10__one", 10, 2, 3)
	c.Put("10__two", true, 2)
	type args struct {
		key   string
		value interface{}
		tags  []uint16
	}
	tests := []struct {
		name   string
		fields CacheStore
		args   args
	}{
		{
			name:   "patch non existing item",
			fields: c,
			args: args{
				key:   "10__three",
				value: 333,
			},
		},
		{
			name:   "patch existing item",
			fields: c,
			args: args{
				key:   "10__two",
				value: 777,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Patch(tt.args.key, tt.args.value, tt.args.tags...)
			if v, err := c.GetItem(tt.args.key); err != nil || v.Value != tt.args.value {
				t.Errorf("cache.Update() value = %v, want %v", v.Value, tt.args.value)
			}
		})
	}
}

func Test_cache_Delete(t *testing.T) {
	c := debugCache()
	c.Put("10__one", 10, 2, 3)
	c.Put("10__two", true, 2)
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := c.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("cache.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cache_Clear(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		capacity  uint64
		incr      uint64
		items     map[string]Item
		defaultlt time.Duration
		auditor   Auditor
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				mu:        tt.fields.mu,
				capacity:  tt.fields.capacity,
				incr:      tt.fields.incr,
				items:     tt.fields.items,
				defaultlt: tt.fields.defaultlt,
				auditor:   tt.fields.auditor,
			}
			c.Clear()
		})
	}
}

func Test_cache_List(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		capacity  uint64
		incr      uint64
		items     map[string]Item
		defaultlt time.Duration
		auditor   Auditor
	}
	tests := []struct {
		name   string
		fields fields
		want   []Item
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				mu:        tt.fields.mu,
				capacity:  tt.fields.capacity,
				incr:      tt.fields.incr,
				items:     tt.fields.items,
				defaultlt: tt.fields.defaultlt,
				auditor:   tt.fields.auditor,
			}
			if got := c.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_Filter(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		capacity  uint64
		incr      uint64
		items     map[string]Item
		defaultlt time.Duration
		auditor   Auditor
	}
	type args struct {
		f func(i Item) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Item
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				mu:        tt.fields.mu,
				capacity:  tt.fields.capacity,
				incr:      tt.fields.incr,
				items:     tt.fields.items,
				defaultlt: tt.fields.defaultlt,
				auditor:   tt.fields.auditor,
			}
			if got := c.Filter(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_ForEach(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		capacity  uint64
		incr      uint64
		items     map[string]Item
		defaultlt time.Duration
		auditor   Auditor
	}
	type args struct {
		f func(i Item)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				mu:        tt.fields.mu,
				capacity:  tt.fields.capacity,
				incr:      tt.fields.incr,
				items:     tt.fields.items,
				defaultlt: tt.fields.defaultlt,
				auditor:   tt.fields.auditor,
			}
			c.ForEach(tt.args.f)
		})
	}
}

func Test_cache_ListValues(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		capacity  uint64
		incr      uint64
		items     map[string]Item
		defaultlt time.Duration
		auditor   Auditor
	}
	tests := []struct {
		name   string
		fields fields
		want   []interface{}
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				mu:        tt.fields.mu,
				capacity:  tt.fields.capacity,
				incr:      tt.fields.incr,
				items:     tt.fields.items,
				defaultlt: tt.fields.defaultlt,
				auditor:   tt.fields.auditor,
			}
			if got := c.ListValues(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.ListValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_ListKeys(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		capacity  uint64
		incr      uint64
		items     map[string]Item
		defaultlt time.Duration
		auditor   Auditor
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				mu:        tt.fields.mu,
				capacity:  tt.fields.capacity,
				incr:      tt.fields.incr,
				items:     tt.fields.items,
				defaultlt: tt.fields.defaultlt,
				auditor:   tt.fields.auditor,
			}
			if got := c.ListKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.ListKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_ExtendLifetime(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		capacity  uint64
		incr      uint64
		items     map[string]Item
		defaultlt time.Duration
		auditor   Auditor
	}
	type args struct {
		key string
		dur time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				mu:        tt.fields.mu,
				capacity:  tt.fields.capacity,
				incr:      tt.fields.incr,
				items:     tt.fields.items,
				defaultlt: tt.fields.defaultlt,
				auditor:   tt.fields.auditor,
			}
			if err := c.ExtendLifetime(tt.args.key, tt.args.dur); (err != nil) != tt.wantErr {
				t.Errorf("cache.ExtendLifetime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cache_Immortalize(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		capacity  uint64
		incr      uint64
		items     map[string]Item
		defaultlt time.Duration
		auditor   Auditor
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				mu:        tt.fields.mu,
				capacity:  tt.fields.capacity,
				incr:      tt.fields.incr,
				items:     tt.fields.items,
				defaultlt: tt.fields.defaultlt,
				auditor:   tt.fields.auditor,
			}
			if err := c.Immortalize(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("cache.Immortalize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
