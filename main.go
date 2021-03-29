package main

import (
	"fmt"
	"go-practice/iniconfig/ini"
	"io/ioutil"
	"log"
)

func main()  {
	data, _ := ioutil.ReadFile("./ini/config.ini")

	conf := &ini.Config{}

	err := ini.UnMarshal(data, conf)
	if err != nil {
		log.Fatal("unmarshal failed: ", err)
	}


	fmt.Println(conf)

}
