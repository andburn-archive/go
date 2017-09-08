package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var moduleName string = "moduleName"
var baseUrl string = "fullURL"
var moduleStr string = "module"
var moduleNum int = 6
var fileExt string = "pdf"

func main() {
	for i := 2; i <= moduleNum; i++ {
		for j := 1; j <= 100; j++ {
			fname := fmt.Sprintf("m%d_u%d.%s", i, j, fileExt)
			url := fmt.Sprintf("%s/%s/%s%d/%s", baseUrl, moduleName, moduleStr, i, fname)
			fmt.Println(url)

			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				break
			}

			fmt.Println(url)

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}

			outDir := "\\" + moduleName
			os.Mkdir(outDir, os.ModePerm)
			err = ioutil.WriteFile(outDir+"\\"+fname, data, 0666)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
