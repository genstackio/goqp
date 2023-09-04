package goqp

import (
	"encoding/json"
	"errors"
	"net/url"
	"testing"
)

type MyStruct struct {
	Name   string
	Owner  string
	Email  string
	Extras string
}

func TestNewQueryParser(t *testing.T) {
	t.Run("string param with no default value return empty string when no present", func(t *testing.T) {
		q := url.Values{}
		d := MyStruct{}
		err := NewQueryParser(&q, &d).String("n", "", func(v string, d *MyStruct) {
			d.Name = v
		}).Parse(func(d *MyStruct) error {
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if d.Name != "" {
			t.Fatalf("bad value: %s instead of %s", d.Name, "hello")
		}
	})
	t.Run("string param with default value return default value when no present", func(t *testing.T) {
		q := url.Values{}
		d := MyStruct{}
		err := NewQueryParser(&q, &d).String("n", "hello", func(v string, d *MyStruct) {
			d.Name = v
		}).Parse(func(d *MyStruct) error {
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if d.Name != "hello" {
			t.Fatalf("bad value: %s instead of %s", d.Name, "hello")
		}
	})
	t.Run("string param with no default value but present and non empty return non empty value", func(t *testing.T) {
		q := url.Values{}
		q.Set("n", "bla")
		d := MyStruct{}
		err := NewQueryParser(&q, &d).String("n", "", func(v string, d *MyStruct) {
			d.Name = v
		}).Parse(func(d *MyStruct) error {
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if d.Name != "bla" {
			t.Fatalf("bad value: %s instead of %s", d.Name, "bla")
		}
	})
	t.Run("error fn param with failing status return error and do not populate remaining fields in struct", func(t *testing.T) {
		q := url.Values{}
		q.Set("n", "bla")
		q.Set("o", "me")
		q.Set("e", "the@email.com")
		d := MyStruct{}
		err := NewQueryParser(&q, &d).
			String("n", "", func(v string, d *MyStruct) {
				d.Name = v
			}).
			String("o", "", func(v string, d *MyStruct) {
				d.Owner = v
			}).
			ErrorFn("k", func(v string, d *MyStruct) error {
				return errors.New("there was an error")
			}).
			String("e", "", func(v string, d *MyStruct) {
				d.Email = v
			}).
			Parse(func(d *MyStruct) error {
				return nil
			})
		if err == nil {
			t.Fatalf("expected error not raised")
		}
		if d.Name != "bla" {
			t.Fatalf("bad value: %s instead of %s", d.Name, "bla")
		}
		if d.Owner != "me" {
			t.Fatalf("bad value: %s instead of %s", d.Owner, "me")
		}
		if d.Email != "" {
			t.Fatalf("bad value: %s instead of %s", d.Email, "")
		}
	})
	t.Run("extras param return params that were not registered", func(t *testing.T) {
		q := url.Values{}
		q.Set("n", "bla")
		q.Set("o", "me")
		q.Set("e", "the@email.com")
		d := MyStruct{}
		err := NewQueryParser(&q, &d).
			String("n", "", func(v string, d *MyStruct) {
				d.Name = v
			}).
			String("o", "", func(v string, d *MyStruct) {
				d.Owner = v
			}).
			Extras(func(extras map[string]string, d *MyStruct) {
				ex, _ := json.Marshal(extras)
				d.Extras = string(ex)
			}).
			Parse(func(d *MyStruct) error {
				return nil
			})
		if err != nil {
			t.Fatal(err)
		}
		if d.Name != "bla" {
			t.Fatalf("bad value: %s instead of %s", d.Name, "bla")
		}
		if d.Owner != "me" {
			t.Fatalf("bad value: %s instead of %s", d.Owner, "me")
		}
		if d.Extras != "{\"e\":\"the@email.com\"}" {
			t.Fatalf("bad value: %s instead of %s", d.Extras, "{\"e\":\"the@email.com\"}")
		}
	})
	t.Run("extras param return params that were not registered #2", func(t *testing.T) {
		q := url.Values{}
		q.Set("n", "bla")
		q.Set("o", "me")
		q.Set("e", "the@email.com")
		d := MyStruct{}
		err := NewQueryParser(&q, &d).
			String("o", "", func(v string, d *MyStruct) {
				d.Owner = v
			}).
			Extras(func(extras map[string]string, d *MyStruct) {
				ex, _ := json.Marshal(extras)
				d.Extras = string(ex)
			}).
			Parse(func(d *MyStruct) error {
				return nil
			})
		if err != nil {
			t.Fatal(err)
		}
		if d.Name != "" {
			t.Fatalf("bad value: %s instead of %s", d.Name, "")
		}
		if d.Owner != "me" {
			t.Fatalf("bad value: %s instead of %s", d.Owner, "me")
		}
		if d.Extras != "{\"e\":\"the@email.com\",\"n\":\"bla\"}" {
			t.Fatalf("bad value: %s instead of %s", d.Extras, "{\"e\":\"the@email.com\",\"n\":\"bla\"}")
		}
	})
}
