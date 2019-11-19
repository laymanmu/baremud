layers:
    * world
        - consumes messages from server
    * server
        - consumes messages from client
            - handles ErrorMessage
            - bubbles all other messages up
        - produces messages for world
    * client
        - produces messages for server

events flow up:
    * client -> server -> world
   


onConnect:
    * client requests username
    * client sends welcome chat
    * client produces:  ClientStartMessage
    * server bubbles:   ClientStartMessage
    * world  consumers: ClientStartMessage
        - broadcasts user joined chat
    
    