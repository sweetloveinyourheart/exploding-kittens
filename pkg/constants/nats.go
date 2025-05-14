package constants

type busKey string
type connectionPoolKey string

var Bus = busKey("bus")

const ServicePrefix = "kittens"

const LobbyRoot = "lobby"
const GameRoot = "game"
const DeskRoot = "desk"
const HandRoot = "hand"

const LobbyStream = ServicePrefix + "-" + LobbyRoot
const GameStream = ServicePrefix + "-" + GameRoot
const DeskStream = ServicePrefix + "-" + DeskRoot
const HandStream = ServicePrefix + "-" + HandRoot

var ConnectionPool = connectionPoolKey("connectionPool")

const NatsChannelBufferSize = 1000
