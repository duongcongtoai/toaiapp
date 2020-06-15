package registry

import (
	"log"

	"github.com/labstack/echo/v4"
)

var instance *Registry

type Component interface {
	SetRegistry(*Registry)
	GetWeight() int
	SetupEcho(*echo.Echo) error
	Shutdown() error
	GetName() string
	SetupFromYaml(string, bool) error
}
type Registry struct {
	clist  []Component
	cnames map[string]Component
	debug  bool
}

func Instance() *Registry {
	if instance != nil {
		return instance
	}

	instance = &Registry{
		cnames: make(map[string]Component),
	}
	return instance
}

func (r *Registry) SetupEcho(e *echo.Echo) error {
	for _, c := range r.clist {
		if err := c.SetupEcho(e); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) SetupFromYaml(configFile string, debug bool) error {
	r.debug = debug
	for _, c := range r.clist {
		if r.debug {
			log.Printf("SetupFromYaml (%d): %s\n", c.GetWeight(), c.GetName())
		}
		if err := c.SetupFromYaml(configFile, debug); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) Shutdown() error {
	i := len(r.clist) - 1
	for true {
		c := r.clist[i]
		if r.debug {
			log.Printf("Shutdown (%d): %s\n", c.GetWeight(), c.GetName())
		}

		if err := c.Shutdown(); err != nil {
			return err
		}
		i = i - 1
		if i < 0 {
			break
		}
	}
	return nil
}

func (r *Registry) Register(name string, c Component) {
	c.SetRegistry(r)
	r.cnames[name] = c

	idx := -1
	found := false
	for i, lcom := range r.clist {
		idx = i
		if lcom.GetWeight() < c.GetWeight() {
			continue
		} else {
			found = true
			break
		}
	}
	if !found {
		idx = idx + 1
	}

	if idx == 0 {
		r.clist = append([]Component{c}, r.clist...)
	} else if idx == len(r.clist) {
		r.clist = append(r.clist, c)
	} else {
		after := r.clist[idx:]
		r.clist = append(r.clist[:idx], c)
		r.clist = append(r.clist, after...)
	}
}
