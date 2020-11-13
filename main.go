package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/text/encoding/charmap"
)

func execCmd(path string, args ...string) {
	cmd := exec.Command(path, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	err = cmd.Wait()

	utf8 := buf.String()

	decoder := charmap.CodePage866.NewDecoder()
	reader, _ := decoder.String(utf8)

	fmt.Printf("Command finished with error: %v\n", err)
	fmt.Printf("Command finished with output: %v\n", reader)
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("usage: wifiaddcmd.exe ssid passwd")
		return
	}

	ssid := os.Args[1]
	passwd := os.Args[2]

	text := `<?xml version="1.0" encoding="UTF-8" ?>
	<WLANProfile xmlns="http://www.microsoft.com/networking/WLAN/profile/v1">
	    <name>` + ssid + `</name>
	    <SSIDConfig>
	        <SSID>
	            <name>` + ssid + `</name>
	        </SSID>
	    </SSIDConfig>
	    <connectionType>ESS</connectionType>
	    <connectionMode>auto</connectionMode>
	    <MSM>
	        <security>
	            <authEncryption>
	                <authentication>WPA2PSK</authentication>
	                <encryption>AES</encryption>
	                <useOneX>false</useOneX>
	            </authEncryption>
	            <sharedKey>
	                <keyType>passPhrase</keyType>
	                <protected>false</protected>
	                <keyMaterial>` + passwd + `</keyMaterial>
	            </sharedKey>
	        </security>
	    </MSM>
	</WLANProfile>`

	path, _ := os.LookupEnv("temp")

	wifixml := path + "\\wifixml.xml"

	file, err := os.Create(wifixml)
	if err != nil {
		fmt.Println(err)
	} else {
		file.WriteString(text)
	}
	file.Close()

	execCmd("cmd", "/c netsh wlan add profile filename="+wifixml)

	os.Remove(wifixml)

	return
}
