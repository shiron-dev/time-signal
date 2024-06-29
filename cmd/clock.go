package cmd

import (
	"fmt"
	"time"
)

// 関数を呼び出すタイミング
var intervals = []int{0, 10, 20, 30, 40, 50}

func Timer(onSecondChange func(seconds int), callFunction func(time time.Time)) {
	var prevSecond int

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("タイムゾーンのロードに失敗しました:", err)
		return
	}

	for {
		now := time.Now().In(loc)
		seconds := now.Second()

		if seconds != prevSecond {
			onSecondChange(seconds)
			prevSecond = seconds
		}

		for _, interval := range intervals {
			if seconds == interval {
				callFunction(now)
				break
			}
		}

		time.Sleep(time.Second)
	}
}
