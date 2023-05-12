package main

import (
	"crypto/rand"
	"crypto/rc4"
	"fmt"
	"math/big"
)

func getPrime(n int) *big.Int {
	q, err := rand.Prime(rand.Reader, n)
	if err != nil {
		panic(err)
	}
	return q
}

func isPrime(n *big.Int) bool {
	return n.ProbablyPrime(20)
}

// 判断一个数是否为素数的原根
func isPrimitiveRoot(a, p *big.Int) bool {
	phi := big.NewInt(0).Sub(p, big.NewInt(1))

	// 遍历
	for i := big.NewInt(2); i.Cmp(phi) < 0; i.Add(i, big.NewInt(1)) {
		if big.NewInt(0).Exp(a, i, p).Cmp(big.NewInt(1)) == 0 {
			return false
		}
	}
	return true
}

// 生成随机数
func randInt(min, max *big.Int) *big.Int {
	n, err := rand.Int(rand.Reader, big.NewInt(0).Sub(max, min))
	if err != nil {
		panic(err)
	}
	return n.Add(n, min)
}

// 生成素数p的原根
func genPrimitiveRoot(p *big.Int) *big.Int {
	phi := big.NewInt(0).Sub(p, big.NewInt(1))
	// 随机生成一个数 a
	for {
		a := randInt(big.NewInt(2), phi)

		// 如果 a 是 p 的原根，则返回 a
		if isPrimitiveRoot(a, p) {
			return a
		}
	}
}

type human struct {
	name       string
	publicKey  *big.Int
	privateKey *big.Int
	key        *big.Int
}

func main() {
	p := getPrime(16)
	g := genPrimitiveRoot(p)
	fmt.Printf("p value:%v\ng value:%v\n", p, g)
	Alice := new(human)
	Bob := new(human)
	Alice.privateKey = randInt(big.NewInt(2), p)
	Bob.privateKey = randInt(big.NewInt(2), p)
	fmt.Printf("Alice.privateKey:%v\nBob.privateKey:%v\n", Alice.privateKey, Bob.privateKey)
	Alice.publicKey = new(big.Int).Exp(g, Alice.privateKey, p)
	// fmt.Printf("pkey1:%v\n", pkey1)
	Bob.publicKey = new(big.Int).Exp(g, Bob.privateKey, p)
	fmt.Printf("Alice.publicKey:%v\nBob.publicKey:%v\n", Alice.publicKey, Bob.publicKey)
	Alice.key = new(big.Int).Exp(Bob.publicKey, Alice.privateKey, p)
	Bob.key = new(big.Int).Exp(Alice.publicKey, Bob.privateKey, p)
	fmt.Printf("Alice.key:%v\nBob.key:%v\n", Alice.key, Bob.key)

	plaintext := []byte("Hello, world!")

	// 使用Alice的密钥进行RC4加密
	cipher, err := rc4.NewCipher(Alice.key.Bytes())
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, len(plaintext))
	cipher.XORKeyStream(ciphertext, plaintext)
	fmt.Printf("Ciphertext: %x\n", ciphertext)

	// 使用Bob的密钥进行RC4解密
	decipher, err := rc4.NewCipher(Bob.key.Bytes())
	if err != nil {
		panic(err)
	}
	decrypted := make([]byte, len(ciphertext))
	decipher.XORKeyStream(decrypted, ciphertext)
	fmt.Printf("Decrypted: %s\n", decrypted)

}
