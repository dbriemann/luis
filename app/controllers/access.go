package controllers

import "github.com/revel/revel"

type Access struct {
	*revel.Controller
}

func (c Access) Login() revel.Result {
	return c.Render()
}
