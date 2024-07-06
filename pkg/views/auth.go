package views

import (
	"html/template"
)

func Login() *template.Template {
	tmpl, err := template.ParseFiles(
		"templates/base.tmpl",
		"templates/login.tmpl",
	)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func AdminLogin() *template.Template {
	tmpl, err := template.ParseFiles(
		"templates/base.tmpl",
		"templates/adminLogin.tmpl",
	)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func Signup() *template.Template {
	tmpl, err := template.ParseFiles(
		"templates/base.tmpl",
		"templates/signup.tmpl",
	)
	if err != nil {
		panic(err)
	}
	return tmpl
}
