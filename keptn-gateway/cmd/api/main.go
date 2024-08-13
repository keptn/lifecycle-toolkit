package main

import (
	"keptn.sh/keptn-gateway/internal/gateway"
)

func main() {
	gw := gateway.NewGateway()
	gw.Serve()
}
