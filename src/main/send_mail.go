package main

import (
    "flag"
    "mail"
    "fmt"
    "strings"
    // "io/ioutil"
    "path"
    "bufio"
    "os"
)

var filename = flag.String("f", "", "file name")

const (
    user      = ""
    // QQ邮箱授权码：http://service.mail.qq.com/cgi-bin/help?subtype=1&&id=28&&no=1001256
    passwd    = ""
    smtphost  = "smtp.qq.com:587"
    to        = ""
)

func main() {
    flag.Parse()
    // fByte, err := ioutil.ReadFile(*filename)
    fd, err := os.Open(*filename)
    defer fd.Close()
    if err != nil {
        fmt.Printf("Open file[%s] err[%s]\n", *filename, err)
        return
    }

    basename := path.Base(*filename)
    subject := "今日房价：" + basename

    body := `
    <html>
    <body>
    <table border="1">
    <tr>
    <th>房屋id</th>
    <th>小区id</th>
    <th>单价</th>
    <th>单价</th>
    <th>总价</th>
    <th>地址</th>
    <th>关注情况</th>
    <th>地铁</th>
    <th>税费</th>
    </tr>
    `
    endStr := `
    </table>
    </body>
    </html>
    `

    scanner := bufio.NewScanner(fd)
    for scanner.Scan() {
        body += `<tr><td>`
        line := scanner.Text()
        fLine := strings.Replace(line, "\t", "</td><td>", -1)
        body += fLine
        body += `</td></tr>`
    }

    body += endStr
    fmt.Printf("format success, now sendmai\n")

    err = mail.SendMail(user, passwd, smtphost, to, subject, body, "html")
    if err != nil {
        fmt.Printf("send mail err[%s]\n", err)
        return
    }

    fmt.Printf("send mail suc\n")
}
