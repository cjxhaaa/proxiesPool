package main

import (
	"ProxyPool/pkg/secret"
	"fmt"
	"os"
	"strings"
)

func main() {
	origin := strings.TrimSpace(os.Args[1])
	key := "1234567812345678"
	fmt.Println("origin:",origin)

	encryptCode := secret.AesEncrypt(origin, key)
	fmt.Println("encrypt result:",encryptCode)

	decryptCode := secret.AesDecrypt(encryptCode, key)
	fmt.Println("decrypt result:", decryptCode)
}
