package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"sync/atomic"
	"time"
)

var (
	state int64
	delay int64 = 100
	exit  int64
	delta int64 = 10
)

var (
	delayItem *systray.MenuItem
	stateItem *systray.MenuItem
	itemQuit  *systray.MenuItem
)

func main() {
	go click()
	go add()
	systray.Run(onReady, onExit)
}

func add() {
	fmt.Println("--- Please press ctrl + shift + q to stop hook ---")
	robotgo.EventHook(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		robotgo.EventEnd()
		atomic.StoreInt64(&exit, 1)
		systray.Quit()
	})
	fmt.Println("--- Please press ctrl + shift + s to start clicking ---")
	robotgo.EventHook(hook.KeyDown, []string{"s", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("start")
		atomic.StoreInt64(&state, 1)
		updateStateItem(true)
	})
	fmt.Println("--- Please press ctrl + shift + p to pause clicking ---")
	robotgo.EventHook(hook.KeyDown, []string{"p", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("pause")
		atomic.StoreInt64(&state, 0)
		updateStateItem(false)
	})
	fmt.Println("--- Please press ctrl + shift + + to increase click delay by", delta, "milliseconds ---")
	robotgo.EventHook(hook.KeyDown, []string{"+", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("inc", atomic.AddInt64(&delay, delta))
		updateDelayMenuItem()
	})
	fmt.Println("--- Please press ctrl + shift + + to decrease click delay by", delta, "milliseconds ---")
	robotgo.EventHook(hook.KeyDown, []string{"-", "ctrl", "shift"}, func(e hook.Event) {
		if atomic.LoadInt64(&delay) > delta {
			fmt.Println("decr", atomic.AddInt64(&delay, -delta))
		} else {
			fmt.Println("can't set delay below 0")
		}
		updateDelayMenuItem()
	})

	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func click() {
	for {
		if atomic.LoadInt64(&state) != 0 {
			robotgo.Click("left", false)
			fmt.Println("click")
		}
		if atomic.LoadInt64(&exit) != 0 {
			return
		}
		clickDelay := atomic.LoadInt64(&delay)
		time.Sleep(time.Microsecond * time.Duration(clickDelay))
	}
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Clicker")
	systray.SetTooltip("Clicker")

	stateItem = systray.AddMenuItemCheckbox("state", "Current state", false)
	stateItem.Disable()
	delayItem = systray.AddMenuItem(delayTitle(), "Click delay")
	delayItem.Disable()
	systray.AddSeparator()
	itemQuit = systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	itemQuit.SetIcon(icon.Data)
	go func() {
		<-itemQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	atomic.StoreInt64(&exit, 1)
	clickDelay := atomic.LoadInt64(&delay)
	time.Sleep(2 * time.Millisecond * time.Duration(clickDelay))
}

func delayTitle() string {
	return fmt.Sprintf("delay: %dms", atomic.LoadInt64(&delay))
}

func updateDelayMenuItem() {
	delayItem.SetTitle(delayTitle())
	fmt.Println("update")
}

func updateStateItem(active bool) {
	if active {
		stateItem.Check()
		return
	}
	stateItem.Uncheck()
}
