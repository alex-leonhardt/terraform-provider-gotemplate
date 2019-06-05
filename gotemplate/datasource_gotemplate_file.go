package gotemplate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/hashicorp/terraform/helper/schema"
)

// ---------------------------------------

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}

type templateRenderError error

func renderFile(d *schema.ResourceData) (string, error) {

	var err error

	tf := template.FuncMap{
		"isInt": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
				return true
			default:
				return false
			}
		},
		"isString": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.String:
				return true
			default:
				return false
			}
		},
		"isSlice": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Slice:
				return true
			default:
				return false
			}
		},
		"isArray": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Array:
				return true
			default:
				return false
			}
		},
		"isMap": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Map:
				return true
			default:
				return false
			}
		},
	}

	for k, v := range sprig.TxtFuncMap() {
		tf[k] = v
	}

	var data string // data from tf
	data = d.Get("data").(string)

	// unmarshal json from data into m
	var m = make(map[string]interface{}) // unmarshal data into m
	if err = json.Unmarshal([]byte(data), &m); err != nil {
		return "", templateRenderError(fmt.Errorf("failed to render %v", err))
	}

	templateFile := d.Get("template").(string)
	baseName := path.Base(templateFile)
	tt, err := template.New(baseName).Funcs(tf).ParseFiles(templateFile)
	if err != nil {
		return "", templateRenderError(fmt.Errorf("failed to render %v", err))
	}

	var contents bytes.Buffer // io.writer for template.Execute
	if tt != nil {
		err = tt.Execute(&contents, m)
		if err != nil {
			return "", templateRenderError(fmt.Errorf("failed to render %v", err))
		}
	} else {
		return "", templateRenderError(fmt.Errorf("error: %v", err))
	}

	return contents.String(), nil
}

func dataSourceFileRead(d *schema.ResourceData, meta interface{}) error {
	rendered, err := renderFile(d)
	if err != nil {
		return err
	}
	d.Set("rendered", rendered)
	d.SetId(hash(rendered))
	return nil
}
