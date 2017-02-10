import socket
import sys
from threading import Thread
from colorama import *
import time


# CLIENT NEEDS A WAY TO SPECIFY IP AND PORT, REWRITE!
# Server IP and socket
TCP_IP = "192.168.0.20" #IP to connect to
TCP_PORT = 5005 #Port to connect to
BUFFERSIZE = 1024


# Create socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((TCP_IP, TCP_PORT))
# Send connection code
s.send("1a92?#qQ,=11")

#Startup message
print (Fore.LIGHTMAGENTA_EX + """//////////////////////////////////////////
//  _____        _____ _           _    //
// |  __ \      / ____| |         | |   //
// | |__) |   _| |    | |__   __ _| |_  //
// |  ___/ | | | |    | '_ \ / _` | __| //
// | |   | |_| | |____| | | | (_| | |_  //
// |_|    \__, |\_____|_| |_|\__,_|\__| //
//         __/ |                        //
//        |___/                         //
//////////////////////////////////////////""")
print (Style.RESET_ALL)
print " "
print "\nWelcome to PyChat!\nType !help to see commands.\n"
nickName = raw_input("Choose a nickname: ")

def helpCommand():
    print "--Commands: ",  "\n--!help shows this list", "\n--!list lists connected users", "\n--!name change your nickname", "\n--!quit quit"


def listCommand():
    print "list"


def nameCommand():
    global nickName
    oldName = nickName
    nickName = raw_input("Choose a nickname: ")
    s.sendall(oldName + " changed name to " + nickName)


def quitCommand():
    s.sendall("10n0m0001", (TCP_IP, TCP_PORT))
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
            data, address = s.recv(1024)
            print data
        except socket.error:
            return False

listenerThread = Thread(target=chatListener)
listenerThread.start()


def inputListener():
    while True:
        message = raw_input(nickName + "-> ")
        if not inputChecker(message):
            s.sendall(nickName + ": " + message)

inputThread = Thread(target=inputListener)
inputThread.start()
