package main

import (
	"fmt"
	"github.com/dunpju/higo-gin/higo"
	"sync"
	"time"
)

func main() {
	m := &sync.Map{}
	v, ok := m.LoadOrStore("j", &sync.Mutex{})
	fmt.Println(fmt.Sprintf("%p", v), ok)
	v1, ok1 := m.LoadOrStore("j", &sync.Mutex{})
	fmt.Println(fmt.Sprintf("%p", v1), ok1)

	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		defer wait.Done()
		//r := rand.New(rand.NewSource(time.Now().UnixNano()))
		//num := r.Intn(128)
		//time.Sleep(time.Duration(num))
		higo.Lock("test", func() {
			fmt.Println("开始执行锁3")
			time.Sleep(1)
			fmt.Println("执行完成3")
		})
	}()
	wait.Add(1)
	go func() {
		defer wait.Done()
		//r := rand.New(rand.NewSource(time.Now().UnixNano()))
		//num := r.Intn(128)
		//time.Sleep(time.Duration(num))
		//higo.Retry(100*time.Millisecond, 3, "test", func() {
		//	fmt.Println("开始执行锁2")
		//	time.Sleep(1)
		//	fmt.Println("执行完成2")
		//})
		higo.Lock("test", func() {
			fmt.Println("开始执行锁2")
			time.Sleep(1)
			fmt.Println("执行完成2")
		})
	}()
	wait.Add(1)
	go func() {
		defer wait.Done()
		//r := rand.New(rand.NewSource(time.Now().UnixNano()))
		//num := r.Intn(128)
		//time.Sleep(time.Duration(num))
		higo.Lock("test", func() {
			fmt.Println("开始执行锁1")
			time.Sleep(10)
			fmt.Println("执行完成1")
		})
	}()
	wait.Wait()

}
