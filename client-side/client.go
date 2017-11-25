package main

import (
	"crypto/rand"
	"regexp"
	"log"
)

const patt = "[\\wöÖçÇğĞüÜıİşŞ]+[\\d\\D]*[\\wöÖçÇğĞüÜıİşŞ]+"

func GetRandomAPIKey() string {
	var dictionary string
	dictionary = "0123ABCDEFGHKLMTUVWXYZ        aaaeeeiiiiooouuuuabcdefghnopqrstuvwxyz"
	var bytes = make([]byte, 35)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v % byte(len(dictionary))]
	}
	return string(bytes)
}
const clearTextPattern = "[\\wöÖçÇğĞüÜıİşŞ]+[\\d\\D]*[\\wöÖçÇğĞüÜıİşŞ]+"
const quotesAdaptPattern = "['`´ʹʻʼʽʾʿˈˊˋ]"
func main() {
	//// Take substring from index 4 to length of string.
	//substring := value[4:len(value)]
	//rxTEMP := regexp.MustCompile("(?:[a-zA-Z]\\.){2,} ")
	rxTEMP := regexp.MustCompile("(?:[a-zA-ZöÖçÇğĞüÜıİşŞ]\\.){2,}[.]?[ ]*$")

	var input = "Ç.L.Ş."
	arrStr := rxTEMP.MatchString(input)
	quotesadapter := regexp.MustCompile(quotesAdaptPattern)
	log.Println(quotesadapter.ReplaceAllString("hayrullah'ın ´kalemʿi `s´dʹcʻzʼsʽzʾetʿerˈmˊzˋ","'"))
	arrIndex := rxTEMP.FindAllStringSubmatchIndex(input,-1)
//	pure := rxTEMP.FindString(input)
//	index := rxTEMP.FindStringIndex(input)[0]
	log.Println(arrStr)
	log.Println(len(arrIndex))
	//substring := value[:2]
	/*rxTEMP := regexp.MustCompile("[\\wöÖçÇğĞüÜıİşŞ]+[,*_-].[\\s]*")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var input string = "Havalimanı'nda unutulan binlerce cep telefonu, laptop, bilgisayar, fotoğraf makinası gibi binlerce bavul dolu eşya yarın satışa çıkıyor"
	fmt.Println(rxTEMP.FindAllStringSubmatch(input, -1))*/
	/*var temp string = "uLevFU3oFneADB"
	var api string = GetRandomAPIKey()
	fmt.Println(api)
	for i := 0; i < 400; i++ {
		for j := 0; j < 25; j++ {
			temp += "deneme   " + api
		}
	}
	temp += temp + temp + temp + temp
	fmt.Println(len(temp))
	start := time.Now()
	*//*rex := regexp.MustCompile("deneme   " + api)
	out := rex.FindAllStringSubmatch(temp, -1)*//*
	fmt.Println(strings.Count(temp,"{tarzı}"))
	time.Sleep(1000 * time.Millisecond)
	elapsed := time.Since(start)
	//println(len(out))
	fmt.Println(elapsed)*/

}

