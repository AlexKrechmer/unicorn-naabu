package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Purple = "\033[35m"
	Gray   = "\033[37m"
)

func main() {
	// ==== FLAG HANDLING =====
	target := flag.String("target", "", "Target IP or hostname")
	ports := flag.String("p", "", "Ports to scan (e.g. 80,443 or 1-65535)")
	fullTCP := flag.Bool("full-tcp", false, "Scan all TCP ports (1-65535)")
	flag.Parse()

	// Support positional target
	if *target == "" && len(flag.Args()) > 0 {
		*target = flag.Args()[0]
	}

	if *target == "" {
		fmt.Println(Red + "Please specify a target with -target or as the first argument" + Reset)
		os.Exit(1)
	}

	// Full TCP overrides ports
	if *fullTCP {
		*ports = "1-65535"
	} else if *ports == "" {
		*ports = "80,443"
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



	
	// ==== Spinner ====
	spinnerDone := make(chan bool)
	go func() {
		spin := []rune{'|', '/', '-', '\\'}
		i := 0
		for {
			select {
			case <-spinnerDone:
				return
			default:
				fmt.Printf("\r%sScanning %s... %c%s", Yellow, *target, spin[i%len(spin)], Reset)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// ==== RUN NAABU ====
	naabuArgs := []string{"-silent", "-host", *target, "-no-ping"}
	if *fullTCP {
		naabuArgs = append(naabuArgs, "-p-")
	} else {
		naabuArgs = append(naabuArgs, "-ports", *ports)
	}

	naabuCmd := exec.Command("sudo", "naabu", naabuArgs...)
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

	openPorts := []string{}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				openPorts = append(openPorts, parts[1])
			} else {
				openPorts = append(openPorts, parts[0])
			}
		}
	}
	naabuCmd.Wait()
	spinnerDone <- true
	fmt.Println("\rScan complete.                 ")

	if len(openPorts) == 0 {
		fmt.Printf("%s[!] No open ports found by Naabu, using default ports %s%s\n", Red, *ports, Reset)
		openPorts = []string{*ports}
	}

	portList := strings.Join(openPorts, ",")
	fmt.Printf("%s[Naabu] Open ports: %s%s\n\n", Purple, portList, Reset)

	// ==== RUN NMAP ====
	var nmapArgs []string
	if *fullTCP {
		nmapArgs = []string{"-sC", "-sV", "-Pn", "-p-", *target}
	} else {
		nmapArgs = []string{"-sC", "-sV", "-Pn", "-p", portList, *target}
	}

	fmt.Println(Cyan + "[*] Running Nmap..." + Reset)
	nmapCmd := exec.Command("sudo", "nmap", nmapArgs...)
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
