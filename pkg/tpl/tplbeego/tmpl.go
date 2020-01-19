package tplbeego

import (
	"bytes"
	"errors"
	"html/template"
)

type Tmpl struct {
	// context data
	Data map[interface{}]interface{}

	debug   bool
	tplPath string

	// template data
	TplName        string
	ViewPath       string
	Layout         string
	LayoutSections map[string]string // the key is the section name and the value is the template name
	TplPrefix      string
	TplExt         string
}

// Init generates default values of controller operations.
func (c *Tmpl) Init(tplExt, viewPath string) {
	c.Data = make(map[interface{}]interface{}, 0)
	c.Layout = ""
	c.TplName = ""
	c.TplExt = tplExt
	c.ViewPath = viewPath
}

func (c *Tmpl) SetTplPath(tplPath string) {
	c.tplPath = tplPath
}

// RenderBytes returns the bytes of rendered template string. Do not send out response.
func (c *Tmpl) RenderBytes() ([]byte, error) {
	buf, err := c.renderTemplate()
	//if the controller has set layout, then first get the tplName's content set the content to the layout
	if err == nil && c.Layout != "" {
		c.Data["LayoutContent"] = template.HTML(buf.String())

		if c.LayoutSections != nil {
			for sectionName, sectionTpl := range c.LayoutSections {
				if sectionTpl == "" {
					c.Data[sectionName] = ""
					continue
				}
				buf.Reset()
				err = ExecuteViewPathTemplate(&buf, sectionTpl, c.viewPath(), c.Data)
				if err != nil {
					return nil, err
				}
				c.Data[sectionName] = template.HTML(buf.String())
			}
		}

		buf.Reset()
		ExecuteViewPathTemplate(&buf, c.Layout, c.viewPath(), c.Data)
	}
	return buf.Bytes(), err
}

func (c *Tmpl) renderTemplate() (bytes.Buffer, error) {
	var buf bytes.Buffer
	if c.tplPath == "" {
		return buf, errors.New("tpl path is empty")
	}
	if c.TplName == "" {
		c.TplName = c.tplPath + "." + c.TplExt
	}
	if c.TplPrefix != "" {
		c.TplName = c.TplPrefix + c.TplName
	}
	if Config().Debug {
		buildFiles := []string{c.TplName}
		if c.Layout != "" {
			buildFiles = append(buildFiles, c.Layout)
			if c.LayoutSections != nil {
				for _, sectionTpl := range c.LayoutSections {
					if sectionTpl == "" {
						continue
					}
					buildFiles = append(buildFiles, sectionTpl)
				}
			}
		}
		BuildTemplate(c.viewPath(), buildFiles...)
	}
	return buf, ExecuteViewPathTemplate(&buf, c.TplName, c.viewPath(), c.Data)
}

func (c *Tmpl) viewPath() string {
	return c.ViewPath
}
