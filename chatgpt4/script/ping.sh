#!/bin/bash

# Install ethtool if not present
if ! command -v ethtool &> /dev/null; then
    echo "Installing ethtool..."
    sudo apt-get install ethtool || sudo yum install ethtool
fi

# Function to perform a single ping and return the latency
function ping_latency() {
    local result=$(ping -c 1 8.8.8.8)
    #echo "Debug: Ping result: $result"   # Debug statement

    if [[ $? -eq 0 ]]; then
        local latency=$(echo "$result" | grep 'time=' | awk -F'time=' '{print $2}' | awk '{print $1}')
        #echo "Debug: Extracted latency: $latency"   # Debug statement
        echo $latency
    else
        echo "failed"
    fi
}

# Ping multiple times and store latencies
echo "Pinging 8.8.8.8..."
declare -a latencies
total_latency=0
valid_counts=0

for i in {1..4}; do
    latency=$(ping_latency)
    #echo $latency
    if [[ $latency != "Ping failed" ]]; then
        latencies[$i]=$latency
        total_latency=$(echo "$total_latency + $latency" | bc -l)
        valid_counts=$((valid_counts + 1))
    else
        latencies[$i]="1000"
    fi
done

# Display latencies in a table format
echo -e "\nPing Latency (ms):"
printf "%-10s %-10s\n" "Attempt" "Latency"
for i in {1..4}; do
    printf "%-10s %-10s\n" "$i" "${latencies[$i]}"
done

# Calculate average latency if there were successful pings
if [ $valid_counts -gt 0 ]; then
    avg_latency=$(echo "$total_latency / $valid_counts" | bc -l)

    # Check if the average ping time is greater than 20ms
    if (( $(echo "$avg_latency > 20" | bc -l) )); then
        echo -e "\nAverage Latency is greater than 20ms. Dumping network configurations..."
        echo "----------------------------------------"

        # Dump network interface configuration
        echo -e "\n### IFCONFIG ###"
        ifconfig
        echo "----------------------------------------"

        # Dump routing table
        echo -e "\n### ROUTE ###"
        route -n
        echo "----------------------------------------"

        # Dump DNS configuration
        echo -e "\n### DNS CONFIG ###"
        cat /etc/resolv.conf
        echo "----------------------------------------"

        # Perform traceroute
        echo -e "\n### TRACEROUTE ###"
        traceroute 8.8.8.8
        echo "----------------------------------------"

        # Dump IP stack related parameters
        echo -e "\n### IP STACK PARAMETERS ###"
        sysctl -a | grep net.ipv4
        echo "----------------------------------------"

        # Dump interface details
        echo -e "\n### INTERFACE DETAILS ###"
        for iface in $(ls /sys/class/net/ | grep -v lo); do 
            echo -e "\nInterface: $iface"
            ethtool $iface
        done
        echo "----------------------------------------"

        # Dump traffic control parameters
        echo -e "\n### TRAFFIC CONTROL PARAMETERS ###"
        tc -s qdisc
        echo "----------------------------------------"

    else
        echo -e "\nAverage Latency is under 20ms. Network seems fine."
    fi
else
    echo -e "\nAll ping attempts failed. Unable to calculate average latency."
fi

