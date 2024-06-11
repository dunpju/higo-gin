package main

import (
	"fmt"
	"github.com/dunpju/higo-gin/higo"
	"sync"
	"time"
)

func main() {
	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		defer wait.Done()
		higo.Lock("test", func() {
			fmt.Println("开始执行锁1")
			time.Sleep(10)
			fmt.Println("执行完成1")
		})
	}()
	wait.Add(1)
	go func() {
		defer wait.Done()
		higo.Lock("test", func() {
			fmt.Println("开始执行锁2")
			time.Sleep(1)
			fmt.Println("执行完成2")
		})
	}()
	wait.Wait()

}
