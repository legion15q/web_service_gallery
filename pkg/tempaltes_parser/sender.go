package sender

import (
	"bytes"
	"text/template"

	"github.com/sirupsen/logrus"
)

type Template struct {
	Template *template.Template
}

type TemplateParser interface {
	GenerateBodyFromHTML(data interface{}) (string, error)
}

func NewTemplateParser(templateFileName string) *Template {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		logrus.Error("failed to parse file", templateFileName)
		return &Template{}
	}
	return &Template{t}
}
func (e *Template) GenerateBodyFromHTML(data interface{}) string {

	buf := new(bytes.Buffer)
	if err := e.Template.Execute(buf, data); err != nil {
		logrus.Error("failed to execute template: ", err)
		return ""
	}
	return buf.String()
}
