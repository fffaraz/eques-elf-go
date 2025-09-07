package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Device struct {
	IP       string
	Mac      string
	Password string
	Status   string
}

func cmdDiscover() []Device {
	conn, err := net.ListenPacket("udp4", "0.0.0.0:0")
	if err != nil {
		log.Println("Error listening on UDP:", err)
		return nil
	}
	defer conn.Close()

	go func() {
		for i := 1; i < 2; i++ {
			command := "lan_phone%mac%nopassword%" + time.Now().Format("2006-01-02-15:04:05") + "%heart"
			ciphertext, _ := hex.DecodeString(aesEcb256Encrypt(command))
			conn.WriteTo(ciphertext, &net.UDPAddr{
				IP:   net.ParseIP("192.168.1.255"), // 255.255.255.255
				Port: 27431,
			})
			fmt.Printf("Sent discovery packet %d\n", i)
			time.Sleep(500 * time.Millisecond)
		}
		conn.Close()
	}()

	devices := make(map[string]Device)
	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			break
		}
		ipv4 := addr.(*net.UDPAddr).IP.To4().String()
		decrypted := aesEcb256Decrypt(hex.EncodeToString(buffer[:n]))
		fields := strings.Split(decrypted, "%")
		if len(fields) == 5 && fields[0] == "lan_device" {
			if _, ok := devices[ipv4]; ok {
				fmt.Printf("Device with IP %s already exists, skipping...\n", ipv4)
				continue
			}
			status := strings.Split(fields[3], "#")
			devices[ipv4] = Device{
				IP:       ipv4,
				Mac:      fields[1],
				Password: fields[2],
				Status:   status[0],
			}
			fmt.Printf("IP: %s | Mac: %s | Password: %s | Status: %s\n", ipv4, fields[1], fields[2], status[0])
			fmt.Printf("-cmd on -ip %s -mac %s -pass %s\n", ipv4, fields[1], fields[2])
			fmt.Println()
		} else {
			fmt.Printf("Received %d bytes from %s: %s\n", n, addr.String(), decrypted)
		}
	}
	if len(devices) == 0 {
		return nil
	}
	var result []Device
	for _, device := range devices {
		result = append(result, device)
	}
	return result
}

func sendCommand(device Device, command string) *Device {
	remoteaddress := net.UDPAddr{
		IP:   net.ParseIP(device.IP),
		Port: 27431,
	}

	conn, err := net.DialUDP("udp4", nil, &remoteaddress)
	if err != nil {
		log.Println("Error dialing UDP:", err)
		return nil
	}
	defer conn.Close()
	conn.SetWriteDeadline(time.Now().Add(2 * time.Second))

	msg, _ := hex.DecodeString(aesEcb256Encrypt("lan_phone%" + device.Mac + "%" + device.Password + "%" + command))
	_, err = conn.Write(msg)
	if err != nil {
		log.Println("Error sending message:", err)
		return nil
	}

	// Read response synchronously with timeout
	_ = conn.SetReadDeadline(time.Now().Add(1500 * time.Millisecond))
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFrom(buffer)
	if err != nil {
		log.Println("Error reading from UDP:", err)
		return &device
	}
	decrypted := aesEcb256Decrypt(hex.EncodeToString(buffer[:n]))
	fmt.Printf("Received %d bytes from %s: %s\n", n, addr.String(), decrypted)
	fields := strings.Split(decrypted, "%")
	if len(fields) == 5 && fields[0] == "lan_device" && fields[1] == device.Mac {
		device.Status = fields[3]
	}
	return &device
}

func sendCommandOn(device Device) *Device {
	return sendCommand(device, "open%relay")
}

func sendCommandOff(device Device) *Device {
	return sendCommand(device, "close%relay")
}

func main() {
	cmdPtr := flag.String("cmd", "", "Command to execute (discover, status, timer, on, off)")
	ipPtr := flag.String("ip", "", "IP address of the Eques elf Smart Plug")
	macPtr := flag.String("mac", "", "Mac address of the Eques elf Smart Plug in the format aa-bb-cc-dd-ee-ff")
	passPtr := flag.String("pass", "", "Password of the Eques elf Smart Plug")
	cliPtr := flag.Bool("cli", false, "Enable CLI mode (default: false)")
	flag.Parse()

	if !*cliPtr {
		runGUI()
		return
	}

	if *cmdPtr == "" || *cmdPtr == "discover" {
		cmdDiscover()
		return
	}
	if *ipPtr == "" || *macPtr == "" || *passPtr == "" {
		fmt.Println("Please provide the required parameters: ip, mac, and password.")
		return
	}

	device := Device{
		IP:       *ipPtr,
		Mac:      *macPtr,
		Password: *passPtr,
	}

	switch *cmdPtr {
	case "status":
		sendCommand(device, "check%relay")
	case "timer":
		sendCommand(device, "check#total%timer")
	case "on":
		sendCommandOn(device)
	case "off":
		sendCommandOff(device)
	}
}
