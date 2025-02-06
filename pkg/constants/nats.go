package constants

type connectionPoolKey string

const ServicePrefix = "kittens"

const LobbyRoot = "lobby"

const LobbyStream = ServicePrefix + "-" + LobbyRoot

var ConnectionPool = connectionPoolKey("edgeConnectionPool")
