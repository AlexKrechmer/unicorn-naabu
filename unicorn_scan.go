package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
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
	fullTCP := flag.Bool("full-tcp", true, "Scan all TCP ports (1-65535) by default")
	fullScan := flag.Bool("full", false, "Run full Naabu + Nmap + OS detection scan")
	useSudo := flag.Bool("sudo", true, "Use sudo for SYN scans by default")
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

	if *fullScan {
		runFullScan(*target, *fullTCP, *useSudo, *minRate, *timing, *saveFile)
		return
	}

	// Standard behavior with retries/backoff
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
		sortPorts(openPorts)
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

// ==== FULL SCAN FUNCTION ====
func runFullScan(target string, fullTCP, useSudo bool, minRate, timing int, saveFile string) {
	fmt.Println(Cyan + "[*] Running FULL Naabu + Nmap + OS detection..." + Reset)

	allPorts := []string{}

	// Pass 1: Full TCP scan
	fmt.Println(Yellow + "[*] Naabu pass 1: Full TCP scan (-p-)" + Reset)
	ports1 := runNaabuLive(target, fullTCP, useSudo, minRate, saveFile)
	allPorts = mergePorts(allPorts, ports1)

	// Pass 2: SYN scan (requires sudo)
	if useSudo {
		fmt.Println(Yellow + "[*] Naabu pass 2: SYN scan" + Reset)
		ports2 := runNaabuLive(target, fullTCP, useSudo, minRate, saveFile)
		allPorts = mergePorts(allPorts, ports2)
	}

	// Pass 3: TCP connect fallback (non-root)
	fmt.Println(Yellow + "[*] Naabu pass 3: TCP connect fallback" + Reset)
	ports3 := runNaabuLive(target, fullTCP, false, minRate, saveFile)
	allPorts = mergePorts(allPorts, ports3)

	if len(allPorts) > 0 {
		sortPorts(allPorts)
		portList := strings.Join(allPorts, ",")
		fmt.Printf("%s[Naabu] Open ports after all passes: %s%s\n\n", Purple, portList, Reset)
		runNmapFullScan(target, portList, useSudo, timing, saveFile)
	} else {
		fmt.Printf("%s[!] No ports found. Running Nmap top 1000 ports fallback.%s\n", Red, Reset)
		runNmapCommon(target, useSudo, timing, saveFile)
	}

	fmt.Println(Green + "[+] Full scan summary complete." + Reset)
}

// ==== NAABU LIVE ====
func runNaabuLive(target string, fullTCP, useSudo bool, minRate int, saveFile string) []string {
    args := []string{"-host", target, "-oJ", "-", "--min-rate", strconv.Itoa(minRate)}
    if fullTCP {
        args = append(args, "-p-")
    }

    var cmd *exec.Cmd
    if useSudo {
        cmd = exec.Command("sudo", append([]string{"naabu"}, args...)...)
    } else {
        cmd = exec.Command("naabu", args...)
    }

    fmt.Println(Cyan+"[*] Starting full Naabu sweep..."+Reset)
    stdout, _ := cmd.StdoutPipe()
    cmd.Stderr = cmd.Stdout

    if err := cmd.Start(); err != nil {
        fmt.Println(Red+"[!] Failed to start Naabu:", err, Reset)
        return nil
    }

    scanner := bufio.NewScanner(stdout)
    openPorts := []string{}
    portSet := make(map[string]bool)

    for scanner.Scan() {
        line := scanner.Text()
        if line == "" { continue }
        // Parse JSON output from -oJ
        if strings.Contains(line, `"port"`) {
            re := regexp.MustCompile(`"port":\s*(\d+)`)
            if m := re.FindStringSubmatch(line); m != nil {
                port := m[1]
                if !portSet[port] {
                    portSet[port] = true
                    openPorts = append(openPorts, port)
                    fmt.Printf(Green+"[+] Open port found: %s%s\n", port, Reset)
                }
            }
        }
    }

    cmd.Wait()
    return openPorts
}
// ==== SORT PORTS ====
func sortPorts(ports []string) {
	sort.Slice(ports, func(i, j int) bool {
		a, _ := strconv.Atoi(ports[i])
		b, _ := strconv.Atoi(ports[j])
		return a < b
	})
}

// ==== NMAP FUNCTIONS ====
func runNmapColor(target, ports string, useSudo bool, timing int, saveFile string) {
	args := []string{"-sS", "-sV", "-Pn", "-p", ports, "-T" + strconv.Itoa(timing), target}
	fmt.Println(Cyan + "[*] Running Nmap on discovered ports..." + Reset)
	runCmdLiveSave("nmap", args, useSudo, saveFile)
}

func runNmapFullScan(target, ports string, useSudo bool, timing int, saveFile string) {
	args := []string{"-sS", "-sV", "-O", "-Pn", "-p", ports, "-T" + strconv.Itoa(timing), target}
	fmt.Println(Cyan + "[*] Running full Nmap scan with OS detection..." + Reset)
	runCmdLiveSave("nmap", args, useSudo, saveFile)
}

func runNmapFull(target string, useSudo bool, timing int, saveFile string) {
	args := []string{"-sS", "-sV", "-Pn", "-p-", "-T" + strconv.Itoa(timing), target}
	fmt.Println(Cyan + "[*] Running full Nmap scan..." + Reset)
	runCmdLiveSave("nmap", args, useSudo, saveFile)
}

func runNmapCommon(target string, useSudo bool, timing int, saveFile string) {
	args := []string{"-sS", "-sV", "-Pn", "-T" + strconv.Itoa(timing), target}
	fmt.Println(Cyan + "[*] Running fast common ports Nmap scan..." + Reset)
	runCmdLiveSave("nmap", args, useSudo, saveFile)
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
			file.WriteString("[" + strings.Title(name) + "] " + line + "\n")
		}
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(Red + "[!] Command execution error:", err, Reset)
	}
}
