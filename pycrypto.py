#note that this will not run unless Pycrypto is installed

from Crypto.Cipher import AES
import os

#make a key and set a message to be encrypted
key = os.urandom(16)
message = 'A super secret message.'

#create a cipher object
cipher = AES.new(key)

def pad(s):
    return s + ((16-len(s) % 16) * '{')

def encrypt(plaintext):
    global cipher
    return cipher.encrypt(pad(plaintext))

def decrypt(ciphertext):
    global cipher
    dec = cipher.decrypt(ciphertext).decode('utf-8')
    l = dec.count('{')
    return dec[:len(dec)-l]

print("*** PYCRYPTO ***")
print("Message: " + message)
encrypted = encrypt(message)
decrypted = decrypt(encrypted)
print("Encrypted: " + encrypted)
print("Decrytped: " + decrypted)

