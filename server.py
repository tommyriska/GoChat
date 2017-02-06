import socket

# Server IP and socket
UDP_IP = "10.224.209.209"
UDP_PORT = 5005

# Create and bind socket
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.bind((UDP_IP, UDP_PORT))

print "Server is running at ", UDP_IP, ": ", UDP_PORT

# List that holds connections
connections = []

# Server loop
while True:
    # Listen for message
    data, address = s.recvfrom(1024)

    # Check for connection code
    if data == "1a92?#qQ,=11":
        # Add IP to connections
        if address not in connections:
            connections.append(address)
            print address, " connected"

    elif data == "10n0m0001":
        if address in connections:
            connections.remove(address)
            print address, " disconnected"

    # If not, print message
    else:
        print address, " ", data

        # Send message to all connections
        for c in connections:
            if c != address:
                s.sendto(data, c)