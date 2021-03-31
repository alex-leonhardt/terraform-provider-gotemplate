package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"text/template"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/giantswarm/k8scloudconfig/v10/pkg/ignition"
	"github.com/giantswarm/microerror"
)

const (
	nestedTemplatesDir = "files"
)

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}

// files is map[string]string (k: filename, v: contents) for files that are fetched from disk
// and then filled with data.
type files map[string]string

// RenderNestedTemplates walks over templatesDir and parses all regular files with
// text/template. Parsed templates are then rendered with ctx, base64 encoded
// and added to returned files.
func renderNestedTemplates(templatesDir string, ctx interface{}) (files, error) {
	files := files{}

	err := filepath.Walk(templatesDir, func(path string, f os.FileInfo, err error) error {
		if f.Mode().IsRegular() {
			tmpl, err := template.ParseFiles(path)
			if err != nil {
				return microerror.Maskf(err, "failed to parse file %#q", path)
			}
			var data bytes.Buffer
			tmpl.Execute(&data, ctx)

			relativePath, err := filepath.Rel(templatesDir, path)
			if err != nil {
				return microerror.Mask(err)
			}
			files[relativePath] = base64.StdEncoding.EncodeToString(data.Bytes())
		}
		return nil
	})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return files, nil
}

func renderMainTemplate(d *schema.ResourceData) (string, error) {

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
	templatesDir := filepath.Join(filepath.Dir(templateFile), nestedTemplatesDir)

	// render nested templates
	m["Files"], err = renderNestedTemplates(templatesDir, m)
	if err != nil {
		panic(err)
	}

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

	// generate ignition if requested
	is_ignition := d.Get("is_ignition").(bool)
	if is_ignition {
		ignition_content, err := ignition.ConvertTemplatetoJSON(contents.Bytes())
		if err != nil {
			panic(err)
		}

		return string(ignition_content), nil
	}

	return contents.String(), nil
}

func dataSourceFileRead(d *schema.ResourceData, meta interface{}) error {
	rendered, err := renderMainTemplate(d)
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
			"is_ignition": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				Description: "return ignition in rendered",
				Default:     false,
			},
			"rendered": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered template",
			},
		},
	}
}
