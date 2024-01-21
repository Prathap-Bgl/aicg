package main

import (
    "fmt"
    "log"
    "net"
    "time"

    "github.com/krolaw/dhcp4"
)

func main() {
    // Setup a DHCP Discover packet
    packet := dhcp4.NewPacket(dhcp4.BootRequest)
    packet.SetXId([]byte{0xAB, 0xCD, 0xEF, 0x01}) // Transaction ID
    packet.SetCHAddr(net.HardwareAddr{0x00, 0x15, 0x5d, 0x10, 0x3a, 0x1c}) // Client MAC address
    packet.AddOption(dhcp4.OptionDHCPMessageType, []byte{byte(dhcp4.Discover)})
    packet.PadToMinSize()

    // Connect to the DHCP server
    serverAddr := "127.0.0.1:12000" // Replace with your server's IP and port
    conn, err := net.Dial("udp", serverAddr)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Send the packet
    _, err = conn.Write(packet)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("DHCP Discover packet sent")

    // Wait and listen for a response
    buffer := make([]byte, 1500) // DHCP packets are typically less than 1500 bytes
    conn.SetReadDeadline(time.Now().Add(10 * time.Second)) // 10-second timeout for the response

    n, err := conn.Read(buffer)
    if err != nil {
        log.Fatal("Failed to read from UDP:", err)
    }

    responsePacket := dhcp4.Packet(buffer[:n])
    printPacket(responsePacket)
}

// Function to print the DHCP packet in a readable format
func printPacket(packet dhcp4.Packet) {
    fmt.Println("Received DHCP Packet:")
    fmt.Printf("OpCode: %d\n", packet.OpCode())
    fmt.Printf("Hardware Type: %d\n", packet.HType())
    fmt.Printf("Hardware Address Length: %d\n", packet.HLen())
    fmt.Printf("Hops: %d\n", packet.Hops())
    fmt.Printf("Transaction ID: %x\n", packet.XId())
    fmt.Printf("Seconds Elapsed: %d\n", packet.Secs())
    fmt.Printf("Bootp Flags: %x\n", packet.Flags())
    fmt.Printf("Client IP Address: %s\n", packet.CIAddr())
    fmt.Printf("Your (client) IP Address: %s\n", packet.YIAddr())
    fmt.Printf("Next Server IP Address: %s\n", packet.SIAddr())
    fmt.Printf("Relay Agent IP Address: %s\n", packet.GIAddr())
    fmt.Printf("Client MAC Address: %s\n", packet.CHAddr())

    // Print DHCP Options
    fmt.Println("Options:")
    for optionType, optionData := range packet.ParseOptions() {
        fmt.Printf(" - Option: %d, Data: %x\n", optionType, optionData)
    }
}

