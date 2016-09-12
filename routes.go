package main

import (
	m "cli-client/models"
	c "cli-client/controllers"
	//"github.com/betacraft/yaag/middleware"
)

var routes = m.Routes{
	m.Route{"Ping","GET","/ping",c.Ping,},

	m.Route{"Index","GET","/app",c.Index,},
	m.Route{"CreateApp","POST","/app",c.CreateApp,},

	m.Route{"RegisterDevice","POST","/device",c.RegisterDevice,},
}