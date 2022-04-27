// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Program creates customers for the simulation of the sleeping barber problem
// implemented in the shop package.



/*
经典问题:
理发师问题描述：

（1）理发店里有一位理发师、一把理发椅和n把供等候理发的顾客坐的椅子
（2）如果没有顾客，理发师便在理发椅上睡觉
（3）一个顾客到来时，它必须叫醒理发师
（4）如果理发师正在理发时又有顾客来到，则如果有空椅子可坐，就坐下来等待，否则就离开

问题分析：

1、对于理发师问题而言，是生产者-消费者（有界缓冲区）模型的一种。其中理发师和顾客之间涉及到进程之间的同步问题，理发师是生产者，顾客是消费者，生产者生产的速度（理发师理发的速度），和消费者消费的速度（顾客来到理发店的时间），这两者肯定是不同的。那么，题目中就涉及到了一个有界缓冲区，即有N把顾客可以坐着等候理发师的椅子。如果顾客来的太快了，就可以先坐在椅子上等候一下理发师，但是如果椅子坐满了，这时候顾客就直接走，不理发了。这个N是有界缓冲区的大小，如果缓冲区放不下消费者了，消费者就不进行消费。

2、同样的，在生产者-消费者（有界缓冲区）模型中，还存在进程之间的互斥，比如多个消费者同时访问缓冲区，那么肯定会改变缓冲区的状态，缓冲区就是临界资源，多个消费者不能同时去改变缓冲区的状态。在这个问题上，就相当于执行顾客任务的进程，就必须有互斥的操作，同样的，理发师改变缓冲区状态的操作也需要互斥。这个问题中，缓冲区的状态，就是还剩多少个在等待的顾客，顾客来一个，肯定等待理发的顾客数目就+1，理发师理一次发，等待理发的顾客数目就-1。
————————————————
版权声明：本文为CSDN博主「不会code的菜鸟」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/CLZHIT/article/details/113799912
*/

package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"./shop"
)

func main() {
	const maxChairs = 10
	s := shop.Open(maxChairs)

	// Create a goroutine that is constantly, but inconsistently, generating
	// customers who are entering the shop.
	go func() {
		var id int64
		for {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			name := fmt.Sprintf("cust-%d", atomic.AddInt64(&id, 1))
			if err := s.EnterCustomer(name); err != nil {
				fmt.Printf("Customer %q told %q\n", name, err)
				if err == shop.ErrShopClosed {
					break
				}
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("Shutting down shop")
	s.Close()
}
