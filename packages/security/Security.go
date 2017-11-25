package security

import "crypto/rand"

func GetRandomAPIKey() string {
	var dictionary string
	dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 35)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v % byte(len(dictionary))]
	}
	return string(bytes)
}
