package main

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const timeTemplate1 = "2006-01-02 15:04:05"

func main() {
	myApp := app.New()
	setting := myApp.Settings()
	setting.SetTheme(theme.LightTheme())
	resource, err := fyne.LoadResourceFromPath("img/time_icon.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	myApp.SetIcon(resource)
	myWindow := myApp.NewWindow("Time")

	gridElems := TimestampToDate()
	gridElems2 := DateToTimestamp()
	gridElems = append(gridElems, gridElems2...)
	grid := container.New(layout.NewGridLayout(3), gridElems...)
	myWindow.SetContent(grid)
	myWindow.Resize(fyne.NewSize(500, 50))
	myWindow.ShowAndRun()
}

func TimestampToDate() []fyne.CanvasObject {
	timeStampInp := widget.NewEntry()
	now := time.Now().UnixNano() / 1e6
	timeStampInp.SetText(strconv.Itoa(int(now)))
	timeStampInp.SetPlaceHolder("time stamp")
	text3 := widget.NewLabel("date")
	click1 := widget.NewButton("click", func() {
		nums, err := strconv.Atoi(timeStampInp.Text)
		if err != nil {
			fmt.Println(err)
			return
		}
		date := time.Unix(int64(nums/1000), 0).Format(timeTemplate1)
		text3.SetText(date)
		//复制到剪切版
		clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
		clipboard.SetContent(date)
	})
	return []fyne.CanvasObject{
		timeStampInp, click1, text3,
	}
}

func DateToTimestamp() []fyne.CanvasObject {
	timeStampInp := widget.NewEntry()
	timeStampInp.SetPlaceHolder("date")
	text3 := widget.NewLabel("time stamp")
	click1 := widget.NewButton("click", func() {
		stamp, err := time.ParseInLocation(timeTemplate1, timeStampInp.Text, time.Local)
		if err != nil {
			fmt.Println(err)
			return
		}
		millionSec := strconv.Itoa(int(stamp.UnixNano() / 1e6))
		text3.SetText(millionSec)
		clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
		clipboard.SetContent(millionSec)
	})
	return []fyne.CanvasObject{
		timeStampInp, click1, text3,
	}
}
