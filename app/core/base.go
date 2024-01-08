package core

import "fastProxy/app/services"

func Start() {
	services.StartHttpServer()
}
