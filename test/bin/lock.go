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
	for i := 10; i > 0; i-- {
		wait.Add(1)
		go func(i int) {
			defer wait.Done()
			//r := rand.New(rand.NewSource(time.Now().UnixNano()))
			//num := r.Intn(128)
			//time.Sleep(time.Duration(num))
			if i == 5 {
				lockOk := higo.Retry(&higo.Returner{Interval: 1 * time.Second, Retry: 4},
					&higo.Locker{Key: "test", Timeout: 1 * time.Second},
					func() {
						fmt.Println(fmt.Sprintf("%d开始执行锁", i))
						time.Sleep(3 * time.Second)
						fmt.Println(fmt.Sprintf("%d执行完成", i))
					})
				fmt.Printf("lockOk: %v %d\n", lockOk, i)
			} else {
				lockOk := higo.Lock(&higo.Locker{Key: "test", Timeout: 1 * time.Second}, func() {
					fmt.Println(fmt.Sprintf("%d开始执行锁", i))
					time.Sleep(2 * time.Second)
					fmt.Println(fmt.Sprintf("%d执行完成", i))
				})
				fmt.Printf("lockOk: %v %d\n", lockOk, i)
			}
		}(i)
	}
	/*
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
		}()*/
	wait.Wait()

}
