package main

import (
	"fmt"
	"log"

	"github.com/blushft/go-diagrams/diagram"
	"github.com/blushft/go-diagrams/nodes/generic"
)

func main() {
	d, err := diagram.New(diagram.Label("my-network"), diagram.Filename("network"))
	if err != nil {
		log.Fatal(err)
	}

	connect := func(start *diagram.Node, end *diagram.Node) {
		d.Connect(start, end, diagram.Bidirectional())
	}

	decix := generic.Place.Datacenter().Label("DE/CIX")
	edge := generic.Network.Switch().Label("NAT Router")
	connect(decix, edge)
	mesh := generic.Place.Datacenter().Label("NYCMesh")
	connect(mesh, edge)
	mesh1340 := generic.Network.Router().Label("nycmesh 1340")
	connect(mesh, mesh1340)
	lightbeam := generic.Network.Router().Label("nycmesh-lbe-1659")
	connect(mesh1340, lightbeam)

	meshGroup := diagram.NewGroup("NYC Mesh").Label("NYC Mesh")
	meshGroup.Add(edge, mesh, mesh1340)
	d.Group(meshGroup)

	gw := generic.Network.Firewall().Label("kibble")
	airport := generic.Network.Switch().Label("Airport Express")
	belkin := generic.Network.Router().Label("Belkin WiFi")
	roku := generic.Device.Tablet().Label("Roku TV")

	aptGroup := diagram.NewGroup("Home").Label("Home")
	aptGroup.Add(airport, belkin, roku)
	d.Group(aptGroup)

	kibbleGroup := diagram.NewGroup("Kibble").Label("Kibble")
	kibbleGroup.Add(gw)
	aptGroup.Group(kibbleGroup)

	bridges := [4]*diagram.Node{}
	for i := range bridges {
		if i == 1 || i == 2 {
			continue
		}
		name := fmt.Sprintf("Bridge %d", i)
		bridge := generic.Compute.Rack().Label(name)
		kibbleGroup.Add(bridge)
		connect(gw, bridge)
		bridges[i] = bridge
	}

	// physicalEthers := [5]*diagram.Node{}
	// for i := range physicalEthers {
	// 	name := fmt.Sprintf("em%d", i)
	// 	physicalEthers[i] = generic.Blank.Blank().Label(name)
	// 	bridges[0].Add(physicalEthers[i])
	// }

	emLabel := func(i int) func(eo *diagram.EdgeOptions) {
		return func(eo *diagram.EdgeOptions) {
			eo.Label = fmt.Sprintf("em%d", i)
			eo.Forward = true
			eo.Reverse = true
		}
	}
	d.Connect(bridges[0], airport, emLabel(1))
	d.Connect(bridges[0], belkin, emLabel(3))
	connect(belkin, roku)

	d.Connect(lightbeam, gw, func(eo *diagram.EdgeOptions) {
		eo.Label = "em0"
		eo.Forward = true
		eo.Reverse = true
	})

	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
}
