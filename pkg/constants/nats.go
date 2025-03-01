package constants

type busKey string
type connectionPoolKey string

var Bus = busKey("bus")

const ServicePrefix = "kittens"

const LobbyRoot = "lobby"

const LobbyStream = ServicePrefix + "-" + LobbyRoot

var ConnectionPool = connectionPoolKey("connectionPool")

const NatsChannelBufferSize = 1000
