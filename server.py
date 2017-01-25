import thread
from ChatDefs import *

### INITIATE CONNECTION VARIABLES ###
WindowTitle = "tChat 0.1 - Server"
s = socket(AF_INET, SOCK_STREAM)
HOST = gethostname()
PORT = 8011
conn = ''
s.bind((HOST, PORT))

### CONNECTION MANAGEMENT ###
def getConnected():
    s.listen(1)
    global conn
    conn, addr = s.accept()
    loadConnectionInfo(chatLog, 'Connected with: ' + str(addr) + '\n--------------------------------------------------')

    while 1:
        try:
            data = conn.recv(1024)
            loadOtherEntry(chatLog, data)
            if base.focus_get() == None:
                flashMyWindow(WindowTitle)
        except:
            loadConnectionInfo(chatLog, '\n [ Your partner has disconnected ]\n [ Waiting for him to connect..] \n')
            getConnected()
    conn.close()

thread.start_new_thread(getConnected,())