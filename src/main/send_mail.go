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
    user      = "jiujiangyuan@qq.com"
    passwd    = ""
    smtphost  = "smtp.qq.com:587"
    to        = "jiujiangyuan@qq.com;873029024@qq.com"
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

    err = mail.SendMail(user, passwd, smtphost, to, subject, body, "html")
    if err != nil {
        fmt.Printf("send mail err[%s]\n", err)
        return
    }

    fmt.Printf("send mail suc\n")
}
