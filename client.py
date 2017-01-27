import socket
import sys
from threading import Thread

# Server IP and socket
UDP_IP = "10.224.240.67"
UDP_PORT = 5005

# Create socket
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.setblocking(0)

print "\nWelcome to fissechat!\nType !help to see commands.\n"
nickName = raw_input("Choose a nickname: ")

# Send connection code
s.sendto("1a92?#qQ,=11", (UDP_IP, UDP_PORT))


def helpCommand():
    print "--Commands: ",  "\n--!help shows this list", "\n--!list lists connected users", "\n--!name change your nickname", "\n--!quit quit"


def listCommand():
    print "list"


def nameCommand():
    global nickName
    oldName = nickName
    nickName = raw_input("Choose a nickname: ")
    s.sendto(oldName + " changed name to " + nickName, (UDP_IP, UDP_PORT))


def quitCommand():
    s.sendto("10n0m0001", (UDP_IP, UDP_PORT))
    sys.exit()

commandList = {"!help": helpCommand, "!list": listCommand, "!name": nameCommand, "!quit": quitCommand}


def inputChecker(input):
    try:
        if input[0] == "!":
            if input.split(" ")[0] in commandList:
                commandList.get(input.split(" ")[0])()
            return True
        else:
            return False
    except IndexError:
        return False


def chatListener():
    while True:
        try:
            data, address = s.recvfrom(1024)
            print data
        except socket.error:
            False

listenerThread = Thread(target=chatListener)
listenerThread.start()


def inputListener():
    while True:
        message = raw_input()
        if not inputChecker(message):
            s.sendto(nickName + ": " + message, (UDP_IP, UDP_PORT))

inputThread = Thread(target=inputListener)
inputThread.start()