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
    urlPattern  = "https://%s.lianjia.com/chengjiao/pg%dc%s/"
    fPattern    = "%s.xiaoqu_%s.%s"
    timePattern = "2006-01-02"
)

var maxNum      = flag.Int("n", 10, "grab max pages")
var filePath    = flag.String("p", "/tmp/", "data store path")
var province    = flag.String("province", "bj", "province you wanna")
var xiaoquId    = flag.String("xiaoqu", "1111027378424", "xiaoqu id you wanna")

func main() {
    var buf bytes.Buffer
    flag.Parse()

    for i := 1; i <= *maxNum; i++ {
        processUrl := fmt.Sprintf(urlPattern, *province, i, *xiaoquId)
        // fmt.Printf("processing url[%s] xiaoqu[%s]\n", processUrl, *xiaoquId)
        doc, err := goquery.NewDocument(processUrl)
        if err != nil {
            log.Printf("url[%s] NewDocument err[%s]\n", processUrl, err)
            continue
        }

        doc.Find("ul.listContent").Find("li").Each(func(i int, s *goquery.Selection) {
            housetitle := strings.TrimSpace(s.Find(".title").Text())
            houseInfo  := s.Find(".houseInfo").Text()
            dealDate   := s.Find(".dealDate").Text()
            totalPrice := s.Find("totalPrice").Text()
            posInfo    := s.Find(".positionInfo").Text()
            unitPrice  := s.Find(".unitPrice").Text()
            dealInfo   := s.Find(".dealHouseInfo").Text()
            dealcyInfo := s.Find(".dealCycleeInfo").Text()

            // houseInfoStr := trimSpace(houseInfo)
            // addrStr      := trimSpace(addr)
            // iPriceStr    := trimSpace(iPrice)
            // iCountStr    := trimSpace(iCount)

            
log.Printf("[%d]th xiaoquId[%s] xiaoqu[%s] houseInfo[%s] dealDate[%s] totalPrice[%s] posInfo[%s] unitPrice[%s] dealInfo[%s] dealcyInfo[%s]\n",
            i, *xiaoquId, housetitle, houseInfo, dealDate, totalPrice, posInfo, unitPrice, dealInfo, dealcyInfo)
            buf.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", 
                housetitle, houseInfo, dealDate, totalPrice, posInfo, unitPrice, dealInfo, dealcyInfo))
        })
    }

    timeStr := time.Now().Format(timePattern)
    fName := fmt.Sprintf(fPattern, *province, *xiaoquId, timeStr)
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
