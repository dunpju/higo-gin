package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

var (
	crtFileName = "higo.crt.pem"
	keyFileName = "higo.key.pem"
)

type Ssl struct {
	outDir string
	crt string
	key string
}

// 构造函数
func NewSsl(outDir string, crt string, key string) *Ssl {
	return &Ssl{outDir,crt,key}
}

// 生成ssl证书
func (this *Ssl) Generate() {

	// 创建输出目录
	this.createOutDir()

	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)

	// 定义：引用IETF的安全领域的公钥基础实施（PKIX）工作组的标准实例化内容
	subject := pkix.Name{
		Organization:       []string{"WWW.YUMI.COM"},
		OrganizationalUnit: []string{"ITs"},
		CommonName:         "YUMI.COM Web",
	}

	// 设置 SSL证书的属性用途
	certificate509 := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(100 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	// 生成指定位数密匙
	pk, _ := rsa.GenerateKey(rand.Reader, 1024)

	// 生成 SSL公匙
	derBytes, _ := x509.CreateCertificate(rand.Reader, &certificate509, &certificate509, &pk.PublicKey, pk)
	certOut, _ := os.Create(this.crt)
	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	_ = certOut.Close()

	// 生成 SSL私匙
	keyOut, _ := os.Create(this.key)
	_ = pem.Encode(keyOut, &pem.Block{Type: "RAS PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	_ = keyOut.Close()
}

// 更新ssl证书
func (this *Ssl) Update()  {

	// 创建输出目录
	this.createOutDir()

	// 判断创建时间
	crtFile := NewFile(this.crt)
	if crtFile.IsExist() {

		createTimestamp := crtFile.GetCreateTimestamp()

		if (CurrentTimestamp() - createTimestamp) > new(RandomUtil).IntHour24ToSecond() {
			fmt.Println("重新生成证书")
			this.Generate() // 重新生成证书
		}
	}
	fmt.Println("更新ssl证书")
}

// 创建输出目录
func (this *Ssl) createOutDir()  {
	// 目录不存在，并创建
	if _, err := os.Stat(this.outDir); os.IsNotExist(err) {
		if os.Mkdir(this.outDir, os.ModePerm) != nil {}
	}

	this.crt = this.outDir + crtFileName
	this.key = this.outDir + keyFileName
}
