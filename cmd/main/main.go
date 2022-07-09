//main.go
package main

import "web_app/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
