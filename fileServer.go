package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"
)


func main() {
	fmt.Println("hello。。。。")
	http.Handle("/", http.FileServer(http.Dir(`/root/极客时间`)))
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/test", test)
	http.ListenAndServe(":1111", nil)
}


func upload(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	if r.Method == "GET"{
		time := time.Now().Unix()
		h := md5.New()
		h.Write([]byte(strconv.FormatInt(time,10)))
		token := hex.EncodeToString(h.Sum(nil))
		t, _ := template.ParseFiles("./view/upload.ctpl")
		t.Execute(w, token)
	}else if r.Method == "POST"{        //把上传的文件存储在内存和临时文件中
		w.Header().Set("content-type","text/html; charset=utf-8")
		r.ParseMultipartForm(32 << 20)        //获取文件句柄，然后对文件进行存储等处理
		file, handler, err := r.FormFile("uploadfile")
		if err != nil{
			log.Println("form file err: ", err)
			return
		}
		defer file.Close()
		names := strings.Split(handler.Filename,".")
		if len(names) != 2 {
			fmt.Fprintf(w, "文件名错误, 除了后缀, 不能有'.'")
			fmt.Fprintf(w, `<a href="http://212.64.16.41:1111">返回</a>`)
			return
		}
		if names[1] != "zip" {
			fmt.Fprintf(w, "只支持zip文件")
			fmt.Fprintf(w, `<a href="http://212.64.16.41:1111">返回</a>`)
			return
		}
		filepathNames,_ := ioutil.ReadDir("/root/极客时间")
		for _,v := range filepathNames {
			if v.Name() == names[0] {
				fmt.Fprintf(w, "文件已经存在")
				fmt.Fprintf(w, `<a href="http://212.64.16.41:1111">返回</a>`)
				return
			}
		}
		//创建上传的目的文件
		f, err := os.OpenFile("./" + handler.Filename, os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil{
			log.Println("open file err: ", err)
			fmt.Fprintf(w, "%v", err)
			fmt.Fprintf(w, `<a href="http://212.64.16.41:1111">返回</a>`)
			return
		}
		defer f.Close()        //拷贝文件
		io.Copy(f, file)

		err = exec.Command("sh","do.sh", names[0],handler.Filename).Run()
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "%v", err)
			fmt.Fprintf(w, `<a href="http://212.64.16.41:1111">返回</a>`)
			return
		}
		fmt.Fprintf(w, "success")
		fmt.Fprintf(w, `<a href="http://212.64.16.41:1111">返回</a>`)
	}

}

func test(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type","text/html; charset=utf-8")
	fmt.Fprintf(w, "success \n")
	fmt.Fprintf(w, `<a href="http://212.64.16.41:1111">返回</a>`)
}