module radithor

go 1.26.3

require golang.org/x/sys v0.49.0

require shared v0.0.0

require github.com/horobimasu/gotils v0.0.0-20261009122752-1a78f85ce208

require (
	github.com/bwmarrin/discordgo v0.29.0
	github.com/horobimasu/windowsx v0.0.0-20261017091157-0ca212c235a1
	github.com/kbinani/screenshot v0.0.0-20250624051815-089614a94018
)

require (
	github.com/gen2brain/shm v0.1.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
)

replace shared => ../shared
