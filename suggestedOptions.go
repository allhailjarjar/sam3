package sam3

import (
	"net"
	"net/http"
	"os"
)

// Examples and suggestions for options when creating sessions.
var (
	// Suitable options if you are shuffling A LOT of traffic. If unused, this
	// will waste your resources.
	Options_Humongous = []string{"inbound.length=3", "outbound.length=3",
		"inbound.lengthVariance=1", "outbound.lengthVariance=1",
		"inbound.backupQuantity=3", "outbound.backupQuantity=3",
		"inbound.quantity=6", "outbound.quantity=6"}

	// Suitable for shuffling a lot of traffic.
	Options_Large = []string{"inbound.length=3", "outbound.length=3",
		"inbound.lengthVariance=1", "outbound.lengthVariance=1",
		"inbound.backupQuantity=1", "outbound.backupQuantity=1",
		"inbound.quantity=4", "outbound.quantity=4"}

	// Suitable for shuffling a lot of traffic quickly with minimum
	// anonymity. Uses 1 hop and multiple tunnels.
	Options_Wide = []string{"inbound.length=1", "outbound.length=1",
		"inbound.lengthVariance=1", "outbound.lengthVariance=1",
		"inbound.backupQuantity=2", "outbound.backupQuantity=2",
		"inbound.quantity=3", "outbound.quantity=3"}

	// Suitable for shuffling medium amounts of traffic.
	Options_Medium = []string{"inbound.length=3", "outbound.length=3",
		"inbound.lengthVariance=1", "outbound.lengthVariance=1",
		"inbound.backupQuantity=0", "outbound.backupQuantity=0",
		"inbound.quantity=2", "outbound.quantity=2"}

	// Sensible defaults for most people
	Options_Default = []string{"inbound.length=3", "outbound.length=3",
		"inbound.lengthVariance=0", "outbound.lengthVariance=0",
		"inbound.backupQuantity=1", "outbound.backupQuantity=1",
		"inbound.quantity=1", "outbound.quantity=1"}

	// Suitable only for small dataflows, and very short lasting connections:
	// You only have one tunnel in each direction, so if any of the nodes
	// through which any of your two tunnels pass through go offline, there will
	// be a complete halt in the dataflow, until a new tunnel is built.
	Options_Small = []string{"inbound.length=3", "outbound.length=3",
		"inbound.lengthVariance=1", "outbound.lengthVariance=1",
		"inbound.backupQuantity=0", "outbound.backupQuantity=0",
		"inbound.quantity=1", "outbound.quantity=1"}

	// Does not use any anonymization, you connect directly to others tunnel
	// endpoints, thus revealing your identity but not theirs. Use this only
	// if you don't care.
	Options_Warning_ZeroHop = []string{"inbound.length=0", "outbound.length=0",
		"inbound.lengthVariance=0", "outbound.lengthVariance=0",
		"inbound.backupQuantity=0", "outbound.backupQuantity=0",
		"inbound.quantity=2", "outbound.quantity=2"}
)

func PrimarySessionString() string {
	_, err := http.Get("http://127.0.0.1:7070")
	if err != nil {
		_, err := http.Get("http://127.0.0.1:7657")
		if err != nil {
			return "MASTER"
		}
		// at this point we're probably running on Java I2P and thus probably
		// have a PRIMARY session. Just to be sure, try to make one, check
		// for errors, then immediately close it.
		testSam, err := NewSAM(SAMDefaultAddr(""))
		if err != nil {
			return "MASTER"
		}
		newKeys, err := testSam.NewKeys()
		if err != nil {
			return "MASTER"
		}
		primarySession, err := testSam.newPrimarySession("PRIMARY", "primaryTestTunnel", newKeys, Options_Small)
		if err != nil {
			return "MASTER"
		}
		primarySession.Close()
		return "PRIMARY"
	}
	return "MASTER"
}

var PrimarySessionSwitch string = PrimarySessionString()

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

var SAM_HOST = getEnv("sam_host", "127.0.0.1")
var SAM_PORT = getEnv("sam_port", "7656")

func SAMDefaultAddr(fallforward string) string {
	if fallforward == "" {
		return net.JoinHostPort(SAM_HOST, SAM_PORT)
	}
	return fallforward
}
