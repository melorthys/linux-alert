package main

import (
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func cpu() {
	cpu := "top -bn2 | grep '%Cpu' | tail -1 | grep -P '(....|...) id,'|awk '{print 100-$8}'"

	cpu_output, err := exec.Command("/bin/sh", "-c", cpu).Output()
	if err != nil {
		fmt.Println("error")
	}

	cpu_string_output := strings.Trim(string(cpu_output), "\n")

	cpu_usage, err := strconv.ParseFloat(cpu_string_output, 8)
	if err != nil {
		fmt.Println("parsing error")
	}

	fmt.Print("CPU  = % " + string(cpu_output))
	if cpu_usage >= 80.0 {
		fmt.Println("CPU usage is over 80%, sending alert email!")
		sendEmail("CPU", cpu_string_output)
	}
}

func ram() {
	ram := "free | grep Mem | awk '{ printf(\"%.1f\", $3/$2 * 100.0) }'"

	ram_output, err := exec.Command("/bin/sh", "-c", ram).Output()
	if err != nil {
		fmt.Println("error")
	}

	ram_usage, err := strconv.ParseFloat(string(ram_output), 8)
	if err != nil {
		fmt.Println("parsing error")
	}

	fmt.Println("RAM  = % " + string(ram_output))
	if ram_usage >= 80.0 {
		fmt.Println("RAM usage is over 80%, sending alert email!")
		sendEmail("RAM", string(ram_output))
	}
}

func disk() {
	disk := "df -h --output=pcent / | awk 'NR==2{print $1}' | rev | cut -c 2- | rev"

	disk_output, err := exec.Command("/bin/sh", "-c", disk).Output()
	if err != nil {
		fmt.Println("error")
	}

	disk_string_output := strings.Trim(string(disk_output), "\n")

	disk_usage, err := strconv.ParseFloat(string(disk_string_output), 8)
	if err != nil {
		fmt.Println("parsing error")
	}

	fmt.Print("DISK = % " + string(disk_output))
	if disk_usage >= 80.0 {
		fmt.Println("DISK usage is over 80%, sending alert email!")
		sendEmail("DISK", disk_string_output)
	}
}

func sendEmail(metric string, usage string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from_name := os.Getenv("MAIL_FROM_NAME")
	from_email := os.Getenv("MAIL_FROM_EMAIL")
	to_name := os.Getenv("MAIL_TO_NAME")
	to_email := os.Getenv("MAIL_TO_EMAIL")
	subj := "[ALERT] High " + metric + " usage!"

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}

	current_time := time.Now().Format(time.RFC1123)

	conn, error := net.Dial("udp", "8.8.8.8:80")
	if error != nil {
		fmt.Println(error)

	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()

	body := metric + " usage is " + usage + "% for " + hostname + " (" + ip + ") at " + current_time

	auth := smtp.PlainAuth("", user, pass, host)

	// NOTE: Using the backtick here ` works like a heredoc, which is why all the
	// rest of the lines are forced to the beginning of the line, otherwise the
	// formatting is wrong for the RFC 822 style
	message := `To: "` + to_name + `" <` + to_email + `>
From: "` + from_name + `" <` + from_email + `>
Subject: ` + subj + `

` + body + `
`

	err = smtp.SendMail(host+":"+port, auth, from_email, []string{to_email}, []byte(message))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	for {
		fmt.Println("Getting system metrics..")
		cpu()
		ram()
		disk()
		fmt.Println("Sleeping for 10 minutes..\n")
		time.Sleep(10 * time.Minute)
	}
}
