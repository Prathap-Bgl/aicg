# **Overview**
This directory houses code for three distinct applications, each with its own purpose and technology stack. The applications demonstrate the practical use of scripting, web development, and network communication.

### Script: Network Latency Checker
*Directory: script*

This directory contains a simple Bash script designed to monitor network latency. The script pings the well-known Google DNS server at 8.8.8.8. If the response time exceeds 20 milliseconds, it automatically prints all relevant network configuration details of the server. This tool is useful for quick network diagnostics and monitoring.

*Usage:*

    Execute the Bash script in the script directory to perform the latency check and network configuration display.

### Web Application: Age-Based Access Control
*Directory: webapp*

This application, written in Python, features a basic web interface with multiple functionalities:

    *Admin Interface (/admin): Allows the addition of entries with fields for name, age, and gender.*
    *User Interface (/user): Users can enter a name to check if the person is eligible to consume alcohol based on their age.*
    *Test Interface (/run-tests): A separate interface to run test cases for the application.*

*Usage:*

    Run app.py from the webapp directory.
    Access the interfaces through a web browser:
        Admin Interface: http://127.0.0.1:5000/admin
        User Interface: http://127.0.0.1:5000/user
        Test Interface: http://127.0.0.1:5000/run-tests

### Network: Basic DHCP Server
*Directory: network*

This directory features a DHCP server application written in Go (Golang). The server listens for DHCP discovery messages. Upon receiving a message, it checks for an entry corresponding to the source MAC address in a local file-based database. If a matching entry is found, the server generates a DHCP response that includes the bootfile name indicated by the file associated with the MAC address.

*Usage:*

    First, build the applications.
    Run the DHCP server using ./dhcp.
    To test, execute ./dhcp-test.
    The list of successful dhcp discover matches can be viewed by http://127.0.0.1/list 

Each application in this directory showcases different aspects of software development, from scripting and web development to network programming, providing a diverse set of tools for various IT needs.
