package torrenttools

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/marksamman/bencode"
)

func DecodeTorrent(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dict, err := bencode.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	for k := range dict {

		if reflect.TypeOf(dict[k]).Kind() == reflect.Map {
			subMap := dict[k].(map[string]interface{})

			for n := range subMap {
				//fmt.Println(newShit[n])
				//if newShit[n]  {
				fmt.Println(n, "(map):")

				if n != "pieces" {
					fmt.Println(subMap[n])
				}

				//}
			}

		} else if reflect.TypeOf(dict[k]).Kind() == reflect.String {
			fmt.Println(k, ": ", dict[k])
			//fmt.Println(dict[k])
		} else if reflect.TypeOf(dict[k]).Kind() == reflect.Slice {
			fmt.Println(k, "(slice):")
			subSlice := dict[k].([]interface{})
			for n := range subSlice {
				fmt.Println(n)
			}

			fmt.Println(dict[k])
			fmt.Println("--------------")

		} else {
			fmt.Println("-------------- Type error:", reflect.TypeOf(dict[k]).Kind())

			fmt.Println(k, ": ", dict[k])
			fmt.Println("-------------- End error --------------")
			//panic(1)
		}

	}

	fmt.Printf("string: %v\n", dict["comment"])

}
