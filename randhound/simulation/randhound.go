package main
//gopkg.in/dedis/onet.v1
import (
	"github.com/BurntSushi/toml"
	"mobilehound/randhound"
	"mobilehound/onet"
	"mobilehound/log"
	"mobilehound/simul"
	"mobilehound/simul/monitor"
)

func init() {
	onet.SimulationRegister("RandHound", NewRHSimulation)
}

// RHSimulation implements a RandHound simulation
type RHSimulation struct {
	onet.SimulationBFTree
	Groups    int
	GroupSize int
	Faulty    int
	Purpose   string
}

// NewRHSimulation creates a new RandHound simulation
func NewRHSimulation(config string) (onet.Simulation, error) {
	rhs := &RHSimulation{}
	_, err := toml.Decode(config, rhs)
	if err != nil {
		return nil, err
	}
	return rhs, nil
}

// Setup configures a RandHound simulation with certain parameters
func (rhs *RHSimulation) Setup(dir string, hosts []string) (*onet.SimulationConfig, error) {
	sim := new(onet.SimulationConfig)
	rhs.CreateRoster(sim, hosts, 2000)
	err := rhs.CreateTree(sim)
	return sim, err
}

// Run initiates a RandHound simulation
func (rhs *RHSimulation) Run(config *onet.SimulationConfig) error {
	randM := monitor.NewTimeMeasure("tgen-randhound")
	bandW := monitor.NewCounterIOMeasure("bw-randhound", config.Server)
	client, err := config.Overlay.CreateProtocol("RandHound", config.Tree, onet.NilServiceID)
	if err != nil {
		return err
	}
	rh, _ := client.(*randhound.RandHound)
	if rhs.Groups == 0 {
		if rhs.GroupSize == 0 {
			log.Fatal("Need either Groups or GroupSize")
		}
		rhs.Groups = rhs.Hosts / rhs.GroupSize
	}
	err = rh.Setup(rhs.Hosts, rhs.Faulty, rhs.Groups, rhs.Purpose)
	if err != nil {
		return err
	}
	if err := rh.Start(); err != nil {
		log.Error("Error while starting protcol:", err)
	}

	select {
	case <-rh.Done:
		log.Lvlf1("RandHound - done")
		random, transcript, err := rh.Random()
		if err != nil {
			return err
		}
		randM.Record()
		bandW.Record()
		log.Lvlf1("RandHound - collective randomness: ok")

		verifyM := monitor.NewTimeMeasure("tver-randhound")
		err = rh.Verify(rh.Suite(), random, transcript)
		if err != nil {
			return err
		}
		verifyM.Record()
		log.Lvlf1("RandHound - verification: ok")

		//case <-time.After(time.Second * time.Duration(rhs.Hosts) * 5):
		//log.Print("RandHound - time out")
	}

	return nil

}

func main() {
	simul.Start()
}
