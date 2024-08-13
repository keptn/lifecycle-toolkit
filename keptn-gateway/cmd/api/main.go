package main

import (
	"fmt"
	"keptn.sh/keptn-gateway/internal/gateway"
)

func main() {
	gw := gateway.NewGateway()

	err := gw.Server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
