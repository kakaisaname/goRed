package main

import (
	"encoding/json"
	"fmt"
	"github.com/kakaisaname/goRed/services"
)

func main() {

	d, e := json.Marshal(&services.AccountTransferDTO{})
	fmt.Println(e)
	fmt.Println(string(d))
}
