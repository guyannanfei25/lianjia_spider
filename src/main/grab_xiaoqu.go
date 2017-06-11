package main

import (
    "io/ioutil"
    "bytes"
    "flag"
    "time"
    "github.com/PuerkitoBio/goquery"
    "log"
    "fmt"
    "path"
    "strings"
    "unicode"
)

const (
    urlPattern  = "http://%s.lianjia.com/xiaoqu/pg%d"
    fPattern    = "%s.xiaoqu.%s"
    timePattern = "2006-01-02"
)

var maxNum      = flag.Int("n", 100, "grab max pages")
var filePath    = flag.String("p", "/tmp/", "data store path")
var province    = flag.String("province", "bj", "province you wanna")

func main() {
    var buf bytes.Buffer
    flag.Parse()

    for i := 1; i <= *maxNum; i++ {
        processUrl := fmt.Sprintf(urlPattern, *province, i)
        doc, err := goquery.NewDocument(processUrl)
        if err != nil {
            log.Printf("url[%s] NewDocument err[%s]\n", processUrl, err)
            continue
        }

        doc.Find("li.xiaoquListItem").Each(func(i int, s *goquery.Selection) {
            xiaoquid, ok := s.Find("a.img").Attr("href")
            fmt.Printf("id[%d] herf[%s]\n", i, xiaoquid)

            if !ok {
                fmt.Printf("id[%d] cannot find xiaoquid\n", i)
            }

            xiaoquId := GetSubStr(xiaoquid, "xiaoqu/", "/")

            housetitle := strings.TrimSpace(s.Find(".title").Text())
            houseInfo  := s.Find(".houseInfo").Text()
            addr       := s.Find(".positionInfo").Text()
            tag        := strings.TrimSpace(s.Find(".tagList").Text())
            iPrice     := s.Find(".xiaoquListItemPrice").Text()
            iCount     := s.Find(".xiaoquListItemSellCount").Text()

            houseInfoStr := trimSpace(houseInfo)
            addrStr      := trimSpace(addr)
            iPriceStr    := trimSpace(iPrice)
            iCountStr    := trimSpace(iCount)

            
            log.Printf("[%d]th xiaoquId[%s] xiaoqu[%s] info[%s] position[%s] tag[%s] price[%s] count[%s]\n",
                i, xiaoquId, housetitle, houseInfoStr, addrStr, tag, iPriceStr, iCountStr)
            buf.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n", 
                xiaoquId, housetitle, houseInfoStr, addrStr, tag, iPriceStr, iCountStr))
        })
    }

    timeStr := time.Now().Format(timePattern)
    fName := fmt.Sprintf(fPattern, *province, timeStr)
    fullPath := path.Join(*filePath, fName)
    if err := ioutil.WriteFile(fullPath, buf.Bytes(), 0766); err != nil {
        log.Printf("Write content[%s] to file[%s] err[%s]\n",
            buf.String(), fullPath, err)
    }

}

func trimSpace(str string) string {
    tmp := strings.FieldsFunc(str, unicode.IsSpace)
    return strings.Join(tmp, " ")
}

func GetSubStr(str, pre, suf string) string {
    begin := strings.Index(str, pre)
    if begin == -1 {
        return ""
    }

    end := strings.Index(str[begin + len(pre):], suf)
    if end == -1 {
        return ""
    }

    return str[begin + len(pre) : begin + len(pre) + end]
}
