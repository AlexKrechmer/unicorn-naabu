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
)

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
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣿⣿⠻⣿⣿⡀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠐⠋⣹⣿⠃⠀⠈⣿⣿⣴⠇
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⣾⠟⠀⠀⠀⠀⠘⠉⠛⠀
`
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

func main() {
	// ==== FLAGS ====
	target := flag.String("target", "", "Target IP or hostname")
	fullTCP := flag.Bool("full-tcp", false, "Scan all TCP ports (1-65535)")
	useSudo := flag.Bool("sudo", false, "Use sudo for Naabu/Nmap")
	timing := flag.Int("T", 5, "Nmap timing template (0-5), default 5=Insane")
	flag.Parse()

	if *target == "" && len(flag.Args()) > 0 {
		*target = flag.Args()[0]
	}
	if *target == "" {
		fmt.Println(Red + "[!] Please specify a target with -target or as the first argument" + Reset)
		os.Exit(1)
	}

	printBanners(*target)

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

	if len(openPorts) > 0 {
		portList := strings.Join(openPorts, ",")
		fmt.Printf("%s[Naabu] Open ports: %s%s\n\n", Purple, portList, Reset)
		runNmap(*target, portList, *fullTCP, *useSudo, *timing)
	} else {
		fmt.Printf("%s[!] No open ports found, running full Nmap scan until cancelled.%s\n", Red, Reset)
		runNmapFull(*target, *useSudo, *timing)
	}

	fmt.Println(Green + "[+] Scan summary complete." + Reset)
}

// ==== RUN NMAP ====
func runNmap(target, ports string, fullTCP, useSudo bool, timing int) {
	var nmapArgs []string
	if fullTCP {
		nmapArgs = []string{"-sC", "-sV", "-Pn", "-p-", "-T" + strconv.Itoa(timing), target}
	} else {
		nmapArgs = []string{"-sC", "-sV", "-Pn", "-p", ports, "-T" + strconv.Itoa(timing), target}
	}
	fmt.Println(Cyan + "[*] Running Nmap..." + Reset)
	runCmdLive("nmap", nmapArgs, useSudo)
}

func runNmapFull(target string, useSudo bool, timing int) {
	args := []string{"-sC", "-sV", "-Pn", "-p-", "-T" + strconv.Itoa(timing), target}
	runCmdLive("nmap", args, useSudo)
}

// ==== COMMAND EXECUTION ====
func runCmdLive(name string, args []string, useSudo bool) {
	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.Command("sudo", append([]string{name}, args...)...)
	} else {
		cmd = exec.Command(name, args...)
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		fmt.Println(Red+"[!] Failed to start command:", err, Reset)
		return
	}

	scannerOut := bufio.NewScanner(stdout)
	scannerErr := bufio.NewScanner(stderr)

	go func() {
		for scannerOut.Scan() {
			fmt.Println(Green + scannerOut.Text() + Reset)
		}
	}()
	go func() {
		for scannerErr.Scan() {
			fmt.Println(Red + scannerErr.Text() + Reset)
		}
	}()

	cmd.Wait()
}

