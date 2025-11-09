package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	Display  DisplayConfig  `toml:"display"`
	Colors   ColorsConfig   `toml:"colors"`
	Disks    DisksConfig    `toml:"disks"`
	Services ServicesConfig `toml:"services"`
}

type ColorsConfig struct {
	Header   string `toml:"header"`
	Title    string `toml:"title"`
	Label    string `toml:"label"`
	User     string `toml:"user"`
	FQDN     string `toml:"fqdn"`
	Dim      string `toml:"dim"`
	Active   string `toml:"active"`
	Inactive string `toml:"inactive"`
	Failed   string `toml:"failed"`
	Green    string `toml:"green"`
	Orange   string `toml:"orange"`
	Red      string `toml:"red"`
}

type DisplayConfig struct {
	LabelWidth  int `toml:"label_width"`
	GreenUntil  int `toml:"green_until"`
	OrangeUntil int `toml:"orange_until"`
}

type DisksConfig struct {
	Paths []string `toml:"paths"`
}

type ServicesConfig struct {
	Monitored []string `toml:"monitored"`
}

var (
	headerStyle   lipgloss.Style
	titleStyle    lipgloss.Style
	labelStyle    lipgloss.Style
	userStyle     lipgloss.Style
	fqdnStyle     lipgloss.Style
	dimStyle      lipgloss.Style
	activeStyle   lipgloss.Style
	inactiveStyle lipgloss.Style
	failedStyle   lipgloss.Style
	greenStyle    lipgloss.Style
	orangeStyle   lipgloss.Style
	redStyle      lipgloss.Style

	// Default vars
	labelWidth  = 15
	GreenUntil  = 65
	OrangeUntil = 85
)

// Default color scheme (Catppuccin Macchiato)
var defaultColors = ColorsConfig{
	Header:   "#ea76cb", // latte pink
	Title:    "#8bd5ca", // macchiatto teal
	Label:    "#c6a0f6", // macchiatto mauve
	User:     "#f5a97f", // macchiatto peach
	FQDN:     "#eed49f", // macchiatto yellow
	Dim:      "#494d64", // macchiatto surface 1
	Active:   "#a6da95", // macchiatto green
	Inactive: "#ee99a0", // macchiatto maroon
	Failed:   "#ed8796", // macchiatto red
	Green:    "#a6da95", // macchiatto green
	Orange:   "#f5a97f", // macchiatto peach
	Red:      "#ed8796", // macchiatto red
}

func main() {
	configPath := flag.String("config", "", "Path to config file")
	showHelp := flag.Bool("help", false, "Show help message")
	flag.BoolVar(showHelp, "h", false, "Show help message (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Repo: https://github.com/BOAScripts/logindash \n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "LoginDash - Shows system information on login\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  --config string\n")
		fmt.Fprintf(os.Stderr, "        Path to config file (default: ~/.config/login-dash/config.toml)\n")
		fmt.Fprintf(os.Stderr, "  -h, --help\n")
		fmt.Fprintf(os.Stderr, "        Show this help message\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -config /path/to/config.toml\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -help\n\n", os.Args[0])
	}

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	config := loadConfig(*configPath)
	// Initialize styles with colors from config (or defaults)
	initStyles(config.Colors)

	// Set label width from config if provided
	if config.Display.LabelWidth > 0 {
		labelWidth = config.Display.LabelWidth
	}

	if config.Display.GreenUntil > 0 {
		GreenUntil = config.Display.GreenUntil
	}

	if config.Display.OrangeUntil > 1 {
		OrangeUntil = config.Display.OrangeUntil
	}

	displayInfo(config)
}

func loadConfig(customPath string) Config {
	var config Config
	var configPath string

	if customPath != "" {
		configPath = customPath
	} else {
		currentUser, err := user.Current()
		if err == nil {
			configPath = filepath.Join(currentUser.HomeDir,
				".config", "logindash", "config.toml")
		}
	}

	if configPath != "" {
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			config = Config{}
		}
	}

	return config
}

func initStyles(colors ColorsConfig) {
	// Merge config colors with defaults
	if colors.Header == "" {
		colors.Header = defaultColors.Header
	}
	if colors.Title == "" {
		colors.Title = defaultColors.Title
	}
	if colors.Label == "" {
		colors.Label = defaultColors.Label
	}
	if colors.User == "" {
		colors.User = defaultColors.User
	}
	if colors.FQDN == "" {
		colors.FQDN = defaultColors.FQDN
	}
	if colors.Dim == "" {
		colors.Dim = defaultColors.Dim
	}
	if colors.Active == "" {
		colors.Active = defaultColors.Active
	}
	if colors.Inactive == "" {
		colors.Inactive = defaultColors.Inactive
	}
	if colors.Failed == "" {
		colors.Failed = defaultColors.Failed
	}
	if colors.Green == "" {
		colors.Green = defaultColors.Green
	}
	if colors.Orange == "" {
		colors.Orange = defaultColors.Orange
	}
	if colors.Red == "" {
		colors.Red = defaultColors.Red
	}

	// Initialize styles with structure + colors
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors.Header)).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colors.Header)).
		Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors.Title))

	labelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors.Label))

	userStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors.User))

	fqdnStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors.FQDN))

	dimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Dim))

	activeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Active))

	inactiveStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Inactive))

	failedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Failed))

	greenStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Green))

	orangeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Orange))

	redStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Red))

}

func displayInfo(config Config) {
	currentUser, _ := user.Current()
	fqdn := getFQDN()
	lastLogin, lastIP := getLastLogin(currentUser.Username)

	header := fmt.Sprintf("%s@%s\n%s",
		userStyle.Render(currentUser.Username),
		fqdnStyle.Render(fqdn),
		dimStyle.Render("Last login: "+lastLogin))

	// Only add "From:" line if there's an IP
	if lastIP != "" {
		header += fmt.Sprintf("\n%s", dimStyle.Render("From: "+lastIP))
	}

	fmt.Println(headerStyle.Render(header))
	fmt.Println()

	displaySystem()
	displayStorage(config.Disks.Paths)
	displayServices(config.Services.Monitored)

	fmt.Println(dimStyle.Render(strings.Repeat("─", 50)))
}

func displaySystem() {
	fmt.Println(titleStyle.Render("System"))

	OSInfo := getOSInfo()
	uptime := getUptime()
	cpu := getCPUUsage()
	cpucores := getCPUCores()
	memUsed, memTotal, memPercent := getMemoryUsage()
	iface := getDefaultInterface()
	ipAddr := getIPAddress(iface)
	gateway := getGateway()
	dns := getDNSServers()

	format := fmt.Sprintf("  %%s %%--%ds %%s\n", labelWidth)
	fmt.Printf(format, labelStyle.Render("▸"), "OS", OSInfo)
	fmt.Printf(format, labelStyle.Render("▸"), "Uptime", uptime)

	formatCPU := fmt.Sprintf("  %%s %%--%ds %%s %%s\n", labelWidth)
	fmt.Printf(formatCPU, labelStyle.Render("▸"), "CPU",
		cpucores, colorizePercentage(cpu))

	formatMem := fmt.Sprintf("  %%s %%--%ds %%s/%%s %%s\n", labelWidth)
	fmt.Printf(formatMem, labelStyle.Render("▸"), "RAM",
		memUsed, memTotal, colorizePercentage(memPercent))

	formatIP := fmt.Sprintf("  %%s %%--%ds %%s (%%s)\n", labelWidth)
	fmt.Printf(formatIP, labelStyle.Render("▸"), "IP Address",
		ipAddr, iface)
	fmt.Printf(format, labelStyle.Render("▸"), "Gateway", gateway)
	fmt.Printf(format, labelStyle.Render("▸"), "DNS", dns)

	fmt.Println()
}

func displayStorage(extraPaths []string) {
	fmt.Println(titleStyle.Render("Storage"))

	format := fmt.Sprintf("  %%s %%--%ds %%s/%%s %%s\n", labelWidth)

	used, total, percent := getDiskUsage("/")
	fmt.Printf(format, labelStyle.Render("▸"), "/",
		used, total, colorizePercentageStr(percent))

	for _, path := range extraPaths {
		if isMountPoint(path) {
			used, total, percent := getDiskUsage(path)
			label := path
			fmt.Printf(format, labelStyle.Render("▸"),
				label, used, total, colorizePercentageStr(percent))
		}
	}

	detectedMounts := autoDetectMounts()
	if len(detectedMounts) > 0 {
		fmt.Println()
		fmt.Println(titleStyle.Render("Storage (/mnt)"))
		for _, path := range detectedMounts {
			// Skip if path is already in extraPaths
			if slices.Contains(extraPaths, path) {
				continue
			}

			if isMountPoint(path) {
				used, total, percent := getDiskUsage(path)
				label := path
				fmt.Printf(format, labelStyle.Render("▸"),
					label, used, total, colorizePercentageStr(percent))
			}
		}
	}
	fmt.Println()
}

func displayServices(services []string) {
	if len(services) == 0 {
		return
	}

	fmt.Println(titleStyle.Render("Services"))

	format := fmt.Sprintf("  %%s %%s %%s\n")

	for _, service := range services {
		state, statusSince := getServiceStatus(service)

		var marker, styledStatusSince string
		switch state {
		case "active":
			marker = activeStyle.Render("●")
			styledStatusSince = dimStyle.Render(statusSince)
		case "inactive":
			marker = inactiveStyle.Render("○")
		default:
			marker = failedStyle.Render("✗")
			styledStatusSince = dimStyle.Render(statusSince)
		}

		label := service
		fmt.Printf(format, marker, label, styledStatusSince)
	}
	fmt.Println()
}

func colorizePercentage(percent float64) string {
	percentStr := fmt.Sprintf("(%.1f%%)", percent)
	if percent <= float64(GreenUntil) {
		return greenStyle.Render(percentStr)
	} else if percent <= float64(OrangeUntil) {
		return orangeStyle.Render(percentStr)
	}
	return redStyle.Render(percentStr)
}

func colorizePercentageStr(percentStr string) string {
	numStr := strings.TrimSuffix(percentStr, "%")
	percent, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return percentStr
	}

	percentFormatted := fmt.Sprintf("(%s)", percentStr)
	if percent <= float64(GreenUntil) {
		return greenStyle.Render(percentFormatted)
	} else if percent <= float64(OrangeUntil) {
		return orangeStyle.Render(percentFormatted)
	}
	return redStyle.Render(percentFormatted)
}

func getFQDN() string {
	out, _ := exec.Command("hostname", "-f").Output()
	return strings.TrimSpace(string(out))
}

func getLastLogin(username string) (string, string) {
	out, _ := exec.Command("last", "-n", "2", "-w", username).Output()
	lines := strings.Split(string(out), "\n")

	if len(lines) > 1 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 4 {
			ip := ""
			dateStart := 3

			// Check if field 2 looks like an IP address or hostname
			if len(fields) > 2 {
				potentialIP := fields[2]
				// Check if it's an IP address (contains dots and numbers)
				// or a hostname (not a date component)
				if strings.Contains(potentialIP, ".") ||
					strings.Contains(potentialIP, ":") {
					ip = potentialIP
					dateStart = 4
				}
			}

			dateInfo := []string{}
			for i := dateStart; i < len(fields); i++ {
				if fields[i] == "still" || fields[i] == "-" {
					break
				}
				dateInfo = append(dateInfo, fields[i])
			}

			if len(dateInfo) > 0 {
				return strings.Join(dateInfo, " "), ip
			}
		}
	}
	return "Unknown", ""
}

func getOSInfo() string {
	out, _ := exec.Command("bash", "-c",
		"cat /etc/os-release | grep 'PRETTY_NAME' | awk  -F '[\"\"]' '{print $2}'").Output()
	return strings.TrimSpace(string(out))
}

func getUptime() string {
	out, _ := exec.Command("uptime", "-p").Output()
	return strings.TrimPrefix(strings.TrimSpace(string(out)), "up ")
}

func getCPUCores() string {
	out, _ := exec.Command("bash", "-c",
		"lscpu | grep -E '^CPU\\(s\\):' | awk '{print $2}'").Output()
	return strings.TrimSpace(string(out))
}

func getCPUUsage() float64 {
	out, _ := exec.Command("bash", "-c",
		"top -bn1 | grep 'Cpu(s)' | awk '{print $2 + $4}'").Output()

	cpuStr := strings.TrimSpace(string(out))
	cpu, err := strconv.ParseFloat(cpuStr, 64)
	if err != nil {
		return 0
	}

	return cpu
}

func getMemoryUsage() (string, string, float64) {
	out, _ := exec.Command("free", "-h").Output()
	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Mem:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				total := fields[1]
				used := fields[2]

				outBytes, _ := exec.Command("free", "-b").Output()
				linesBytes := strings.Split(string(outBytes), "\n")
				for _, lineB := range linesBytes {
					if strings.HasPrefix(lineB, "Mem:") {
						fieldsB := strings.Fields(lineB)
						totalB, _ := strconv.ParseFloat(fieldsB[1], 64)
						usedB, _ := strconv.ParseFloat(fieldsB[2], 64)
						percent := (usedB / totalB) * 100
						return used, total, percent
					}
				}
			}
		}
	}
	return "0", "0", 0
}

func getDefaultInterface() string {
	out, _ := exec.Command("bash", "-c",
		"ip route | grep default | awk '{print $5}' | head -1").Output()
	return strings.TrimSpace(string(out))
}

func getIPAddress(iface string) string {
	out, _ := exec.Command("bash", "-c",
		fmt.Sprintf("ip -4 addr show %s | grep inet | "+
			"awk '{print $2}' | head -1", iface)).Output()
	return strings.TrimSpace(string(out))
}

func getGateway() string {
	out, _ := exec.Command("bash", "-c",
		"ip route | grep default | awk '{print $3}' | head -1").Output()
	return strings.TrimSpace(string(out))
}

func getDNSServers() string {
	out, _ := exec.Command("bash", "-c",
		"resolvectl status | grep -E '^\\s+DNS Servers:' | "+
			"awk '{print $3}' | tr '\\n' ' '").Output()
	dns := strings.TrimSpace(string(out))

	if dns == "" {
		out, _ = exec.Command("bash", "-c",
			"grep nameserver /etc/resolv.conf | "+
				"awk '{print $2}' | tr '\\n' ' '").Output()
		dns = strings.TrimSpace(string(out))
	}

	return dns
}

func getDiskUsage(path string) (string, string, string) {
	out, _ := exec.Command("df", "-h", path).Output()
	lines := strings.Split(string(out), "\n")
	if len(lines) >= 2 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 5 {
			return fields[2], fields[1], fields[4]
		}
	}
	return "0", "0", "0%"
}

func isMountPoint(path string) bool {
	out, err := exec.Command("mountpoint", "-q", path).CombinedOutput()
	return err == nil && len(out) == 0
}

func autoDetectMounts() []string {
	var mounts []string
	entries, _ := os.ReadDir("/mnt")
	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join("/mnt", entry.Name())
			if isMountPoint(path) {
				mounts = append(mounts, path)
			}
		}
	}
	return mounts
}

func getServiceStatus(service string) (string, string) {
	stateOut, _ := exec.Command("systemctl", "is-active", service).Output()
	state := strings.TrimSpace(string(stateOut))

	statusOut, _ := exec.Command("systemctl", "status", service).Output()
	re := regexp.MustCompile(`Active: ([^\s]+) \(([^)]+)\) since [^;]+; (.+)`)
	matches := re.FindStringSubmatch(string(statusOut))

	if len(matches) >= 4 {
		statusSince := fmt.Sprintf("since %s", matches[3])
		return state, statusSince
	}

	return state, ""
}
