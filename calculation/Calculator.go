package calculation

import (
	"../stypes"
	"../packages/adapters"
	"database/sql"
	"fmt"
	"time"
	"log"
	"strings"
	"runtime"
	"sync"
	"math/rand"
	"regexp"
	"unicode/utf8"
)

var maxNumCpu = 8
var rxTEMP *regexp.Regexp

type Calculator struct {
	input    string
	preview1 string
	preview2 string
	output1  string
	output2  string
	numCpu   int
}
type Part struct {
	subs []string
}

func (c *Calculator) SetInput(s string) {
	c.input = s
}
func (c *Calculator) GetSpin1() *string {
	return &c.output1
}
func (c *Calculator) GetSpin2() *string {
	return &c.output2
}
func (c *Calculator) GetPreview1() *string {
	return &c.preview1
}
func (c *Calculator) GetPrevies2() *string {
	return &c.preview2
}
func (c *Calculator) Calculate() {
	start := time.Now()
	rx := regexp.MustCompile("Xq7TpK4pX")
	fmt.Printf("%q\n", c.input)
	if !utf8.ValidString(c.input) {
		v := make([]rune, 0, len(c.input))
		for i, r := range c.input {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(c.input[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		c.input = string(v)
	}
	fmt.Printf("%q\n", c.input)
	quotesadapter := regexp.MustCompile(quotesAdaptPattern)
	c.input = quotesadapter.ReplaceAllString(c.input, "'")
	arr := rx.Split(c.input, -1)
	max := runtime.GOMAXPROCS(maxNumCpu)
	l := len(arr)
	if l == 0 {
		ProgressWithOld(c)
	} else if l > 0 && l < max {
		c.numCpu = l
		ProgressWithNew(arr, l, c)
	} else {
		c.numCpu = max
		ProgressWithNew(arr, l, c)
	}
	elapsed := time.Since(start)
	log.Println(c.output1)
	log.Println(c.output2)
	fmt.Println(elapsed)
}
func ProgressWithNew(jobs []string, l int, c *Calculator) {
	//var parts = make([]*string, c.numCpu)
	//var parts2 = make([]*string, c.numCpu)
	X := l / c.numCpu
	left := l % c.numCpu
	//var partJobs [c.numCpu][X + 1]string
	var partJobs = make([]Part, c.numCpu)
	for i := 0; i < c.numCpu; i++ {
		if i < left {
			partJobs[i].subs = make([]string, X + 1)
		} else {
			partJobs[i].subs = make([]string, X)
		}

	}
	/*for i := 0; i < target; i++ {
		k1 := i % c.numCpu
		k2 := i / c.numCpu
		log.Printf("%d , %d   i= %d %s\n\n\n", k1, k2, i,jobs[i])
		partJobs[k1].subs[k2] = jobs[i]
	}*/
	var count int = 0
	for i := 0; i < c.numCpu; i++ {
		for j := 0; j < X; j++ {
			partJobs[i].subs[j] = jobs[count]
			count++
		}
		if i < left {
			partJobs[i].subs[X] = jobs[count]
			count++
		}

	}

	var parts = make([]*string, c.numCpu)
	var parts2 = make([]*string, c.numCpu)
	var prevs = make([]*string, c.numCpu)
	var prevs2 = make([]*string, c.numCpu)
	var str string = ""
	var str2 string = ""
	var prev string = "<div> "
	var prev2 string = "<div> "
	if c.numCpu == 1 {
		db := adapters.ConnectDB(1)
		defer db.Close()
		p1, p2, temp1, temp2 := DoWorkNew(partJobs[0].subs, db)
		str = string(*temp1)
		str2 = string(*temp2)
		prev = string(*p1)
		prev2 = string(*p2)

	} else {
		var wg sync.WaitGroup
		wg.Add(c.numCpu)
		for i := 0; i < c.numCpu; i++ {
			go WorkNew(&wg, i, &partJobs, prevs, prevs2, parts, parts2)
		}
		wg.Wait()
		for i := 0; i < len(parts); i++ {
			str += string(*parts[i])
			str2 += string(*parts2[i])
			prev += string(*prevs[i])
			prev2 += string(*prevs2[i])
		}
	}
	c.preview1 = prev + " </div>"
	c.preview2 = prev2 + " </div>"
	c.output1 = str
	c.output2 = str2
}

func ProgressWithOld(c *Calculator) {
	mod := (len(c.input) / 600)
	if mod == 0 {
		c.numCpu = 1
	} else if mod > 0 && mod <= 1 {
		c.numCpu = 2
	} else if mod > 1 && mod <= 2 {
		c.numCpu = 3
	} else if mod > 2 && mod <= 4 {
		c.numCpu = 4
	} else if mod > 4  && mod <= 8 {
		c.numCpu = 5
	} else if mod > 8 && mod <= 16 {
		c.numCpu = 6
	} else if mod > 16 && mod <= 32 {
		c.numCpu = 7
	} else if mod > 32 {
		c.numCpu = 8
	}
	runtime.GOMAXPROCS(maxNumCpu)
	tasks := DistributeTasks(&c.input, &c.numCpu)
	var parts = make([]*string, c.numCpu)
	var parts2 = make([]*string, c.numCpu)
	var str string = ""
	var str2 string = ""
	if c.numCpu == 1 {
		temp1, temp2, _ := DoWork(tasks[0], 1)
		str = string(*temp1)
		str2 = string(*temp2)

	} else {
		var wg sync.WaitGroup
		wg.Add(c.numCpu)

		for i := 0; i < c.numCpu; i++ {
			go Work(&wg, i, tasks, parts, parts2)
		}
		wg.Wait()
		for i := 0; i < len(parts); i++ {
			str += string(*parts[i])
			str2 += string(*parts2[i])
		}
	}
	c.output1 = str
	c.output2 = str2
}

func Work(wg *sync.WaitGroup, cpu int, tasks[]string, parts[]*string, parts2[]*string) {
	defer wg.Done()
	r, r2, t := DoWork(tasks[cpu], cpu + 1)
	parts[t - 1] = r
	parts2[t - 1] = r2
}
func WorkNew(wg *sync.WaitGroup, cpu int, tasks *[]Part, prevs[]*string, prevs2[]*string, parts[]*string, parts2[]*string) {
	defer wg.Done()
	log.Println(cpu)
	db := adapters.ConnectDB(cpu + 1)
	defer db.Close()
	p1, p2, r, r2 := DoWorkNew((*tasks)[cpu].subs, db)
	prevs[cpu] = p1
	prevs2[cpu] = p2
	parts[cpu] = r
	parts2[cpu] = r2
}
func SearchLookAtRegion(s *string, avg int, index int) int {
	i := avg * index
	var result, alterna, foo int
	var isFound, isFoundAlterna, canBeFound = false, false, false

	for count := 0; count < int(float32(avg) * 0.7); count++ {
		switch (*s)[i + count]{
		case '.':
			if !canBeFound {
				result = i + count
				canBeFound = true
			}

		case ',', '\n', ':', ';':
			result = i + count
			isFound = true
			break
		case ' ':
			if canBeFound {
				isFound = true
				break
			}
			if !isFoundAlterna {
				alterna = i + count
				isFoundAlterna = true
			}

		default:
			if canBeFound {
				foo++
				if foo > 2 {
					canBeFound = false
					foo = 0
				}
			}
		}
	}
	if isFound {
		return result
	} else {
		if isFoundAlterna {
			return alterna
		} else {
			//TODO:: işlem bu sefer tam tersi olarak yapılması gerekiyor.
			return avg * index
		}
	}
}
func DistributeTasks(in *string, cpu *int) []string {
	var tasks = make([]string, *cpu)
	long := len(*in)
	avg := long / *cpu
	var old int
	for i := 1; i < *cpu; i++ {
		r := SearchLookAtRegion(in, avg, i)
		if i == 1 {
			//first task
			tasks[i - 1] = (*in)[:r + 1]
		} else {
			//middle tasks
			tasks[i - 1] = (*in)[old + 1:r + 1]
		}
		old = r
	}
	if *cpu == 1 {
		//one task
		tasks[*cpu - 1] = (*in)
	} else {
		//last task
		tasks[*cpu - 1] = (*in)[old + 1:]

	}
	return tasks
}
func DoWork(job string, cpu int) (*string, *string, int) {
	db := adapters.ConnectDB(cpu)
	defer db.Close()
	list := strings.Split(job, " ")
	var result string = ""
	var result2 string = ""
	var buff int = 1
	for i := 0; i < len(list); {
		stack := stypes.NewStack(11)
		row := list[i]
		buff = 1
		if strings.TrimSpace(row) != "" {
			stack.Push(row)
			buff++
			temp := row
			for count := 1; i + count < len(list) && count <= 10; count++ {
				temp += " " + list[i + count]
				if !isPunctuation(getLastChar(&temp)) {
					stack.Push(temp)
					buff++
				} else {
					break
				}
			}
			var r1, isFound bool = false, false
			var rs, rs2 *string
			for {
				temp := stack.Pop()
				if temp != nil {
					buff--
					temps := temp.(string)
					r1, rs, rs2 = Search(&temps, db)
					if r1 {
						isFound = true
						result += " " + *rs
						result2 += " " + *rs2
						break
					}
				} else {
					break
				}
			}
			if (!isFound) {
				result += " " + *rs
				result2 += " " + *rs2
			}
		}
		i += buff
	}
	return &result, &result2, cpu
}
func DoWorkNew(job[] string, db *sql.DB) (*string, *string, *string, *string) {
	rx := regexp.MustCompile(clearTextPattern)
	rxabv := regexp.MustCompile(abbreviationPattern)
	var result string = ""
	var result2 string = ""
	var pre1 string = ""
	var pre2 string = ""
	for _, row := range job {
		p1, p2, s1, s2 := DoSpin(row, db, rx, rxabv)
		result += *s1
		result2 += *s2
		pre1 += *p1
		pre2 += *p2
	}
	return &pre1, &pre2, &result, &result2
}
func DoSpin(job string, db *sql.DB, rxClearText *regexp.Regexp, rxCheckAbbrav *regexp.Regexp) (*string, *string, *string, *string) {
	log.Printf("İş------ %s", job)
	var pure string = job
	var index = len(job)
	var isAbbrav bool = rxCheckAbbrav.MatchString(job)
	if isAbbrav == false {
		pure = rxClearText.FindString(job)
		//TODO: index out of range hatası!
		if pure != "" {
			index = rxClearText.FindStringIndex(job)[0]
		}

	}
	var result string = job[:index]
	var result2 string = job[:index]
	var prev1 string = result
	var prev2 string = result2
	index += len(pure)
	list := strings.Split(pure, " ")
	var buff int = 1
	for i := 0; i < len(list); {
		stack := stypes.NewStack(11)
		row := list[i]
		buff = 1
		if strings.TrimSpace(row) != "" {
			stack.Push(row)
			buff++
			temp := row
			for count := 1; i + count < len(list) && count <= 8; count++ {
				temp += " " + list[i + count]
				if !isPunctuation(getLastChar(&temp)) {
					stack.Push(temp)
					buff++
				} else {
					stack.Push(temp)
					buff++
					break
				}
			}
			var r1, isFound bool = false, false
			var rs, rs2 *string
			for {
				temp := stack.Pop()
				if temp != nil {
					buff--
					temps := temp.(string)
					r1, rs, rs2 = Search(&temps, db)
					if r1 {
						isFound = true
						result += " " + *rs
						result2 += " " + *rs2
						prev1 += colored + *rs + "</font>"
						prev2 += colored + *rs2 + "</font>"
						break
					}
				} else {
					break
				}
			}
			if (!isFound) {
				result += " " + *rs
				result2 += " " + *rs2
				prev1 += " " + *rs
				prev2 += " " + *rs2
			}
		}
		i += buff
	}
	if (isAbbrav == false) {
		result += job[index:]
		result2 += job[index:]
		prev1 += job[index:]
		prev2 += job[index:]
	}
	return &prev1, &prev2, &result, &result2
}
func Search(s *string, db *sql.DB) (bool, *string, *string) {
	//if isPunctuation(s[0])
	var Spin sql.NullString
	log.Printf("Önceki %s\n", *s)
	//var Spin string
	q := ClearText(s)

	//log.Println(*q)
	//query := "SELECT Spin FROM generalspins where ( Spin like '%|" + string(*s) + "|%' or Spin like '{" + string(*s) + "|%' or Spin like '%|" + string(*s) + "}' ) and GeneralSpinsID < 30000;"
	if (strings.ContainsAny(*q, ",\":;!")) {
		log.Printf("İptal  %s\n", *q)
		return false, s, s
	}
	log.Printf("Temiz  %s\n", *q)
	err := db.QueryRow(
		"SELECT group_concat(Spin) as Spin FROM spins " +
			"WHERE GroupID = (SELECT GroupID FROM spins " +
			"where Spin = ? LIMIT 1) " +
			"AND Spin != ? LIMIT 1;", string(*q), string(*q)).Scan(&Spin);

	switch {
	case err == sql.ErrNoRows:
		//spin = line;
		//log.Println("0", *s)
		return false, s, s
	case !Spin.Valid:
		return false, s, s
	case err != nil:
		log.Println(err.Error())
		return false, s, s
	default:
		r1, r2 := GetStringCreateArray(&Spin.String, s)
		rr1 := strings.Replace(*s, *q, *r1, 1)
		rr2 := strings.Replace(*s, *q, *r2, 1)
		return true, &rr1, &rr2
	}
}
//func isPunctuation(c char)

const delim = "\"':?!;,*(){}[]/%"
const sentencePattern = "[^.!?\\s][^.!?]*(?:[.!?](?!['\"]?\\s|$)[^.!?]*)*[.!?]?['\"]?(?=\\s|$)"
const clearTextPattern = "[\\wöÖçÇğĞüÜıİşŞ]+[\\d\\D]*[\\wöÖçÇğĞüÜıİşŞ]+"
const abbreviationPattern = "(?:[a-zA-ZöÖçÇğĞüÜıİşŞ]\\.){2,}[.]?[ ]*$"
const colored = "<font class=\"colored\">"
const quotesAdaptPattern = "['`´ʹʻʼʽʾʿˈˊˋ]"

func ClearText(s *string) *string {
	a := strings.TrimSpace(*s)
	if (len(a) > 2) {
		c := getFirstChar(&a)
		if isPunctuation(c) {
			a = strings.TrimLeft(a, c)
		}
		c = getLastChar(&a)
		if isPunctuation(c) {
			a = strings.TrimRight(a, c)
		}
	}
	return &a
}
func TrimSuffix(s string) string {
	if strings.HasSuffix(s, delim) {
		s = s[:len(s) - len(delim)]
	}
	return s
}
func getFirstChar(s *string) string {
	return string([]rune(*s)[0])
}
func getLastChar(s *string) string {
	t := string(*s)
	if (len(t) > 1) {
		return string([]rune(t[len(t) - 1:])[0])
	} else {
		return string(*s)
	}

}
func isPunctuation(c string) bool {
	if strings.Contains(delim, c) {
		return true
	}
	return false
}

func GetString(s *string, orginal *string) (*string, *string) {
	temp := []byte(*s)
	arr := strings.Split(string(temp[1:len(temp) - 1]), "|")
	var randomed int
	for {
		randomed = rand.Intn(len(arr))
		if arr[randomed] != *orginal {
			break
		}
	}
	var randomed2 int
	for {
		randomed2 = rand.Intn(len(arr))
		if arr[randomed2] != *orginal {
			break
		}
	}
	return &arr[randomed], &arr[randomed2]

}
func GetStringCreateArray(s *string, orginal *string) (*string, *string) {
	arr := strings.Split(*s, ",")
	var randomed int
	for {
		randomed = rand.Intn(len(arr))
		if arr[randomed] != *orginal {
			break
		}
	}
	var randomed2 int
	for {
		randomed2 = rand.Intn(len(arr))
		if arr[randomed2] != *orginal {
			break
		}
	}
	return &arr[randomed], &arr[randomed2]

}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}