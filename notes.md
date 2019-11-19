layers:
    * world
        - top level
        - consumes messages from server
            - handles all messages
    * server
        - mid level
        - consumes messages from client
            - handles ErrorMessage
            - bubbles all other messages up
        - produces messages for world
    * client
        - low level
        - handles commands (input from user)
        - handles writes   (output to user)
        - produces messages for server

events flow up:
    * client -> server -> world

onConnect:
    * user connects
    * client requests username
    * client sends welcome chat
    * client produces:  ClientStartMessage
    * server bubbles:   ClientStartMessage
    * world  consumes:  ClientStartMessage
        - broadcasts user joined chat

onExit:
    * user types exit command
    * client closes network socket
    * client produces: ClientStopMessage
    * client ends go routines for that connection
    * server deletes client from tracking
    * server bubbles: ClientStopMessage
    * world handles ClientStopMessage
        - broadcasts user left chat
    
    