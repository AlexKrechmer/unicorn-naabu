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
	useSudo := flag.Bool("sudo", false, "Use sudo for SYN scans")
	timing := flag.Int("T", 5, "Nmap timing template (0-5)")
	minRate := flag.Int("min-rate", 5000, "Naabu minimum rate for speed")
	retries := flag.Int("retries", 3, "Number of retries if no ports found")
	backoff := flag.Int("backoff", 3, "Backoff seconds between retries")
	saveFile := flag.String("save", "", "Optional file to save Naabu/Nmap outputs")
	flag.Parse()

	if *target == "" && len(flag.Args()) > 0 {
		*target = flag.Args()[0]
	}
	if *target == "" {
		fmt.Println(Red + "[!] Please specify a target with -target or as the first argument" + Reset)
		os.Exit(1)
	}

	if *useSudo && os.Geteuid() != 0 {
		fmt.Println(Red + "[!] SYN scans require sudo/root. Please run this script with sudo." + Reset)
		os.Exit(1)
	}

	printBanners(*target)

	openPorts := []string{}
	for attempt := 1; attempt <= *retries; attempt++ {
		fmt.Printf(Cyan+"[*] Naabu scan attempt %d/%d...%s\n", attempt, *retries, Reset)
		newPorts := runNaabuLive(*target, *fullTCP, *useSudo, *minRate, *saveFile)
		openPorts = mergePorts(openPorts, newPorts)
		if len(openPorts) > 0 {
			break
		}
		if attempt < *retries {
			fmt.Printf(Yellow+"[!] No ports found. Backing off for %d seconds before retry...%s\n", *backoff, Reset)
			time.Sleep(time.Duration(*backoff) * time.Second)
		}
	}

	if len(openPorts) > 0 {
		portList := strings.Join(openPorts, ",")
		fmt.Printf("%s[Naabu] Open ports found: %s%s\n\n", Purple, portList, Reset)
		runNmapColor(*target, portList, *useSudo, *timing, *saveFile)
	} else {
		fmt.Printf("%s[!] No open ports found, running fast common ports Nmap scan first.%s\n", Red, Reset)
		runNmapCommon(*target, *useSudo, *timing, *saveFile)
		fmt.Printf("%s[*] Escalating to full TCP Nmap scan...%s\n", Yellow, Reset)
		runNmapFull(*target, *useSudo, *timing, *saveFile)
	}

	fmt.Println(Green + "[+] Scan summary complete." + Reset)
}

// ==== NAABU LIVE ====
func runNaabuLive(target string, fullTCP, useSudo bool, minRate int, saveFile string) []string {
	args := []string{"-silent", "-host", target, "--min-rate", strconv.Itoa(minRate)}
	if fullTCP {
		args = append(args, "-p-")
	}

	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.Command("sudo", append([]string{"naabu"}, args...)...)
	} else {
		cmd = exec.Command("naabu", args...)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(Red+"[!] Failed to get stdout:", err, Reset)
		return nil
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		fmt.Println(Red+"[!] Failed to start Naabu:", err, Reset)
		return nil
	}

	scanner := bufio.NewScanner(stdout)
	openPorts := []string{}
	portSet := make(map[string]bool)

	var file *os.File
	if saveFile != "" {
		file, err = os.OpenFile(saveFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(Red+"[!] Could not open save file:", err, Reset)
		}
		defer file.Close()
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if file != nil {
			file.WriteString("[Naabu] " + line + "\n")
		}

		port := line
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			port = strings.Split(parts[1], " ")[0]
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

// ==== MERGE PORTS ====
func mergePorts(a, b []string) []string {
	portSet := make(map[string]bool)
	for _, p := range a {
		portSet[p] = true
	}
	for _, p := range b {
		portSet[p] = true
	}
	result := []string{}
	for p := range portSet {
		result = append(result, p)
	}
	return result
}

// ==== NMAP COLOR OUTPUT ====
func runNmapColor(target, ports string, useSudo bool, timing int, saveFile string) {
	args := []string{"-sC", "-sV", "-Pn", "-p", ports, "-T" + strconv.Itoa(timing), target}
	fmt.Println(Cyan + "[*] Running Nmap on discovered ports..." + Reset)

	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.Command("sudo", append([]string{"nmap"}, args...)...)
	} else {
		cmd = exec.Command("nmap", args...)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(Red+"[!] Failed to get stdout:", err, Reset)
		return
	}
	cmd.Stderr = cmd.Stdout

	var file *os.File
	if saveFile != "" {
		file, err = os.OpenFile(saveFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(Red+"[!] Could not open save file:", err, Reset)
		}
		defer file.Close()
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(Red+"[!] Failed to start Nmap:", err, Reset)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		colored := line
		if strings.Contains(line, "/tcp") || strings.Contains(line, "/udp") {
			switch {
			case strings.Contains(line, "open"):
				colored = Green + line + Reset
			case strings.Contains(line, "closed"):
				colored = Red + line + Reset
			case strings.Contains(line, "filtered"):
				colored = Yellow + line + Reset
			}
		}
		fmt.Println(colored)
		if file != nil {
			file.WriteString("[Nmap] " + line + "\n")
		}
	}
	cmd.Wait()
}

// ==== COMMON & FULL NMAP ====
func runNmapCommon(target string, useSudo bool, timing int, saveFile string) {
	args := []string{"-sC", "-sV", "-Pn", "-T" + strconv.Itoa(timing), target}
	fmt.Println(Cyan + "[*] Running fast common ports Nmap scan..." + Reset)
	runCmdLiveSave("nmap", args, useSudo, saveFile)
}

func runNmapFull(target string, useSudo bool, timing int, saveFile string) {
	args := []string{"-sC", "-sV", "-Pn", "-p-", "-T" + strconv.Itoa(timing), target}
	fmt.Println(Cyan + "[*] Running full Nmap scan..." + Reset)
	runCmdLiveSave("nmap", args, useSudo, saveFile)
}

// ==== COMMAND EXECUTION WITH OPTIONAL SAVE ====
func runCmdLiveSave(name string, args []string, useSudo bool, saveFile string) {
	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.Command("sudo", append([]string{name}, args...)...)
	} else {
		cmd = exec.Command(name, args...)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(Red + "[!] Failed to get stdout:", err, Reset)
		return
	}
	cmd.Stderr = cmd.Stdout

	var file *os.File
	if saveFile != "" {
		file, err = os.OpenFile(saveFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(Red + "[!] Could not open save file:", err, Reset)
		}
		defer file.Close()
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(Red + "[!] Failed to start command:", err, Reset)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if file != nil {
			file.WriteString("[" + strings.Title(name) + "] " + line + "\n")
		}
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(Red + "[!] Command execution error:", err, Reset)
	}
}
