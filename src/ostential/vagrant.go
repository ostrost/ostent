package ostential
import (
	"os"
	"syscall"
	"os/user"
	"io/ioutil"
	"encoding/json"
	"html/template"
)

type vagrantExtra struct {
	Box *struct{
		Name     *string
		Provider *string
		Version  * string
	}
}

type vagrantMachine struct {
	            UUID     string
	            UUIDHTML template.HTML // !
	Vagrantfile_pathHTML template.HTML // !

	Local_data_path      string
	Name                 string
	Provider             string
	State                string
	Vagrantfile_path     string

// 	Vagrantfile_name *[]string       // unused
// 	Updated_at         *string       // unused
// 	Extra_data         *vagrantExtra // unused
}

type vagrantMachinesMap map[string]vagrantMachine // the key is UUID

type vagrantStatus struct {
	Machines *vagrantMachinesMap
	Version int // unused
}
type vagrantMachines struct {
	List []vagrantMachine
}

func vagrantmachines() (*vagrantMachines, error) {
	currentUser, _ := user.Current()
	lock_filename  := currentUser.HomeDir + "/.vagrant.d/data/machine-index/index.lock"
	index_filename := currentUser.HomeDir + "/.vagrant.d/data/machine-index/index"

	lock_file, err := os.Open(lock_filename);
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(lock_file.Fd()), syscall.LOCK_EX); err != nil {
		return nil, err
	}
	defer syscall.Flock(int(lock_file.Fd()), syscall.LOCK_UN)

	text, err := ioutil.ReadFile(index_filename)
	if err != nil {
		return nil, err
	}

	status := new(vagrantStatus)
	if err := json.Unmarshal(text, status); err != nil {
		return nil, err
	}
	machines := new(vagrantMachines)
	if status.Machines != nil {
		for uuid, machine := range *status.Machines {
			machine.UUID                 = uuid
			machine.UUIDHTML             = tooltipable(7, uuid)
			machine.Vagrantfile_pathHTML = tooltipable(50, machine.Vagrantfile_path)
			// (*status.Machines)[uuid] = machine
			machines.List = append(machines.List, machine)
		}
	}
	return machines, nil
}
