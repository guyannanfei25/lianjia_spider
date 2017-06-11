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
)

const (
    urlPattern  = "http://%s.lianjia.com/ershoufang/pg%d"
    fPattern    = "%s.ershoufang.%s"
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
    fName := fmt.Sprintf(fPattern, *province, timeStr)
    fullPath := path.Join(*filePath, fName)
    if err := ioutil.WriteFile(fullPath, buf.Bytes(), 0766); err != nil {
        log.Printf("Write content[%s] to file[%s] err[%s]\n",
            buf.String(), fullPath, err)
    }

}
