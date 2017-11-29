package pac

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arces-wot/SEPA-Go/sepa"
	"io"
	"strings"
	"text/template"
)

type Variable struct {
	Type  string
	Value string
}

type SparqlType struct {
	Sparql         string
	ForcedBindings map[string]Variable
}

type profileInternal struct {
	Parameters sepa.Configuration
	Namespaces map[string]string
	Updates    map[string]SparqlType
	Queries    map[string]SparqlType
}

type Profile struct {
	sepa.Configuration
	namespaces map[string]string
	queries map[string]template.Template
	updates map[string]template.Template
}

func newProfile(config sepa.Configuration) Profile {
	qt := make(map[string]template.Template)
	ut := make(map[string]template.Template)
	ns := make(map[string]string)
	p := Profile{config,ns, qt, ut}
	return p
}

const DELIMITER = "??"
var baseTmpl = template.New("___BASE___").Delims(DELIMITER,DELIMITER)

func ProfileFromJson(reader io.Reader) (Profile, error) {
	config := sepa.DefaultConfig()
	profile := &profileInternal{Parameters: config}


	if err := json.NewDecoder(reader).Decode(profile); err != nil {
		return Profile{}, err
	}

	return processRawProfile(*profile)
}

func (p *Profile) AddNamespace(ns string) {
	//TODO: Add namespace handling
	panic("Not implemented")
}

func (p *Profile) AddQuery(name string, query string) error {
	return addTemplate(p.queries, name, query)
}

func (p *Profile) AddUpdate(name string, update string) error {
	return addTemplate(p.updates, name, update)
}

func (p *Profile) ForgetQuery(name string) {
	deleteTemplate(p.queries, name)
}

func (p *Profile) ForgetUpdate(name string) {
	deleteTemplate(p.updates, name)
}

func (p *Profile) GetQuery(id string, data interface{}) (result string, err error) {
	return get(p.queries, id, data)
}

func (p *Profile) GetUpdate(id string, data interface{}) (string, error) {
	return get(p.updates, id, data)
}

func (p *Profile) ContainsQuery(name string) bool {
	_, ok := p.queries[name]
	return ok
}

func (p *Profile) ContainsUpdate(name string) bool {
	_, ok := p.updates[name]
	return ok
}

func get(templates map[string]template.Template, name string, data interface{}) (result string, err error) {
	resWriter := bytes.NewBufferString(result)

	if val, ok := templates[name]; ok {
		if err = val.Execute(resWriter, data); err == nil {
			result = resWriter.String()
		}
	} else {
		err = fmt.Errorf("no sparql found for %s", name)
	}

	return
}

func addTemplate(templates map[string]template.Template, name string, sparql string) error {

	if _, ok := templates[name]; ok {
		return errors.New("query already defined")
	}

	var tmp *template.Template
	var e error

	if tmp, e = baseTmpl.New(name).Parse(sparql); e != nil {
		return e
	}

	templates[name] = *tmp

	return nil
}

func deleteTemplate(templates map[string]template.Template, s string) {
	delete(templates, s)
}

func processRawProfile(p profileInternal) (Profile,error) {
	result := newProfile(p.Parameters)
	for query, data := range p.Queries {
		tmpl := transformInGoTemplate(data.Sparql, data.ForcedBindings)
		if err := result.AddQuery(query, tmpl); err != nil {
			return result, err
		}
	}

	for update, data := range p.Updates {
		tmpl := transformInGoTemplate(data.Sparql, data.ForcedBindings)
		if err := result.AddUpdate(update, tmpl); err != nil {
			return result,err
		}
	}
	return result,nil
}

func transformInGoTemplate(sparql string, binding map[string]Variable) string {
	template := sparql
	for name, bind := range binding {
		var templater *strings.Replacer
		if bind.Value != "" {
			templater = strings.NewReplacer("?"+name, bind.Value)
		} else {
			templater = strings.NewReplacer("?"+name, DELIMITER+"."+strings.Title(name)+DELIMITER)
		}
		template = templater.Replace(template)
	}
	return template
}
