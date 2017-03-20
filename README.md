
# What is PythonChat?
A studentproject where the goal is to develop a chat-service which will use AES encryption to protect the users and their messages. 


# How does PythonChat work?
PythonChat is ment to be an encrypted group chat were clients can send messages safely over the web. The chat is implemented with code that is considered safe, and uses public key encryption (RSA). 

Our first prototype is an UDP-version of an echo server. The following functions are implemented: 

* the server listens for new connections
* recieves messages
* sends messages to other connected clients 


# Who uses PythonChat?
We target people that care about their own privacy and want a fast, simple and safe way to communicate with other people. 


# Instructions for developers
Our project is developed on Python and GoLang. The first prototype is only written in Python, but we have decided to use GoLang for further development. Feedback from the community have stated that use of threads and encryption libraries are easy to use in GoLang, and that GoLang supports multithreading. GoLang is also a very fast programming language, and we want to take advantage of that. 

You should check the following before submitting any code to the project. Be aware that you need to be added by admin before contributing. 

* code is tested 
* code is well documented with comments to each method


# License
This project is developed under the GNU General Public License (GPL) 2.0. 

