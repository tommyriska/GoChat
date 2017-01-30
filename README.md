#Future features

## Client

#### Design:
- [x] Add ASCII-art
- [ ] Let users choose text color for their own chat
- [x] Your own name shows up before your message
- [ ] Clear the screen before startup/welcome message
- [ ] Timestamps on messages
- [ ] Create a GUI

#### Networking
- [x] Change from UDP to TCP connections

#### Other features
- [ ] Add some sort of encryption to secure all data is private
- [ ] Let users specify ip and port to connect to
- [x] Add !help command
- [x] Add !list command
- [x] Add !quit command
- [ ] Add command functions for client module

#### Known problems
- Server wont accept two connections at the same time
- Cannot restart server on the same host or port as last initiation
- 

## Server

#### Design:
- [x] Add ASCII-art
- [ ] Let users choose text color for their own chat
- [x] Your own name shows up before your message
- [ ] Clear the screen before startup/welcome message
- [ ] Timestamps on messages
- [ ] Create a GUI

#### Networking
- [x] Change from UDP to TCP connections

#### Other features
- [ ] Add some sort of encryption to secure all data is private
- [ ] Let users specify ip and port to connect to
- [x] Add !help command for client
- [ ] Add !list command for client
- [x] Add !quit command for client
- [ ] Add command functions for client module

#### Known problems
- Server wont accept two connections at the same time
    - FIXED 30.01.17: Had to set socket options so that the socket made the address [reusable](https://docs.python.org/2/library/socket.html#socket.socket.getsockopt).
- Cannot restart server on the same host or port as last initiation 
- 

#### Suggestions for gitflow

- Everyone makes branches for one fix
- We always merge into development branch before pushing a full release to master
- Merges to master should always contain a functional version of this software