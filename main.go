package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/manifoldco/promptui"
)

func main() {

	//declare file names
	fName := "date.csv"
	fNameLink := "link.csv"

	//create files
	file, err := os.Create(fName)
	if err != nil {
		fmt.Println("faili ei loodud")
	}

	fileLink, err := os.Create(fNameLink)
	if err != nil {
		fmt.Println("file ei loodud")
	}

	//closing files
	defer fileLink.Close()
	defer file.Close()

	//idk, maybe creates writers for files and....
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writerlink := csv.NewWriter(fileLink)
	defer writerlink.Flush()

	//mis lehti tohib scrapida
	c := colly.NewCollector(
		colly.AllowedDomains("www.postimees.ee", "postimees.ee"),
	)

	//get headlines
	c.OnHTML(".list-article__text", func(e *colly.HTMLElement) {

		writer.Write([]string{
			e.Text,
		})
	})

	//get links
	c.OnHTML(".list-article__url", func(e *colly.HTMLElement) {

		writerlink.Write([]string{
			e.Attr("href"),
		})
	})

	//go get shit done
	c.Visit("https://www.postimees.ee")

	//open date.csv
	data, err := os.Open("date.csv")
	if err != nil {
		fmt.Println("ei avatud date.csv")
	}
	//read date.csv
	r := csv.NewReader(data)
	lines, err := r.ReadAll()
	if err != nil {
		fmt.Printf("vahepeal ei toota sest : %s", err)
		os.Exit(1)
	}
	//create slice for use in terminal for date.csv
	pealkiriSlice := []string{}
	for _, dd := range lines {
		for _, d := range dd {
			pealkiriSlice = append(pealkiriSlice, d)
		}
	}

	//start opening link.csv
	//open link.csv
	dataLink, err := os.Open("link.csv")
	if err != nil {
		fmt.Println("error link")
	}
	//read link.csv
	rLink := csv.NewReader(dataLink)
	linesLink, err := rLink.ReadAll()
	if err != nil {
		fmt.Printf("sama nuss mis data.csv %s", err)
	}
	//create a slice for terminal link.csv
	linkSlice := []string{}
	for _, dd := range linesLink {
		for _, d := range dd {
			linkSlice = append(linkSlice, d)
		}
	}

	//terminal stuff i dont understand
	prompt := promptui.Select{
		Label: "Pealkirjad",
		Items: pealkiriSlice,
		Size:  35,
	}
	//starts terminal app
	index, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	_ = result
	//PRINT WHAT WAS CHOSEN
	//fmt.Printf("Pealkiri:   %s \n", result)
	fmt.Println()

	//use index of the selected item and get url
	newsLinkFromIndex := linkSlice[index]

	//get content of selected item
	// lehesisu := ""
	// c.OnHTML(".article-body__item.article-body__item--htmlElement.article-body__item--lead", func(e *colly.HTMLElement) {

	// 	// writerlink.Write([]string{
	// 	// 	e.Text,
	// 	// })
	// 	//fmt.Println(e.Text)
	// 	e.ChildText("p")
	// 	lehesisu = e.ChildText("p")
	// 	fmt.Printf("kontroll %s", e.Text)
	// })

	// c.Visit("https://www.postimees.ee/7285762/venemaalt-eestisse-peagi-enam-nii-lihtsalt-ei-paase")
	// fmt.Printf("lehe sisu on %s\n", lehesisu)

	displayContent(newsLinkFromIndex)

}

func displayContent(leheURL string) {

	var funktsiooniSisuText string
	funktsiooniAlamSisu := ""
	c := colly.NewCollector(
	//colly.AllowedDomains("www.postimees.ee", "postimees.ee"),
	)

	c.OnHTML(".article-body__item.article-body__item--htmlElement.article-body__item--lead", func(e *colly.HTMLElement) {

		funktsiooniSisuText = string(e.ChildText("p"))
	})
	c.OnHTML(".article-body.article-body--left", func(e *colly.HTMLElement) {

		funktsiooniAlamSisu = e.Text

	})

	c.Visit(leheURL)

	fmt.Println(word_wrap(funktsiooniSisuText, 5))
	fmt.Println()
	fmt.Println(word_wrap(funktsiooniAlamSisu, 9))

}

//stole it all

func word_wrap(s string, limit int) string {

	if strings.TrimSpace(s) == "" {
		return s
	}

	// convert string to slice
	strSlice := strings.Fields(s)

	var result string = ""

	for len(strSlice) >= 1 {
		// convert slice/array back to string
		// but insert \r\n at specified limit

		result = result + strings.Join(strSlice[:limit], " ") + "\r\n"

		// discard the elements that were copied over to result
		strSlice = strSlice[limit:]

		// change the limit
		// to cater for the last few words in
		//
		if len(strSlice) < limit {
			limit = len(strSlice)
		}

	}

	return result

}
