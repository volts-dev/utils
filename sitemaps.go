package utils

import (
	"encoding/xml"
	//"log"
	//"os"
	//"strconv"
	//"time"
)

type (
	//
	TSitemaps struct {
		Sitemapindex []TSitemapindex `xml:"sitemapindex"`
	}

	// Sitemapindex 标签
	TSitemapindex struct {
		XMLName  xml.Name   `xml:"sitemapindex"` // 该结构转换成XML的名称
		Xmlns    string     `xml:"xmlns,attr"`   // 该标签的xmlns属性
		Sitemaps []TSitemap `xml:"sitemap"`
	}

	// Sitemap 标签
	TSitemap struct {
		XMLName xml.Name `xml:"sitemap"`
		Loc     string   `xml:"loc"`
	}

	// Urlset 标签
	TUrlset struct {
		XMLName xml.Name `xml:"urlset"`
		Xmlns   string   `xml:"xmlns,attr"`
		Url     []TUrl   `xml:"url"`
	}

	// Url 标签
	TUrl struct {
		XMLName    xml.Name `xml:"url"`
		Loc        string   `xml:"loc"`
		Lastmod    string   `xml:"lastmod"`
		Changefreq string   `xml:"changefreq"`
		Priority   string   `xml:"priority"`
	}
)

/*
func NewSitemaps() {
	sitemap_index := &TSitemapindex{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	var total int = 0
	var index int = 1

	for total <= 1000 {
		Url_set := &TUrlset{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
		var tick int = 0
		for ; tick <= 50; tick++ {
			log.Println("for2")
			new_url := TUrl{
				Loc:        "aa",
				Lastmod:    time.Now().Format("2006-01-02"),
				Changefreq: "weekly",
				Priority:   "0.4"}
			Url_set.Url = append(Url_set.Url, new_url)
			log.Println("完成URL")
		}
		sitemap_file, err := xml.Marshal(Url_set)
		filepath := "./sitemaps/sitemap" + strconv.Itoa(index) + ".xml"

		f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			//return err
			log.Println(err)
		}
		_, err = f.Write([]byte(xml.Header + string(sitemap_file)))
		if err != nil {
			//return err
			log.Println(err)
		}

		sm := TSitemap{Loc: filepath} //创建一个Sitemap标签
		sitemap_index.Sitemap = append(sitemap_index.Sitemap, sm)
		total = total + tick //
		log.Println(total)
		index++
	}

	log.Println("完成sitemaps.xml")
	sitemaps_file, err := xml.Marshal(sitemap_index)
	f, err := os.OpenFile("./sitemaps/sitemaps.xml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		//return err
	}
	_, err = f.Write([]byte(xml.Header + string(sitemaps_file)))
}
*/
