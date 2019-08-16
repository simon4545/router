package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var OnlineIps=[]string{}
func IsOnline(ip string) bool{
	cmd := exec.Command("ping", "-c","1", "-W","1",ip)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	if runtime.GOOS=="darwin"{
		return strings.Contains(string(opBytes), "1 packets received")
	}else if runtime.GOOS=="linux"{
		return strings.Contains(string(opBytes), "64 bytes")
	}
	return false
}
func IsPingable(wg *sync.WaitGroup,ip string,ips *[]string){
	defer wg.Done()
	cmd := exec.Command("ping", "-c","1", "-W","1",ip)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	if runtime.GOOS=="darwin"{
		if strings.Contains(string(opBytes), "1 packets received")==true{
			*ips=append(*ips, ip)
			//fmt.Println("扫描到IP:",ip)
		}else{
			//fmt.Println("IP嗅探失败:",ip)
		}
	}else if runtime.GOOS=="linux"{
		if strings.Contains(string(opBytes), "64 bytes")==true{
			*ips=append(*ips, ip)
			//fmt.Println("扫描到IP:",ip)
		}else{
			//fmt.Println("IP嗅探失败:",ip)
		}
	}
}

/**
扫描局域网
 */
func LanScan() (ips []string) {
	var threadGroup = sync.WaitGroup{}
	threadGroup.Add(255-2)
	for i := 2; i < 255; i++ {
		//fmt.Println("正在扫描IP:","192.168.2."+strconv.Itoa(i))
		go IsPingable(&threadGroup,"192.168.100."+strconv.Itoa(i),&ips)
	}
	//等待所有线程完成
	threadGroup.Wait()
	sort.Strings(ips)
	fmt.Println(ips)
	return
}

/*
比较所有的Ip并返回差集
OnlineIps是上次在线的IP
 */
func LanCompare() (diff []string) {
	scanIps:=LanScan()
	if len(OnlineIps)<len(scanIps){
		OnlineIps=scanIps
		return
	}
	return Difference(OnlineIps,scanIps)
}

func ReDial() bool{
	cmd := exec.Command("/etc/init.d/network", "restart")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return false
	}
	// 读取输出结果
	_, err = ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func SendMail(title string,content string) {
	auth := smtp.PlainAuth("", "simon4545@qq.com", "simon518418", "smtp.qq.com")
	to := []string{"xlzhou@forke.cn","alinger@forke.cn","simon4545@126.com","1176877783@qq.com"}
	nickname := "simon4545"
	user := "simon4545@qq.com"
	subject := title+"网络故障"
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := content
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail("smtp.qq.com:25", auth, user, to, msg)
	if err != nil {
		fmt.Printf("send mail error: %v", err)
	}
}


func HostName() string{
	host, err := os.Hostname()
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		fmt.Printf("%s", host)
	}
	return host
}