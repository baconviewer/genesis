package deploy

import (
	"../blockchains/helpers"
	"../db"
	netem "../net"
	"../ssh"
	"../testnet"
)

// PurgeTestNetwork goes into each given ssh client and removes all the nodes and the networks.
// Increments the build state len(clients) * 2 times and sets it stag to tearing down network,
// if buildState is non nil.
func PurgeTestNetwork(tn *testnet.TestNet) error {
	if tn.BuildState != nil {
		tn.BuildState.SetBuildStage("Tearing down the previous testnet")
	}
	DockerStopServices(tn)
	return helpers.AllServerExecCon(tn, func(client *ssh.Client, server *db.Server) error {
		DockerKillAll(client)
		if tn.BuildState != nil {
			tn.BuildState.IncrementDeployProgress()
		}
		DockerNetworkDestroyAll(client)
		if tn.BuildState != nil {
			tn.BuildState.IncrementDeployProgress()
		}
		netem.RemoveAllOnServer(client, server.Nodes)

		return nil
	})
}
// Destroy is an alias of PurgeTestNetwork
func Destroy(tn *testnet.TestNet) error {
	return PurgeTestNetwork(tn)
}
