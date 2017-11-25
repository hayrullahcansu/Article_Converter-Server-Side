package packages

import (
	"../stypes"
	"log"
)

type APIRegister struct {
	users map[int]string
}

func NewRegister() *APIRegister {
	return &APIRegister{
		users:    make(map[int]string),
	}
}
func (h *APIRegister) CheckRegister(user *stypes.User) bool {
	var api string = h.users[user.UserID]
	if api == "" {
		log.Println("yok")
		return false
	} else {
		log.Println("var")
		return true
	}
}
func (h *APIRegister) CompareKeyUserID(key *string, id *int) int {
	var api string = h.users[*id]
	if api == "" {
		return -2
	} else if api == *key {
		return 1
	} else {
		return -1
	}

}
func (h *APIRegister) AddRegister(user *stypes.User, apikey *string) {
	h.users[user.UserID] = *apikey
}

func (h *APIRegister) RemoveRegister(user *stypes.User) {
	delete(h.users, user.UserID)
}