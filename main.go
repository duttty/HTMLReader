package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var (
	ctsMap   = make(map[string][]string)
	bookList = make([]string, 0, 20)
	rootPath = getAbs()
	str      = `
				<!DOCTYPE html>
				<html>
					<body>		
						<ul>
							{{range $k,$v := .}}
							<li>
								<a href="{{$v}}">{{$v}}</a>
							</li>
							{{end}}
						</ul>
					</body>
				</html>
			`
	right = `
		<div style="
			position: fixed;
			right: 0;
			top: 50%;
			margin-top: -10px;
		"><a href="{{.}}" style="font-weight: bold;">→</a>
		</div>
	`
)

func main() {
	initContents()
	r := httprouter.New()
	//查看书单
	r.GET("/", rootHandler)
	//查看目录
	r.GET("/:path/", pathHandler)
	//开始阅读
	r.GET("/:path/:name", bookHandler)

	log.Fatal(http.ListenAndServe(":8089", r))
}

//主目录，显示有哪些书
func rootHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t, err := template.New("test").Parse(str)
	if err != nil {
		log.Fatalln(err)
	}
	t.Execute(w, bookList)
}

//进入图书
func pathHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//获取cookie
	cookie, err := r.Cookie("read-" + ps.ByName("path"))
	path := ps.ByName("path")
	//不存在则展示目录
	if err != nil {
		cts := ctsMap[path]
		t, err := template.New("test").Parse(str)
		if err != nil {
			log.Fatalln(err)
		}
		t.Execute(w, cts)
		return
	}
	//跳转至上次阅读
	http.Redirect(w, r, cookie.Value, http.StatusTemporaryRedirect)
}

//读书
func bookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := ps.ByName("path")
	name := ps.ByName("name")
	//写入cookie
	http.SetCookie(w, &http.Cookie{
		Name:  fmt.Sprintf("read-%s", path),
		Value: url.QueryEscape(name),
	})
	fp := fmt.Sprintf("%s/%s/%s", rootPath, path, name)
	//读入文件
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Printf("read %s err:\n%s", fp, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	temStr := string(b)
	idx := strings.LastIndex(temStr, "</body>")
	temStr = fmt.Sprintf("%s%s%s", temStr[:idx], right, temStr[idx:])
	t, err := template.New("right").Parse(temStr)
	if err != nil {
		fmt.Println("tem err : ", err)
		return
	}
	this := 0
	list := ctsMap[path]
	for k, v := range list {
		if v == name {
			this = k
			break
		}
	}
	if this == len(list)-1 {
		t.Execute(w, name)
	} else {
		t.Execute(w, list[this+1])
	}

}

//获取绝对路径
func getAbs() string {
	// 获取可执行文件相对于当前工作目录的相对路径
	rel := filepath.Dir(os.Args[0])
	// 根据相对路径获取可执行文件的绝对路径
	abs, err := filepath.Abs(rel)
	if err != nil {
		log.Fatalln("get abs err : ", err)
	}

	return abs
}

//获取目录以及删除非html文件
func getList(dir string) []string {
	flag := false
	fInfoL, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("getList err:\ndir=%s\nerr:%s", dir, err)
	}
	fList := make([]string, 0, 200)
	for _, v := range fInfoL {
		name := v.Name()
		if strings.HasSuffix(name, ".html") {
			fList = append(fList, name)
		} else {
			if !flag {
				fmt.Printf("是否删除%s?按Y确定", dir)
				var f string
				fmt.Scanln(&f)
				if f == "Y" {
					flag = true
				} else {
					log.Fatalln("目录存在重要文件", dir)
				}
			}
			err := os.RemoveAll(dir + "/" + name)
			if err != nil {
				fmt.Println("delete err:", name)
			}
		}
	}
	return fList
}

//初始化目录信息
func initContents() {
	//获取文件信息
	fInfo, err := ioutil.ReadDir(rootPath)
	if err != nil {
		log.Fatalln("init cts err : ", err)
	}
	for _, v := range fInfo {
		if v.IsDir() {
			bookList = append(bookList, v.Name())
			ctsMap[v.Name()] = getList(rootPath + "/" + v.Name())
		}
	}
}