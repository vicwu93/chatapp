module main

go 1.15

replace example.com/constants => ../constants

require (
	example.com/constants v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.4.2
)
