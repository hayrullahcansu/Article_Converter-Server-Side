package main

import (
	"./calculation"
	"./packages/adapters"
	"./packages"
	"./stypes"
	"./packages/security"
	"log"
	"net/http"
	"flag"
	"io/ioutil"
	"fmt"
	"os"
	"database/sql"
	_ "database/mysql"
	"bufio"
	"strings"
	"time"
	"math/rand"
	"net/mux"
	"encoding/json"
	"runtime"
	"io"
)

var addr = flag.String("addr", "127.0.0.1:3291", "http service address")
var addrWSsocket = flag.String("websocket", "/wsocket", "http websocket handler function path")
var addrSProvider = flag.String("serviceprovider", "/sprovider", "http service provider handler function path")
var addrSpinner = flag.String("spinner", "/spinner", "http spinner handler function path")
var addrSignUp = flag.String("signup", "/signup", "http spinner handler function path")
var addrSpinFree = flag.String("spinnerfree", "/spinnerfree", "http spinner free handler");

var appLiveVersion = flag.String("app_live_version", "1.0.0", "http spinner handler function path")

var numCPU = runtime.NumCPU()
var err error
var APIRegister = packages.NewRegister()

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	htmlData, err := ioutil.ReadAll(r.Body) //<--- here!
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// print out
	fmt.Println(string(htmlData)) //<-- here !
	//packages.ServeWs(server, w, r)
}
func ServiceProviderHandlerOLD(w http.ResponseWriter, r *http.Request) {
	log.Println(formatRequest(r))
	decoder := json.NewDecoder(r.Body)
	var spmsg stypes.SPMessage
	err := decoder.Decode(&spmsg)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	var msg *stypes.SPMessage
	log.Printf("%d %s %s ", spmsg.ContentCode, spmsg.Client, spmsg.Content)
	if spmsg.ContentCode == 100 {
		i, j, u := CheckServiceRequest(&spmsg.Content)
		if i {
			switch j{
			case "ok":
				ApiKey := security.GetRandomAPIKey()
				APIRegister.AddRegister(u, &ApiKey)
				msg = PrepareServicesInfo(u, &ApiKey)
			case "alreadyconnected":
				APIRegister.RemoveRegister(u)
				ApiKey := security.GetRandomAPIKey()
				APIRegister.AddRegister(u, &ApiKey)
				msg = PrepareServicesInfo(u, &ApiKey)
			}
		} else {
			switch j {
			case "notfounduser":
				msg = PrepareNotFoundUser()
			case "version":
			case "os":
			case "isntconfirmed":
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(msg)
	if err != nil {
		panic(err)
	}

}
func ServiceProviderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(formatRequest(r))
	decoder := json.NewDecoder(r.Body)
	var spmsg stypes.SPMessage
	err := decoder.Decode(&spmsg)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	var msg *stypes.SPMessage
	log.Printf("%d %s %s ", spmsg.ContentCode, spmsg.Client, spmsg.Content)
	if spmsg.ContentCode == 100 {
		var user stypes.User
		err := json.Unmarshal([]byte(spmsg.Content), &user)
		if err != nil {
			msg = PrepareError()
			panic(err)
		}
		if APIRegister.CheckRegister(&user) {
			ApiKey := security.GetRandomAPIKey()
			APIRegister.AddRegister(&user, &ApiKey)
			msg = PrepareServicesInfo(&user, &ApiKey)
		} else {
			APIRegister.RemoveRegister(&user)
			ApiKey := security.GetRandomAPIKey()
			APIRegister.AddRegister(&user, &ApiKey)
			msg = PrepareServicesInfo(&user, &ApiKey)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(msg)
	if err != nil {
		panic(err)
	}

}
func PrepareNotFoundUser() *stypes.SPMessage {
	return &stypes.SPMessage{
		Client:    "server",
		ContentCode: 102,
		Content: "notfounduser",
		UserID: 0,
		APIKey: "",
	}
}
func PrepareError() *stypes.SPMessage {
	return &stypes.SPMessage{
		Client:    "server",
		ContentCode: 109,
		Content: "error",
		UserID: 0,
		APIKey: "",
	}
}
func PrepareServicesInfo(u *stypes.User, api *string) *stypes.SPMessage {
	return &stypes.SPMessage{
		Client:    "server",
		ContentCode: 101,
		Content: "logined",
		UserID: u.UserID,
		APIKey: *api,
	}
}
func CheckServiceRequest(s *string) (bool, string, *stypes.User) {
	var temp stypes.ClientInfo
	err := json.Unmarshal([]byte(*s), &temp)
	if err != nil {
		panic(err)
	}
	r, u := CheckUserDatabaseAndKey(&temp.Mail, &temp.Password)
	switch r {
	case 1:
		return true, "ok", u
	case 2:
		return true, "alreadyconnected", u
	default:
		return false, "notfounduser", nil
	}
}
func CheckUserDatabaseAndKey(mail *string, pass *string) (int, *stypes.User) {
	db := adapters.ConnectDB(1)
	defer db.Close()
	var user stypes.User
	log.Println(*pass)
	log.Println(*mail)
	row := db.QueryRow("SELECT UserID,Mail FROM users WHERE Mail = ? AND Password = ?", *mail, *pass)
	/*var p []byte
	err := db.QueryRow("SELECT data FROM message WHERE data->>'id'=$1", id).Scan(&p)
	if err != nil {
		// handle error
	}
	var m Message
	err := json.Unmarshal(p, &m)
	if err != nil {
		// handle error
	}*/
	err := row.Scan(&user.UserID, &user.Mail)
	switch {
	case err == sql.ErrNoRows:
		return 3, nil
	case err != nil:
		log.Println(err.Error())
		return 3, nil
	default:

	}
	if APIRegister.CheckRegister(&user) == true {
		return 2, &user
	} else {
		return 1, &user
	}

}
func SpinnerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(formatRequest(r))
	decoder := json.NewDecoder(r.Body)
	var sm stypes.SMessage
	err := decoder.Decode(&sm)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Printf("%d %s %s %s %s ", sm.ContentCode, sm.Client, sm.Content, sm.Content, sm.APIKey);
	if sm.ContentCode == 300 {
		switch APIRegister.CompareKeyUserID(&sm.APIKey, &sm.UserID) {
		case -1:
			sm.Client = "Server"
			sm.ContentCode = 404
			sm.Content = "invalidKey"
			sm.APIKey = ""
			sm.Spin1 = ""
			sm.Spin2 = ""
			sm.UserID = 0
		case -2:
			sm.Client = "Server"
			sm.ContentCode = 405
			sm.Content = "notfoundkey"
			sm.APIKey = ""
			sm.Spin1 = ""
			sm.Spin2 = ""
			sm.UserID = 0
		case 1:
			var calc calculation.Calculator
			calc.SetInput(sm.Content)
			calc.Calculate()
			sm.Spin1 = *calc.GetSpin1()
			sm.Spin2 = *calc.GetSpin2()
			sm.Preview1 = *calc.GetPreview1()
			sm.Preview2 = *calc.GetPrevies2()
			sm.Client = "Server"
			sm.ContentCode = 301
			sm.Content = "ok"
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(sm)
	if err != nil {
		panic(err)
	}
}
func SpinnerFreeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(formatRequest(r))
	decoder := json.NewDecoder(r.Body)
	var sm stypes.SMessageFree
	err := decoder.Decode(&sm)
	if err != nil {
		panic(err)
	}
	var calc calculation.Calculator
	calc.SetInput(sm.Content)
	calc.Calculate()
	sm.Spin1 = *calc.GetSpin1()
	sm.Spin2 = *calc.GetSpin2()
	sm.Content = "ok"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(sm)
	if err != nil {
		panic(err)
	}
}
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(formatRequest(r))
	decoder := json.NewDecoder(r.Body)
	var spmsg stypes.SignUpMessage
	err := decoder.Decode(&spmsg)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Println(spmsg.Mail, spmsg.Name, spmsg.Password)
	if spmsg.ContentCode == 20 {
		t, _ := SignUpNewUser(&spmsg)
		switch t {
		case -1:
			spmsg.ContentCode = 23
			spmsg.Client = "Server"
			spmsg.Content = "occur error when sign up"
			spmsg.Password = ""
			spmsg.Mail = ""
			spmsg.Name = ""
		case 1:
			spmsg.ContentCode = 21
			spmsg.Client = "Server"
			spmsg.Content = "your account created"
			spmsg.Password = ""
			spmsg.Mail = ""
			spmsg.Name = ""
		case 2:
			spmsg.ContentCode = 22
			spmsg.Client = "Server"
			spmsg.Content = "this account already created"
			spmsg.Password = ""
			spmsg.Mail = ""
			spmsg.Name = ""
		}
	} else {
		spmsg.ContentCode = 24
		spmsg.Client = "Server"
		spmsg.Content = "wrong request"
		spmsg.Password = ""
		spmsg.Mail = ""
		spmsg.Name = ""
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(spmsg)
	if err != nil {
		panic(err)
	}

}
func SignUpNewUser(info *stypes.SignUpMessage) (int, *stypes.User) {
	db := adapters.ConnectDB(2)
	defer db.Close()
	var user stypes.User
	row := db.QueryRow("SELECT UserID,Mail,IsConfirmed FROM users WHERE Mail = ?", info.Mail)
	err := row.Scan(&user.UserID, &user.Mail, &user.IsConfirmed)
	switch {
	case err == sql.ErrNoRows:
		r, _ := func() (int, *stypes.User) {
			var query string = "Insert into users (Name,Mail,Password) VALUES(?,?,?)"
			result, err := db.Exec(query, info.Name, info.Mail, info.Password)
			if err == nil {
				i, _ := result.RowsAffected()
				if i > 0 {
					return 1, nil
				} else {
					return -1, nil
				}
			} else {
				return -1, nil
			}
		}()
		return r, nil
	case err != nil:
		log.Println(err.Error())
		return -1, nil
	default:
		return 2, &user
	}

}

func YeniEkle() {
	var path = "spin_genel.txt"
	log.Println("basladı")
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	line, err := r.ReadString('\n')    // line defined once
	db := adapters.ConnectDB(3)
	defer db.Close()
	counter := 0
	for err != io.EOF {
		if line != "" {
			query := "Insert into generalspins (Spin) VALUES(?)"
			result, _ := db.Exec(query, line)
			i, err := result.RowsAffected()
			if err != nil {
				db.Close()
				db = adapters.ConnectDB(4)
			}
			counter += int(i)
			line, err = r.ReadString('\n') //  line was defined before
		}
	}
	log.Println(counter, "row effected")

}

func YeniEkle2() {
	var path = "yeni_eklenecek.txt"
	log.Println("basladı")
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	line, err := r.ReadString('\n')    // line defined once
	db := adapters.ConnectDB(3)
	defer db.Close()
	counter := 369094
	for err != io.EOF {
		if line != "" {
			rows := strings.Split(line, "|")
			for _, a := range rows {
				query := "Insert into spins (Spin,GroupID) VALUES(?,?)"
				a = strings.Replace(a, "{", "", -1)
				a = strings.Replace(a, "}", "", -1)
				result, _ := db.Exec(query, a, counter)
				_, err := result.RowsAffected()
				if err != nil {
					db.Close()
					db = adapters.ConnectDB(4)
				}
			}
			counter++
			line, err = r.ReadString('\n') //  line was defined before
		}
	}
	log.Println(counter, "row effected")

}
func YeniEkleKeyword() {
	var path = "keywords.txt"
	log.Println("basladı")
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	line, err := r.ReadString('\n')    // line defined once
	db := adapters.ConnectDB(3)
	defer db.Close()
	counter := 0;
	for err != io.EOF {
		if line != "" {
			query := "Insert into keywords (Keyword) VALUES(?)"
			result, _ := db.Exec(query, line)
			_, err := result.RowsAffected()
			if err != nil {
				db.Close()
				db = adapters.ConnectDB(4)
			}
			counter++
			line, err = r.ReadString('\n') //  line was defined before
		}
	}
	log.Println(counter, "row effected")

}
func CreateServer() {
	//server := packages.NewServer()
	//go server.Run()

}
func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "spinner:271312cs@/spinner")
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return db

}
func readFileWithReadString(fn string) (err error) {
	/*
	c := make(chan int, numCPU)  // Buffering optional but sensible.
	for i := 0; i < numCPU; i++ {
		go v.DoSome(i*len(v)/numCPU, (i+1)*len(v)/numCPU, u, c)
	}
	// Drain the channel.
	for i := 0; i < numCPU; i++ {
		<-c    // wait for one task to complete
	}
	*/


	db := ConnectDB()
	for {
		var err error
		var result string = ""
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		input = strings.Replace(input, "\n", "", -1)
		//input = strings.Replace(input, "'", "\\'", -1)
		//input = strings.Replace(input, ".", " ", -1)
		list := strings.Split(input, " ")
		start := time.Now()
		var buff bool = false
		var buffString string = ""
		var buffTemp string = ""
		for i, row := range list {
			_, rs := Search(&row, db)
			buffTemp = " " + *rs
			buff = false
			temp := row
			for count := 1; i + count < len(list) && count <= 1; count++ {
				temp += " " + list[i + count]
				r, rs := Search(&temp, db)
				if r {
					buff = true
					buffString = " " + *rs
				}
			}
			if buff {
				result += buffString
			} else {
				result += buffTemp
			}
		}
		fmt.Println(result)
		elapsed := time.Since(start)
		fmt.Println(elapsed)
	}

	if err != nil {
		return err
	}
	return
}
func Search(s *string, db *sql.DB) (bool, *string) {
	var Spin string
	//query := "SELECT Spin FROM generalspins where ( Spin like '%|" + string(*s) + "|%' or Spin like '{" + string(*s) + "|%' or Spin like '%|" + string(*s) + "}' ) and GeneralSpinsID < 30000;"
	err = db.QueryRow(
		"SELECT Spin FROM " +
			"(SELECT * FROM generalspins WHERE MATCH(Spin) AGAINST ('?' IN BOOLEAN MODE)) as results " +
			"WHERE ( Spin like '%|?|%' or Spin like '{?|%' or Spin like '%|?}' )",
		string(*s), string(*s), string(*s), string(*s)).Scan(&Spin)
	switch {
	case err == sql.ErrNoRows:
		//spin = line;
		return false, s
	case err != nil:
		log.Fatal(err)
		return false, s
	default:
		return true, GetString(&Spin, s)
	}
}
func GetString(s *string, orginal *string) *string {
	temp := []byte(*s)
	arr := strings.Split(string(temp[1:len(temp) - 3]), "|")
	var randomed int
	for {
		randomed = rand.Intn(len(arr))
		if arr[randomed] != *orginal {
			break
		}
	}
	return &arr[randomed]

}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")

}
/*func Listele() {

	s := " çÇöÖüÜıİşŞ"
	for _, r := range s {
		fmt.Printf("%c - %d\n", r, r)
	}
	str := " çÇöÖüÜıİşŞ"
	rune, _ := utf8.DecodeRuneInString(string(str[2]))
	fmt.Println(string(rune))

	for len(str) > 0 {

		r, size := utf8.DecodeRuneInString(str)
		fmt.Printf("%d %v\n", r, size)
		str = str[size:]
	}
}*/

const delim = "\"':?!.;,*(){}[]/%"

func TrimSuffix(s string) string {
	if strings.HasSuffix(s, "!") {
		s = s[:len(s) - len("!")]
	}
	return s
}
func TestDB() {
	db := adapters.ConnectDB(1)
	defer db.Close()
	db.Ping()
	/*var user stypes.User
	row := db.QueryRow("SELECT UserID,Mail FROM users WHERE Mail = ? AND Password = ?", *mail, *pass)
	*//*var p []byte
	err := db.QueryRow("SELECT data FROM message WHERE data->>'id'=$1", id).Scan(&p)
	if err != nil {
		// handle error
	}
	var m Message
	err := json.Unmarshal(p, &m)
	if err != nil {
		// handle error
	}*//*
	err := row.Scan(&user.UserID, &user.Mail)
	switch {
	case err == sql.ErrNoRows:
		return 3, nil
	case err != nil:
		log.Println(err.Error())
		return 3, nil
	default:

	}*/
}
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	//YeniEkleKeyword()
	//TestDB()
	//Listele();
	//YeniEkle2()
	//var input = "istanbul'da kentsel dönüşüm kapsamında yıkılmakta olan eski binalar örnek gösterilerek asbest maruziyetinin bizler için yeni olmadığı, havagazı fabrikası faciasının büyütüldüğü dillendiriliyor. yıkımına başlanan bina alelade bir apartman değil; havagazı üretmek üzere inşaa edilmiş, dolayısıyla ısı yalıtımının önem arz ettiği ve 350 ton kadar asbestin her tarafında bol kepçeden kullanıldığı atıl durumdaki bir fabrikadır ve böyle bir fabrikanın usule aykırı biçimde yıkılması halinde etrafa saçılacak asbestin etkileri de yıkıcı olacaktır. yıkımı durdurulan havagazı fabrikasının bulunduğu arazi tez elden çevrilerek karantina altına alınmalıdır. sıfır noktasından ve civar semtlerden düzenli aralıklarla numuneler alınarak, serbest haldeki asbest lifi miktarının ölçümleri yapılmalı ve limit değeri aşan bölgelerde vatandaşların tahliyesi gerçekleştirilmelidir.olayın toplumsal histeri boyutuna geldiğini iddia eden kimselerin ne yazık ki asbest hususunda henüz yeterli bilince ulaşmamış oldukları kanaatindeyim. hepimiz hayatımıza kaldığımız yerden devam etmek istiyoruz. yetkili ağızlardan rahatlatıcı açıklamalar duymak, ertesi sabah çocuklarımızı gönül rahatlığıyla okul servisine bindirmek, işimize odaklanmak, soluduğumuz havanın bizi ölüme hızla yaklaştırmadığı geçmişe dönmek istiyoruz. uyuşmak istiyoruz; zira zihnimizin uzun süreli yoğun stres altında sağlıklı biçimde işleyişini sürdürmesi mümkün değil. akıl sağlığımızı koruyabilmemiz için olayın gerçekliğini ve etkilerini gözardı etmek, onu reddetmek, onu unutmak eğilimindeyiz. fakat vaziyeti olduğu gibi kabul edip göz göre göre zehir solumamız, hayatlarımıza pervasızca devam etmemiz, ömrümüzün bir kısmından vazgeçmemiz kabul edilemez. uyanık kalmalıyız. ankara halkının bu faciayı minimum hasarla atlatmasını sağlayıcı tüm tedbirlerin alınması elzemdir. öyle durumlar vardır ki; sizi radikal kararlar almak zorunda bırakır. bir hortum, bir tsunami, yahut bir yanardağda hareketlilik beklendiğinde afetten etkilenme ihtimali olan insanlar bölgeden ivedilikle tahliye edilir. akl-ı selim sahibi hiçbir kimse yoktur ki başına gelecek bir felaketin varlığını ve gerçekliğini idrak ettiği halde harekete geçmesin, hiçbir önlem almasın. mimarlar odası ankara şubesi ve kimyagerler odasının yıkımı gerçekleştirilen binanın çevresinde yapmış oldukları ölçümler ve kamuoyuyla paylaştıkları veriler neticesinde, ankara'da bir felaketin var ve yaşanmakta olduğundan emin olabilirsiniz. ankara'da vuku bulan bu olay, herhangi bir doğal afetten daha az ciddi değildir. bir doğal afetin sonuçları gerçek zamanlı olarak tecrübe edilmekte iken; asbest denen bu musibet sinsice ilerleyecektir ve neticelerini insanlar ancak seneler sonra çok acı bir biçimde tecrübe edecektir, tek fark budur.ihmalkarlık, iş bilmezlik, liyakatsizlik, şark kurnazlığı ve cahil cesareti sonucu başkent ankara'nın göbeğinde meydana gelen bu vahim tablonun ankara halkına daha fazla zarar vermemesi için olaya bir afet ciddiyetiyle yaklaşmak elzemdir. bölgede aldığımız her nefes, havada asılı duran asbest liflerine biraz daha maruz kalmamıza, ciddi solunum hastalıklarının vücudumuzda vuku bulma riskini artırmamıza sebep olacaktır. ankara halkı, kendi hayatları ve çocuklarının hayatları için harekete geçmelidir. birey olarak sorumluluklarımız etrafımızdakileri bilinçlendirmek; maltepe, kızılay ve civarındaki semtlerde koruyucu maske kullanmak, içinde en azından biraz vicdan kırıntısı taşıdığına inandığınız yetkililerin açıklamalarını takip etmek ve talimatlarını dinlemek, sosyal medya yoluyla kamuoyu oluşturmaya yardımcı olmak, ilgili kurum ve kuruluşları telefon ve email yoluyla göreve davet etmek, havagazı fabrikasının yıkımının teknik yeterliğe sahip bir yıkım firması tarafından gerekli tüm tedbirler alınarak usule uygun biçimde sürdürülmesini talep etmek ve tatmin edici bir netice alana kadar pes etmemek, durumu kanıksamamaktır.asbest konusunda toplumsal farkındalık sağlanması ve yalnızca bu olay özelinde değil; türkiye genelinde asbest ihtiva eden tüm binaların yıkım işlemlerinde azami ölçüde tedbir alınmasının ve yıkımın mevzuata uygun olarak gerçekleştirilmesinin talep edilmesi hususunda hepimize büyük sorumluluk düşmektedir. "

	/*var calc calculation.Calculator
	for {
		var err error
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		//input = "istanbul'da kentsel dönüşüm kapsamında yıkılmakta olan eski binalar örnek gösterilerek asbest maruziyetinin bizler için yeni olmadığı, havagazı fabrikası faciasının büyütüldüğü dillendiriliyor. yıkımına başlanan bina alelade bir apartman değil; havagazı üretmek üzere inşaa edilmiş, dolayısıyla ısı yalıtımının önem arz ettiği ve 350 ton kadar asbestin her tarafında bol kepçeden kullanıldığı atıl durumdaki bir fabrikadır ve böyle bir fabrikanın usule aykırı biçimde yıkılması halinde etrafa saçılacak asbestin etkileri de yıkıcı olacaktır. yıkımı durdurulan havagazı fabrikasının bulunduğu arazi tez elden çevrilerek karantina altına alınmalıdır. sıfır noktasından ve civar semtlerden düzenli aralıklarla numuneler alınarak, serbest haldeki asbest lifi miktarının ölçümleri yapılmalı ve limit değeri aşan bölgelerde vatandaşların tahliyesi gerçekleştirilmelidir.olayın toplumsal histeri boyutuna geldiğini iddia eden kimselerin ne yazık ki asbest hususunda henüz yeterli bilince ulaşmamış oldukları kanaatindeyim. hepimiz hayatımıza kaldığımız yerden devam etmek istiyoruz. yetkili ağızlardan rahatlatıcı açıklamalar duymak, ertesi sabah çocuklarımızı gönül rahatlığıyla okul servisine bindirmek, işimize odaklanmak, soluduğumuz havanın bizi ölüme hızla yaklaştırmadığı geçmişe dönmek istiyoruz. uyuşmak istiyoruz; zira zihnimizin uzun süreli yoğun stres altında sağlıklı biçimde işleyişini sürdürmesi mümkün değil. akıl sağlığımızı koruyabilmemiz için olayın gerçekliğini ve etkilerini gözardı etmek, onu reddetmek, onu unutmak eğilimindeyiz. fakat vaziyeti olduğu gibi kabul edip göz göre göre zehir solumamız, hayatlarımıza pervasızca devam etmemiz, ömrümüzün bir kısmından vazgeçmemiz kabul edilemez. uyanık kalmalıyız. ankara halkının bu faciayı minimum hasarla atlatmasını sağlayıcı tüm tedbirlerin alınması elzemdir. öyle durumlar vardır ki; sizi radikal kararlar almak zorunda bırakır. bir hortum, bir tsunami, yahut bir yanardağda hareketlilik beklendiğinde afetten etkilenme ihtimali olan insanlar bölgeden ivedilikle tahliye edilir. akl-ı selim sahibi hiçbir kimse yoktur ki başına gelecek bir felaketin varlığını ve gerçekliğini idrak ettiği halde harekete geçmesin, hiçbir önlem almasın. mimarlar odası ankara şubesi ve kimyagerler odasının yıkımı gerçekleştirilen binanın çevresinde yapmış oldukları ölçümler ve kamuoyuyla paylaştıkları veriler neticesinde, ankara'da bir felaketin var ve yaşanmakta olduğundan emin olabilirsiniz. ankara'da vuku bulan bu olay, herhangi bir doğal afetten daha az ciddi değildir. bir doğal afetin sonuçları gerçek zamanlı olarak tecrübe edilmekte iken; asbest denen bu musibet sinsice ilerleyecektir ve neticelerini insanlar ancak seneler sonra çok acı bir biçimde tecrübe edecektir, tek fark budur.ihmalkarlık, iş bilmezlik, liyakatsizlik, şark kurnazlığı ve cahil cesareti sonucu başkent ankara'nın göbeğinde meydana gelen bu vahim tablonun ankara halkına daha fazla zarar vermemesi için olaya bir afet ciddiyetiyle yaklaşmak elzemdir. bölgede aldığımız her nefes, havada asılı duran asbest liflerine biraz daha maruz kalmamıza, ciddi solunum hastalıklarının vücudumuzda vuku bulma riskini artırmamıza sebep olacaktır. ankara halkı, kendi hayatları ve çocuklarının hayatları için harekete geçmelidir. birey olarak sorumluluklarımız etrafımızdakileri bilinçlendirmek; maltepe, kızılay ve civarındaki semtlerde koruyucu maske kullanmak, içinde en azından biraz vicdan kırıntısı taşıdığına inandığınız yetkililerin açıklamalarını takip etmek ve talimatlarını dinlemek, sosyal medya yoluyla kamuoyu oluşturmaya yardımcı olmak, ilgili kurum ve kuruluşları telefon ve email yoluyla göreve davet etmek, havagazı fabrikasının yıkımının teknik yeterliğe sahip bir yıkım firması tarafından gerekli tüm tedbirler alınarak usule uygun biçimde sürdürülmesini talep etmek ve tatmin edici bir netice alana kadar pes etmemek, durumu kanıksamamaktır.asbest konusunda toplumsal farkındalık sağlanması ve yalnızca bu olay özelinde değil; türkiye genelinde asbest ihtiva eden tüm binaların yıkım işlemlerinde azami ölçüde tedbir alınmasının ve yıkımın mevzuata uygun olarak gerçekleştirilmesinin talep edilmesi hususunda hepimize büyük sorumluluk düşmektedir. "
		//input = "Havalimanı'nda unutulan binlerce cep telefonu, laptop, bilgisayar, fotoğraf makinası gibi binlerce bavul dolu eşya yarın satışa çıkıyor.THY Dış Hatlarda unutulan eşyaların ihalesi geçtiğimiz günlerde yapıldı. Her yıl ihaleyi alan Alican Korkmaz, bu yıl da ihaleyi aldı. Avcılar Çağatay İş Merkezi'nde yarın sabah başlayacak olan satışta, binlerce son model cep telefonları, laptop, bilgisayar, fotoğraf makinası, bebek arabaları, bebek araç koltukları, gözlükler, parfümler gibi on binlerce ürün yarı fiyatına satışa çıkacak.Satışı gerçekleştiren Alican Kormaz, \"Yarın sabah satışa çıkaracağımız ürünler piyasanın yarı fiyatından da çok düşüğüne alıcılarını bekliyor\" diye konuştu.Her şey mükemmeldi. Ta ki Dubai'ye ayak basana kadar. Ben bu kadar kompleksli bir polis takımı yada havaalanı çalışanı görmedim. Sadece Fast Track çizgisinden bir tık dışarı çıktık diye, onunla ilgili memura soru sorduk diye karakolda bekletildik. Fast track kartlarımız elimizden alındı. Bize geri verilmedi. Yüzümüze dahi bakılmadı. Aşağılandık. Bir küfür yemediğimiz kaldı.  Dubai havalimanındaki polisler bana ilkokul öğretmenimi hatırlatıyor. Konuşma! Gülme! Bakma! Kalkma! İnanılır gibi değilsiniz... 4 saat sonra çıkabildik Dubai Havalanı'ndan. Tam 4 saat sonra! Sebebini bilmiyoruz! Kesinlikle bilmiyoruz! Hiçbir şey söylemiyorlar... sonra gelip özür dilediler bizden. \"Kadir Doğulu bir süre sonra ikinci kez aynı muameleye maruz kaldığını duyurdu:"
		input = "  THY Dış Hatlarda unutulan eşyaların ihalesi geçtiğimiz günlerde yapıldı. Her yıl ihaleyi alan Alican Korkmaz, bu yıl da ihaleyi aldı. Avcılar Çağatay İş Merkezi'nde yarın sabah başlayacak olan satışta, binlerce son model cep telefonları, laptop, bilgisayar, fotoğraf makinası, bebek arabaları, bebek araç koltukları, gözlükler, parfümler gibi on binlerce ürün yarı fiyatına satışa çıkacak. Satışı gerçekleştiren Alican Kormaz, \"Yarın sabah satışa çıkaracağımız ürünler piyasanın yarı fiyatından da çok düşüğüne alıcılarını bekliyor\" diye konuştu."

		if err != nil {
			log.Println(err)
		}
		input = strings.Replace(input, "\n", "", -1)
		calc.SetInput(input)
		calc.Calculate()
	}*/
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(*addrSProvider, ServiceProviderHandler)
	router.HandleFunc(*addrWSsocket, WebSocketHandler)
	router.HandleFunc(*addrSpinner, SpinnerHandler)
	router.HandleFunc(*addrSignUp, SignupHandler)
	router.HandleFunc(*addrSpinFree, SpinnerFreeHandler)
	CreateServer()
	err := http.ListenAndServe(*addr, router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//readFileWithReadString("inputlarge.txt")

}
