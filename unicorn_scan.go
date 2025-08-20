package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ANSI color codes
const (
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Gray   = "\033[37m"
	Reset  = "\033[0m"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./unicorn_scan <target> [ports] [-fast|-full]")
		return
	}

	target := os.Args[1]
	portArg := ""
	scanMode := "-fast"

	// Parse optional port and mode
	for _, arg := range os.Args[2:] {
		if strings.HasPrefix(arg, "-") {
			scanMode = arg
		} else {
			portArg = arg
		}
	}

	// Set default ports based on scan mode
	if portArg == "" {
		if scanMode == "-full" {
			portArg = "-p-" // scan all TCP ports
		} else {
			portArg = "80,443" // default fast scan ports
		}
	}

	// Unicorn ASCII art
	unicornArt := `
⠀⠀⠀⠀⠀⠑⢦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠙⢷⣦⣀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⣷⣿⣾⣿⣧⣄⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⣿⣿⣇⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⣿⣿⣿⣿⣿⣿⣿⣿⣥⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⠟⠉⠉⢹⣿⣿⣿⣿⣿⣿⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣿⣿⣿⣿⣿⣿⣿⣿⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣠⣶⣶⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡄⠀⠀⠀⠀⠀⠀⠀⠀
⢀⣴⠿⠛⠉⢸⡏⠁⠉⠙⠛⠻⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣄⡀⠀⠀⠀⠀⠀
⠉⠉⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦⡀⠀⠀⠀
⠀⠀⠀⠀⠀⠈⠿⠀⠀⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠛⠻⢿⣿⣿⣿⣿⣿⣿⣧⡀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⣿⠟⢿⣷⡄
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⡟⠀⢠⣾⣿⣿
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣿⣀⣾⣿⡿⠃
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⠏⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣿⣿⠻⣿⣿⡀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠐⠋⣹⣿⠃⠀⠈⣿⣿⣴⠇
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⣾⠟⠀⠀⠀⠀⠘⠉⠛⠀
`

	// Naabu banner ASCII
	naabuBanner := `
  ___  ___  ___ _/ /  __ __
 / _ \/ _ \/ _ \/ _ \/ // /
/_//_/\_,_/\_,_/_.__/\_,_/

                projectdiscovery.io
`
	// Print unicorn and banner
	fmt.Println(Purple + unicornArt + Reset)
	fmt.Println(Purple + naabuBanner + Reset)
	fmt.Printf("%s[*] Scanning target: %s%s\n\n", Purple, target, Reset)


	// Spinner goroutine
	spinnerDone := make(chan bool)
	go func() {
		spin := []rune{'|', '/', '-', '\\'}
		i := 0
		for {
			select {
			case <-spinnerDone:
				return
			default:
				fmt.Printf("\r%sScanning %s... %c%s", Yellow, target, spin[i%len(spin)], Reset)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Build Naabu command
	naabuArgs := []string{"-silent", "-host", target, "-no-ping"}
	if scanMode == "-full" {
		naabuArgs = append(naabuArgs, "-p-")
	} else if portArg != "" && portArg != "-p-" {
		naabuArgs = append(naabuArgs, "-ports", portArg)
	}

	naabuCmd := exec.Command("naabu", naabuArgs...)
	stdout, err := naabuCmd.StdoutPipe()
	if err != nil {
		fmt.Println(Red+"[!] Error capturing Naabu output:"+err.Error()+Reset, "")
		return
	}
	naabuCmd.Stderr = os.Stderr
	if err := naabuCmd.Start(); err != nil {
		fmt.Println(Red+"[!] Error starting Naabu:"+err.Error()+Reset, "")
		return
	}

	ports := []string{}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				ports = append(ports, parts[1])
			} else {
				ports = append(ports, parts[0])
			}
		}
	}
	naabuCmd.Wait()
	spinnerDone <- true
	fmt.Println("\rScan complete.                 ")

	if len(ports) == 0 {
		fmt.Printf("%s[!] No open ports found by Naabu, using default ports %s%s\n", Red, portArg, Reset)
		ports = []string{portArg}
	}

	portList := strings.Join(ports, ",")
	fmt.Printf("%s[Naabu] Open ports: %s%s\n\n", Purple, portList, Reset)

	// Build Nmap command
	var nmapArgs []string
	if portList == "-p-" {
		nmapArgs = []string{"-sC", "-sV", "-Pn", "-p-", target}
	} else {
		nmapArgs = []string{"-sC", "-sV", "-Pn", "-p", portList, target}
	}

	fmt.Println(Cyan + "[*] Running Nmap..." + Reset)
	nmapCmd := exec.Command("nmap", nmapArgs...)
	nmapOut, err := nmapCmd.CombinedOutput()
	if err != nil {
		fmt.Println(Red+"[!] Error running Nmap:"+err.Error()+Reset, "")
		return
	}

	// Colorize Nmap output
	for _, line := range strings.Split(string(nmapOut), "\n") {
		switch {
		case strings.Contains(line, "open"):
			fmt.Println(Green + line + Reset)
		case strings.Contains(line, "filtered"):
			fmt.Println(Yellow + line + Reset)
		case strings.Contains(line, "closed"):
			fmt.Println(Gray + line + Reset)
		default:
			fmt.Println(Cyan + line + Reset)
		}
	}

	fmt.Println(Green + "[+] Scan summary complete." + Reset)
}
