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
    sj  "github.com/guyannanfei25/go-simplejson"
)

const (
    urlPattern  = "http://%s.lianjia.com/ershoufang/%s/pg%dp3p4"
    fPattern    = "%s.ershoufang.%s.%s"
    timePattern = "2006-01-02"
)

var maxNum      = flag.Int("n", 100, "grab max pages")
var filePath    = flag.String("p", "/tmp/", "data store path")
var province    = flag.String("province", "bj", "province you wanna")
var area        = flag.String("a", "haidian", "area you wanna, 比如：海淀区(haidian)、朝阳区(chaoyang)")

func main() {
    var buf bytes.Buffer
    flag.Parse()
    first := true

    for i := 1; i <= *maxNum; i++ {
        processUrl := fmt.Sprintf(urlPattern, *province, *area, i)
        doc, err := goquery.NewDocument(processUrl)
        if err != nil {
            log.Printf("url[%s] NewDocument err[%s]\n", processUrl, err)
            continue
        }

        if first {
            // find total num in case be forbidden
            totalStr, ok := doc.Find("div.house-lst-page-box").Attr("page-data")
            if !ok {
                log.Printf("not found page-box\n")
                continue
            }

            totalJson, err := sj.NewJson([]byte(totalStr))
            if err != nil {
                log.Printf("new json er[%s]\n", err)
                continue
            }

            totalPage := totalJson.Get("totalPage").MustInt()
            log.Printf("total page[%d]\n", totalPage)
            *maxNum = totalPage

            first = false
        }

        doc.Find("li.clear").Each(func(i int, s *goquery.Selection) {
            id, ok := s.Find(".unitPrice").Attr("data-hid")
            if !ok {
                fmt.Printf("[%d]th cannot find id\n", i)
            }

            xiaoquId, ok := s.Find(".unitPrice").Attr("data-rid")
            if !ok {
                fmt.Printf("[%d]th cannot find xiaoqu id\n", i)
            }

            price, ok := s.Find(".unitPrice").Attr("data-price")
            if !ok {
                fmt.Printf("[%d]th cannot find price\n", i)
            }

            priceStr := s.Find(".unitPrice").Text()
            totalPrice := s.Find(".totalPrice").Text()
            addr     := s.Find(".address").Find(".houseInfo").Text()
            follow   := s.Find(".followInfo").Text()
            subway   := s.Find(".subway").Text()
            taxfree  := s.Find(".taxfree").Text()
            
            log.Printf("[%d]th id[%s] xiaoquId[%s] price[%s] priceStr[%s] totalPrice[%s] addr[%s] follow[%s] subway[%s] taxfree[%s]\n",
                i, id, xiaoquId, price, priceStr, totalPrice, addr, follow, subway, taxfree)
            buf.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", 
                id, xiaoquId, price, priceStr, totalPrice, addr, follow, subway, taxfree))
        })
    }

    timeStr := time.Now().Format(timePattern)
    fName := fmt.Sprintf(fPattern, *province, *area, timeStr)
    fullPath := path.Join(*filePath, fName)
    if err := ioutil.WriteFile(fullPath, buf.Bytes(), 0766); err != nil {
        log.Printf("Write content[%s] to file[%s] err[%s]\n",
            buf.String(), fullPath, err)
    }

}
