package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func runGUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Eques Elf Smart Plug Remote Control")

	status := widget.NewLabel("Idle")
	deviceList := container.NewVBox()
	scroller := container.NewVScroll(deviceList)
	scroller.SetMinSize(fyne.NewSize(500, 300))

	var discoverBtn *widget.Button
	discoverBtn = widget.NewButton("Discover", func() {
		discoverBtn.Disable()
		status.SetText("Discovering devices...")
		go func() {
			devices := cmdDiscover()
			deviceList.Objects = nil
			if len(devices) == 0 {
				status.SetText("No devices found")
				discoverBtn.Enable()
				deviceList.Refresh()
				return
			}
			for _, d := range devices {
				dev := d // capture
				statusLbl := widget.NewLabel(dev.Status)
				var onBtn, offBtn *widget.Button
				onBtn = widget.NewButton("On", func() {
					onBtn.Disable()
					offBtn.Disable()
					go func() {
						resp := sendCommandOn(dev)
						fyne.DoAndWait(func() {
							if resp != nil {
								statusLbl.SetText(resp.Status)
							}
							onBtn.Enable()
							offBtn.Enable()
						})
					}()
				})
				offBtn = widget.NewButton("Off", func() {
					offBtn.Disable()
					onBtn.Disable()
					go func() {
						resp := sendCommandOff(dev)
						fyne.DoAndWait(func() {
							if resp != nil {
								statusLbl.SetText(resp.Status)
							}
							onBtn.Enable()
							offBtn.Enable()
						})
					}()
				})
				row := container.NewHBox(
					widget.NewLabel(dev.IP),
					widget.NewLabel(dev.Mac),
					statusLbl,
					onBtn,
					offBtn,
				)
				deviceList.Add(row)
			}
			fyne.DoAndWait(func() {
				deviceList.Refresh()
				status.SetText(fmt.Sprintf("Found %d device(s)", len(devices)))
				discoverBtn.Enable()
			})
		}()
	})

	myWindow.SetContent(
		container.NewBorder(
			container.NewHBox(discoverBtn),
			status,
			nil,
			nil,
			scroller,
		),
	)
	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()
}
