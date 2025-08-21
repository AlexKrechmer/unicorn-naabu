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
⠀⠀⠀⢀⣠⣶⣶⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡄⠀⠀⠀⠀⠀⠀⠀⠀
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

// ==== MAIN ====
func main() {
	target := flag.String("target", "", "Target IP or hostname")
	fullTCP := flag.Bool("full-tcp", false, "Scan all TCP ports (1-65535)")
	useSudo := flag.Bool("sudo", false, "Use sudo for Naabu/Nmap")
	timing := flag.Int("T", 5, "Nmap timing template (0-5), default 5=Insane")
	minRate := flag.Int("min-rate", 5000, "Naabu minimum rate for speed")
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

	// ==== RUN NAABU LIVE ====
	openPorts := runNaabuLive(*target, *fullTCP, *useSudo, *minRate)
	spinnerDone <- true
	fmt.Println("\rScan complete.                 ")

	if len(openPorts) > 0 {
		portList := strings.Join(openPorts, ",")
		fmt.Printf("%s[Naabu] Open ports found: %s%s\n\n", Purple, portList, Reset)
		runNmap(*target, portList, *useSudo, *timing)
	} else {
		fmt.Printf("%s[!] No open ports found, running full Nmap scan until cancelled.%s\n", Red, Reset)
		runNmapFull(*target, *useSudo, *timing)
	}

	fmt.Println(Green + "[+] Scan summary complete." + Reset)
}

// ==== NAABU LIVE ====
func runNaabuLive(target string, fullTCP, useSudo bool, minRate int) []string {
	args := []string{"-silent", "-host", target, "-no-ping", "--min-rate", strconv.Itoa(minRate)}
	if fullTCP {
		args = append(args, "-p-")
	}

	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.Command("sudo", append([]string{"naabu"}, args...)...)
	} else {
		cmd = exec.Command("naabu", args...)
	}

	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println(Red+"[!] Failed to start Naabu:", err, Reset)
		return nil
	}

	scanner := bufio.NewScanner(stdout)
	openPorts := []string{}
	portSet := make(map[string]bool)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		port := ""
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			port = strings.Split(parts[1], " ")[0]
		} else {
			port = line
		}
		if _, err := strconv.Atoi(port); err == nil {
			if !portSet[port] {
				portSet[port] = true
				openPorts = append(openPorts, port)
				fmt.Printf(Green+"[+] New open port found: %s%s\n", port, Reset)
			}
		}
	}

	cmd.Wait()
	return openPorts
}

// ==== NMAP ====
func runNmap(target, ports string, useSudo bool, timing int) {
	args := []string{"-sC", "-sV", "-Pn", "-p", ports, "-T" + strconv.Itoa(timing), target}
	fmt.Println(Cyan + "[*] Running Nmap on discovered ports..." + Reset)
	runCmdLive("nmap", args, useSudo)
}

func runNmapFull(target string, useSudo bool, timing int) {
	args := []string{"-sC", "-sV", "-Pn", "-p-", "-T" + strconv.Itoa(timing), target}
	fmt.Println(Cyan + "[*] Running full Nmap scan..." + Reset)
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
