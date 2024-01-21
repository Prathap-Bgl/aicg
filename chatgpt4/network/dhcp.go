package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/krolaw/dhcp4"
)

var (
    successfulSearches []string
    mutex              sync.Mutex
)

// Start a web server to list successful DHCP discoveries
func startWebServer() {
    http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
        mutex.Lock()
        defer mutex.Unlock()

        for _, search := range successfulSearches {
            fmt.Fprintf(w, "%s\n", search)
        }
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Search the MAC address in the file database
func searchDatabase(macAddress string) (string, bool) {
    file, err := os.Open("database.txt")
    if err != nil {
        log.Println(err)
        return "", false
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.Split(line, " ")
        if len(parts) == 2 && strings.EqualFold(parts[0], macAddress) {
            return parts[1], true
        }
    }

    if err := scanner.Err(); err != nil {
        log.Println(err)
    }

    return "", false
}

// Log a successful search
func logSuccessfulSearch(macAddress string, filename string) {
    mutex.Lock()
    successfulSearches = append(successfulSearches, fmt.Sprintf("%s %s %s", time.Now().Format(time.RFC3339), macAddress, filename))
    mutex.Unlock()
}

func handleDHCPPacket(conn *net.UDPConn) {
    buffer := make([]byte, 1500) // DHCP packets are typically less than 1500 bytes
    n, addr, err := conn.ReadFromUDP(buffer)
    if err != nil {
        log.Println(err)
        return
    }

    packet := dhcp4.Packet(buffer[:n])
    options := packet.ParseOptions()
    if options[dhcp4.OptionDHCPMessageType] != nil && options[dhcp4.OptionDHCPMessageType][0] == byte(dhcp4.Discover) {
        macAddress := packet.CHAddr().String()
        if filename, found := searchDatabase(macAddress); found {
            logSuccessfulSearch(macAddress, filename)

            // Normally, you would construct and send a DHCP offer here
            // For simplicity, this example just logs the success
            fmt.Printf("Success: %s found for MAC %s\n", filename, macAddress)

            // Constructing and sending a DHCP offer (highly simplified)
            response := dhcp4.NewPacket(dhcp4.BootReply)
            response.SetXId(packet.XId())
            response.SetYIAddr(net.IP{192, 168, 1, 100}) // Example IP, change as needed
            response.SetSIAddr(net.IP{192, 168, 1, 1})   // Example Server IP, change as needed
            response.SetGIAddr(packet.GIAddr())
            response.SetCHAddr(packet.CHAddr())
            response.AddOption(dhcp4.OptionDHCPMessageType, []byte{byte(dhcp4.Offer)})
            response.AddOption(dhcp4.OptionServerIdentifier, []byte{192, 168, 1, 1}) // Example Server IP, change as needed
            response.AddOption(dhcp4.OptionRouter, []byte{192, 168, 1, 1})           // Example Router IP, change as needed
            response.AddOption(dhcp4.OptionSubnetMask, []byte{255, 255, 255, 0})
            response.AddOption(dhcp4.OptionDomainNameServer, []byte{8, 8, 8, 8}) // Example DNS, change as needed
            response.AddOption(dhcp4.OptionHostName, []byte(filename))            // Filename as Hostname
            // Send the response
            _, err := conn.WriteToUDP(response, addr)
            if err != nil {
                log.Println("Failed to send DHCP response:", err)
            }
        }
    }
}

func main() {
    go startWebServer()

    addr := net.UDPAddr{
        Port: 12000,
        IP:   net.ParseIP("0.0.0.0"),
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    for {
        handleDHCPPacket(conn)
    }
}

