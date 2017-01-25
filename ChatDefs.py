import re
import urllib
from Tkinter import *
from socket import *
import win32gui

def getExternalIP():
    url = "http://checkip.dyndns.org"
    request = urllib.urlopen(url).read()
    return str(re.findAll(r"\d{1,3]\.\d{1,3]\.\d{1,3]\.\d{1,3]", request))

def getInternalIP():
    return str(gethostbyname(getfqdn()))
