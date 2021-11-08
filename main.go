package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const timeTemplate1 = "2006-01-02 15:04:05"

//go:embed img/time_icon.jpg
var icon []byte

func main() {
	myApp := app.New()
	setting := myApp.Settings()
	setting.SetTheme(theme.LightTheme())
	resource := fyne.NewStaticResource("time_icon", icon)
	myApp.SetIcon(resource)
	myWindow := myApp.NewWindow("Time")

	timerGridElems3, ticker := TimeNow()
	defer ticker.Stop()
	gridElems := TimestampToDate()
	gridElems2 := DateToTimestamp()
	gridElems3 := DayToSec()
	gridElems4 := SecToDayHMS()
	timerGridElems3 = append(timerGridElems3, gridElems...)
	timerGridElems3 = append(timerGridElems3, gridElems2...)
	timerGridElems3 = append(timerGridElems3, gridElems3...)
	timerGridElems3 = append(timerGridElems3, gridElems4...)
	grid := container.New(layout.NewGridLayout(3), timerGridElems3...)
	myWindow.SetContent(grid)
	myWindow.Resize(fyne.NewSize(500, 50))
	myWindow.ShowAndRun()
}

func TimeNow() ([]fyne.CanvasObject, *time.Ticker) {
	now := time.Now().UnixNano() / 1e6
	millionSec := strconv.Itoa(int(now))
	timeStampInp := widget.NewLabel(millionSec)
	timeStampInp.SetText(strconv.Itoa(int(now)))
	text3 := widget.NewLabel("DATE")
	nums, err := strconv.Atoi(timeStampInp.Text)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	date, err := timeStampToDate(nums)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	text3.SetText(date)
	ticker := time.NewTicker(time.Second)
	go func(label, text3 *widget.Label) {
		for t := range ticker.C {
			millionSec := strconv.Itoa(int(t.UnixNano() / 1e6))
			label.SetText(millionSec)
			nums, err := strconv.Atoi(timeStampInp.Text)
			if err != nil {
				fmt.Println(err)
				return
			}
			date, err := timeStampToDate(nums)
			if err != nil {
				return
			}
			text3.SetText(date)
		}
	}(timeStampInp, text3)
	click1 := widget.NewButton("COPY", func() {
		//复制到剪切版
		copyClipBoard(text3.Text)
	})
	click1.SetIcon(theme.ContentCopyIcon())
	return []fyne.CanvasObject{
		timeStampInp, click1, text3,
	}, ticker
}

func TimestampToDate() []fyne.CanvasObject {
	timeStampInp := widget.NewEntry()
	now := time.Now().UnixNano() / 1e6
	timeStampInp.SetText(strconv.Itoa(int(now)))
	timeStampInp.SetPlaceHolder("TIMESTAMP")
	text3 := widget.NewLabel("DATE")
	click1 := widget.NewButton("CLICK", func() {
		nums, err := strconv.Atoi(timeStampInp.Text)
		if err != nil {
			fmt.Println(err)
			return
		}
		date, err := timeStampToDate(nums)
		if err != nil {
			return
		}
		text3.SetText(date)
		//复制到剪切版
		copyClipBoard(date)
		fmt.Printf("%d--->%s\n", nums, date)
	})
	return []fyne.CanvasObject{
		timeStampInp, click1, text3,
	}
}

func DateToTimestamp() []fyne.CanvasObject {
	timeStampInp := widget.NewEntry()
	now := time.Now().UnixNano() / 1e6
	date, _ := timeStampToDate(int(now))
	timeStampInp.SetText(date)
	timeStampInp.SetPlaceHolder("DATE")
	text3 := widget.NewLabel("TIMESTAMP")

	click1 := widget.NewButton("CLICK", func() {
		stamp, err := time.ParseInLocation(timeTemplate1, timeStampInp.Text, time.Local)
		if err != nil {
			fmt.Println(err)
			return
		}
		millionSec := strconv.Itoa(int(stamp.UnixNano() / 1e6))
		text3.SetText(millionSec)
		copyClipBoard(millionSec)
		fmt.Printf("%s--->%s\n", stamp, millionSec)
	})
	return []fyne.CanvasObject{
		timeStampInp, click1, text3,
	}
}

func SecToDayHMS() []fyne.CanvasObject {
	timeStampInp := widget.NewEntry()
	timeStampInp.SetPlaceHolder("SECONDS")
	text3 := widget.NewLabel("D H M S")

	click1 := widget.NewButton("CLICK", func() {
		seconds, err := strconv.Atoi(timeStampInp.Text)
		if err != nil {
			fmt.Println(err)
			return
		}
		str := secToDayHMS(seconds)
		text3.SetText(str)
		copyClipBoard(str)
		fmt.Printf("%d--->%s\n", seconds, str)
	})
	return []fyne.CanvasObject{
		timeStampInp, click1, text3,
	}
}

func DayToSec() []fyne.CanvasObject {
	timeStampInp := widget.NewEntry()
	timeStampInp.SetPlaceHolder("D H M S")
	text3 := widget.NewLabel("SECONDS")

	click1 := widget.NewButton("CLICK", func() {
		str := dayHMSToSec(timeStampInp.Text)
		if str == "" {
			return
		}
		text3.SetText(str)
		copyClipBoard(str)
		fmt.Printf("%s--->%s\n", timeStampInp.Text, str)
	})
	return []fyne.CanvasObject{
		timeStampInp, click1, text3,
	}
}

func timeStampToDate(t int) (date string, err error) {
	nums, err := strconv.Atoi(strconv.Itoa(t))
	if err != nil {
		return
	}
	date = time.Unix(int64(nums/1000), 0).Format(timeTemplate1)
	return
}

//SecToDayHMS 秒转为天小时分钟秒钟
func secToDayHMS(seconds int) string {
	var day, hour, minute, second int
	day = seconds / (24 * 3600)
	hour = (seconds - day*3600*24) / 3600
	minute = (seconds - day*24*3600 - hour*3600) / 60
	second = seconds - day*24*3600 - hour*3600 - minute*60
	return strconv.Itoa(day) + "d " + strconv.Itoa(hour) + "h " + strconv.Itoa(minute) + "m " + strconv.Itoa(second) + "s"
}

//DayHMSToSec 天小时分钟秒钟转为秒
func dayHMSToSec(str string) string {
	byStr := strings.Split(str, " ")
	if len(byStr) != 4 {
		return ""
	}
	day, _ := strconv.Atoi(byStr[0])
	hour, _ := strconv.Atoi(byStr[1])
	minute, _ := strconv.Atoi(byStr[2])
	second, _ := strconv.Atoi(byStr[3])

	if day < 0 || hour < 0 || minute < 0 || second < 0 {
		return ""
	}
	return strconv.Itoa(day*3600*24 + hour*3600 + minute*60 + second)
}

func copyClipBoard(context string) {
	clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
	clipboard.SetContent(context)
	fmt.Println("success clipboard", context)
}
