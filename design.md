# baremud

## objects:

### network:

#### server
* handles all networking
* adds/removes clients for connections
* handles messages to/from clients
* holds inbox for messages from clients
* TODO: outbox for messages to clients

#### client
* handles network for 1 connection
* holds Player
* passes messages to server inbox
* can read/write to user

### data:

#### world
* handles top level game logic
* holds Rooms, Gates, Players
* TODO: inbox from server

#### room
* handles game logic for single location
* holds Gates, Players

#### gate
* handles one-way game logic link to a room
* holds Room (destination)

#### player
* handles game logic for 1 user
* holds Room

## flow:

* user connects with telnet
    - server creates a client which creates a new player
* user logs in with creds
    - server handles login
    - player data is updated
    - player is sent through a start gate into a room
    - room is displayed to client
* user can enter commands
    - look, exit, chat
    - messages from clients go to server inbox to be handled
* user exits
    - server removes player from world
    - server disconnects and deletes client
    - server broadcasts that player has left





























