package multitemplate

import (
	"html/template"
	"path/filepath"

	"github.com/gin-gonic/gin/render"
)

type Render struct {
	templates map[string]*template.Template
	left      string
	right     string
}

var _ render.HTMLRender = Render{}

func New() Render {
	return Render{
		templates: make(map[string]*template.Template),
		left:      "{{",
		right:     "}}",
	}
}

func (r *Render) SetDelimiter(left, right string) {
	r.left = left
	r.right = right
}

func (r Render) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	r.templates[name] = tmpl
}

func (r Render) AddFromFiles(name string, files ...string) *template.Template {
	tmpl := template.Must(template.New(filepath.Base(files[0])).Delims(r.left, r.right).ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

func (r Render) AddFromGlob(name, glob string) *template.Template {
	filenames, err := filepath.Glob(glob)
	if err != nil {
		panic(err)
	}

	return r.AddFromFiles(name, filenames...)
}

func (r *Render) AddFromString(name, templateString string) *template.Template {
	tmpl := template.Must(template.New("").Delims(r.left, r.right).Parse(templateString))
	r.Add(name, tmpl)
	return tmpl
}

func (r Render) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.templates[name],
		Data:     data,
	}
}
