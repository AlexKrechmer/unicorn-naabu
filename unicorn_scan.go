package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
	// ==== FLAGS ====
	target := flag.String("target", "", "Target IP or hostname")
	ports := flag.String("p", "", "Ports to scan (e.g. 80,443)")
	fullTCP := flag.Bool("full-tcp", false, "Scan all TCP ports (1-65535)")
	useSudo := flag.Bool("sudo", false, "Use sudo for Naabu/Nmap")
	flag.Parse()

	if *target == "" && len(flag.Args()) > 0 {
		*target = flag.Args()[0]
	}
	if *target == "" {
		fmt.Println(Red + "[!] Please specify a target with -target or as the first argument" + Reset)
		os.Exit(1)
	}

	if *fullTCP {
		*ports = "1-65535"
	} else if *ports == "" {
		*ports = "80,443"
	}

// ==== PRINT BANNERS ====
func printBanners(target string) {
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
	fmt.Println(Purple + unicornArt + Reset)
	fmt.Println(Purple + naabuBanner + Reset)
	fmt.Printf("%s[*] Scanning target: %s%s\n\n", Purple, target, Reset)
}

// ==== SPINNER ====
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

	var naabuCmd *exec.Cmd
	if *useSudo {
		naabuCmd = exec.Command("sudo", append([]string{"naabu"}, naabuArgs...)...)
	} else {
		naabuCmd = exec.Command("naabu", naabuArgs...)
	}

	stdout, err := naabuCmd.StdoutPipe()
	if err != nil {
		fmt.Println(Red+"[!] Error capturing Naabu output: "+err.Error()+Reset, "")
		return
	}
	naabuCmd.Stderr = os.Stderr

	if err := naabuCmd.Start(); err != nil {
		fmt.Println(Red+"[!] Error starting Naabu: "+err.Error()+Reset, "")
		return
	}

	openPorts := []string{}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		port := ""
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			port = strings.TrimSpace(parts[len(parts)-1])
		} else {
			port = line
		}
		if _, err := strconv.Atoi(port); err == nil {
			openPorts = append(openPorts, port)
		}
	}
	naabuCmd.Wait()
	spinnerDone <- true
	fmt.Println("\rScan complete.                 ")

	if len(openPorts) == 0 {
		fmt.Printf("%s[!] No open ports found by Naabu, skipping Nmap.%s\n", Red, Reset)
		return
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

	var nmapCmd *exec.Cmd
	if *useSudo {
		nmapCmd = exec.Command("sudo", append([]string{"nmap"}, nmapArgs...)...)
	} else {
		nmapCmd = exec.Command("nmap", nmapArgs...)
	}

	nmapOut, err := nmapCmd.CombinedOutput()
	if err != nil {
		fmt.Println(Red+"[!] Error running Nmap: "+err.Error()+Reset, "")
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
