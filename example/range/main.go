package main

import (
	"fmt"
	"github.com/NICEXAI/ghost"
	"os"
	"path"
	"strings"
)

type UserList []User

type User struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func main() {
	currentPath, _ := os.Getwd()

	options := make(map[string]interface{})
	//dataList := make([]map[string]interface{}, 0)
	//
	//userList := UserList{
	//	{
	//		Name: "JieKe",
	//		Age: 12,
	//	},
	//	{
	//		Name: "Mari",
	//		Age: 15,
	//	},
	//}
	//
	//bUsers, _ := json.Marshal(userList)
	//_ = json.Unmarshal(bUsers, &dataList)
	//
	//options["data_list"] = dataList

	options["data_list"] = []string{"jieKe", "Mari"}

	originFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/range/origin")
	targetFolder := path.Join(strings.ReplaceAll(currentPath, `\`, `/`), "example/range/dist")

	if err := ghost.ParseAll(originFolder, targetFolder, options); err != nil {
		fmt.Println(err)
	}
}
