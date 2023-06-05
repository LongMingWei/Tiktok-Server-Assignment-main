# TikTok Tech Immersion 2023 Server (Backend) Assignment 

![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

Done by: Long Ming Wei 

About: This code is an implementation of the backend portion of an Instant Messaging system, handling 
requests from users to send to and receive messages from other users. The http-server is responsible 
for receiving HTTP requests from clients and makes an RPC call to the rpc-server. The rpc-server then 
reads messages from an SQL database if the HTTP requests is PULL and writes messages to update the 
database if the HTTP request is SEND. The messages read/written will then be sent back to the http-server
and then to the client as feedback. 

- main.go in http-server initializes itself and uses the API in idl_http.proto to tell the rpc-server whether the request is
sendMessage() (api/send) or pullMessage() (api/pull) and transmits details of the request as well such
as message content and chat/roomID. 
- main.go in rpc-server initializes itself and the SQL database and client with the help of docker-compose.yml
- db.go in rpc-server is the database client and is in charge of saving and getting messages from the 
SQL database based on the request sent by the rpc-server.
- handler.go in rpc-server is in charge of handling the request by breaking down its details such as 
message and sender and telling the database client to save or get messages from its database.
