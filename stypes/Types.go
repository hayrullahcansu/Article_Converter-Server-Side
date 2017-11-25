package stypes

type Message struct {
	Client      string `json:"client"`
	ContentCode int `json:"content_code"`
	Content     string `json:"content"`
}
type SignUpMessage struct {
	Client      string `json:"client"`
	ContentCode int `json:"content_code"`
	Content     string `json:"content"`
	Name        string `json:"name"`
	Mail        string `json:"mail"`
	Password    string `json:"password"`
}
type RMessage struct {
	Client      string `json:"client"`
	ContentCode int `json:"content_code"`
	Content     string `json:"content"`
	APIKey      string `json:"api_key"`
}
type WSMessage struct {
	Client      string `json:"client"`
	ContentCode int `json:"content_code"`
	Content     string `json:"content"`
	APIKey      string `json:"api_key"`
}
type SPMessage struct {
	Client      string `json:"client"`
	ContentCode int `json:"content_code"`
	Content     string `json:"content"`
	UserID      int `json:"user_id"`
	APIKey      string `json:"api_key"`
}
type SMessage struct {
	Client      string `json:"client"`
	ContentCode int `json:"content_code"`
	UserID      int `json:"user_id"`
	APIKey      string `json:"api_key"`
	Content     string `json:"content"`
	Preview1    string `json:"preview1"`
	Preview2    string `json:"preview2"`
	Spin1       string `json:"spin1"`
	Spin2       string `json:"spin2"`
}
type SMessageFree struct {
	Content string `json:"content"`
	Spin1   string `json:"spin1"`
	Spin2   string `json:"spin2"`
}
type ClientInfo struct {
	Info     string `json:"info"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
}
type User struct {
	UserID       int `json:"UserID"`
	Mail         string `json:"Mail"`
	Password     string `json:"Password"`
	RegisterDate string `json:"RegisterDate"`
	IsConfirmed  bool `json:"IsConfirmed"`
}

func (self *Message) String() string {
	return self.Client + " says " + self.Content
}