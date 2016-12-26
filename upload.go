package main
import ( 
    "fmt" 
    "os" 
    "strings" 
    "io"
    "net/http"
)

const (
    UPLOAD_PATH string = "/live/www/html/livepic/"
//    ACCESS_URL string = "http://live.66boss.com/livepic/"
)
func getUUID() string { 
    f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0) 
    b := make([]byte, 16) 
    f.Read(b) 
    f.Close() 
    uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]) 
    return uuid
}
func getFileName(name string) string {

    var temp string = "error"
    i := strings.Index(name, ".")
    if i > 0 {
        uuid := getUUID()
        temp = uuid + name[i:]
    }
    return temp
    

}
 
 
func uploadHandle(w http.ResponseWriter, r *http.Request) {

    if r.Method == "GET" {
        io.WriteString(w, "<html><head><title>我的第一个页面</title></head><body><form action='' method=\"post\" enctype=\"multipart/form-data\"><label>上传图片</label><input type=\"file\" name='file'  /><br/><label><input type=\"submit\" value=\"上传图片\"/></label></form></body></html>")
    } else {
        file, head, err := r.FormFile("file")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()
        
        temp := getFileName(head.Filename)
        uuidFile := UPLOAD_PATH + temp
        fW, err := os.Create(uuidFile)
        if err != nil {
            fmt.Println("create file error")
            return
        }
        defer fW.Close()
        _, err = io.Copy(fW, file)
        if err != nil {
            fmt.Println("copy file error")
            return
        }
        //io.WriteString(w, (ACCESS_URL + temp))
        io.WriteString(w, temp)
        //http.Redirect(w, r, "/hello", http.StatusFound)
    }
}
 
