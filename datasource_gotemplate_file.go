package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"text/template"

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

	var data string // data from tf
	data = d.Get("data").(string)

	// unmarshal json from data into m
	var m = make(map[string]interface{}) // unmarshal data into m
	if err = json.Unmarshal([]byte(data), &m); err != nil {
		panic(err)
	}

	templateFile := d.Get("template").(string)
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}

	var contents bytes.Buffer // io.writer for template.Execute
	tt := t.Funcs(tf)
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

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFileRead,

		Schema: map[string]*schema.Schema{
			"template": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "path to go template file",
			},
			"data": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "variables to substitute",
				ValidateFunc: nil,
			},
			"rendered": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered template",
			},
		},
	}
}
