package filter

import (
	"encoding/json"

	"entgo.io/ent/schema"
)

type Annotation struct {
	FilterFields []*Opt
	Sort         *Opt
	ReturnTotal  *Opt
	Page         *Opt
	ItemsPerPage *Opt
	NoPagination bool
}

// Merge implements ent.Merger interface.
func (a Annotation) Merge(o schema.Annotation) schema.Annotation {
	var ant Annotation
	switch o := o.(type) {
	case Annotation:
		ant = o
	case *Annotation:
		if o != nil {
			ant = *o
		}
	default:
		return a
	}

	if ant.FilterFields != nil {
		if a.FilterFields == nil {
			a.FilterFields = make([]*Opt, 0)
		}
		a.FilterFields = append(a.FilterFields, ant.FilterFields...)
	}

	if ant.Sort != nil && a.Sort == nil {
		a.Sort = ant.Sort
	}

	if ant.ReturnTotal != nil && a.ReturnTotal == nil {
		a.ReturnTotal = ant.ReturnTotal
	}

	if ant.Page != nil && a.Page == nil {
		a.Page = ant.Page
	}

	if ant.ItemsPerPage != nil && a.ItemsPerPage == nil {
		a.ItemsPerPage = ant.ItemsPerPage
	}

	if ant.NoPagination {
		a.NoPagination = ant.NoPagination
	}

	return a
}

func (Annotation) Name() string {
	return "ListOperations"
}

func WithFieldFilter(fields ...string) Annotation {
	fs := make([]*Opt, 0)
	for _, f := range fields {
		fs = append(fs, &Opt{
			In:   "query",
			Name: f,
		})
	}

	return Annotation{
		FilterFields: fs,
	}
}

// Decode from ent.
func (a *Annotation) Decode(o interface{}) error {
	buf, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, a)
}

func WithFilterField(name string, opts ...OptConfig) MutatorOpt {
	return func(a *Annotation) {
		o := &Opt{
			Name: name,
			In:   "query",
		}
		for _, opt := range opts {
			opt(o)
		}

		a.FilterFields = append(a.FilterFields, o)
	}
}

func WithSort(opts ...OptConfig) MutatorOpt {
	return func(a *Annotation) {
		a.Sort = &Opt{
			Name: "sort",
			In:   "query",
		}

		for _, opt := range opts {
			opt(a.Sort)
		}
	}
}

func WithNoPagination() MutatorOpt {
	return func(a *Annotation) {
		a.NoPagination = true
	}
}

func WithPage(opts ...OptConfig) MutatorOpt {
	return func(a *Annotation) {
		a.Page = &Opt{
			Name: "page",
			In:   "query",
		}
		for _, opt := range opts {
			opt(a.Page)
		}
	}
}

func WithItemsPerPage(opts ...OptConfig) MutatorOpt {
	return func(a *Annotation) {
		a.ItemsPerPage = &Opt{
			Name: "itemsPerPage",
			In:   "query",
		}

		for _, opt := range opts {
			opt(a.ItemsPerPage)
		}
	}
}

func WithReturnTotal(opts ...OptConfig) MutatorOpt {
	return func(a *Annotation) {
		a.ReturnTotal = &Opt{
			Name: "total",
			In:   "header",
		}

		for _, opt := range opts {
			opt(a.ReturnTotal)
		}
	}
}
