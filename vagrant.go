package ostent

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"os"
	"os/user"
	"sort"
	"syscall"
	"time"

	"github.com/howeyc/fsnotify"
)

type vagrantMachine struct {
	UUID     string
	UUIDHTML template.HTML // !

	Vagrantfile_pathHTML template.HTML // !
	Vagrantfile_path     string
	Local_data_path      string

	Name     string
	Provider string
	State    string

	// 	Vagrantfile_name *[]string   // unused
	// 	Updated_at         *string   // unused
	// 	Extra_data         *struct { // unused
	//		Box *struct{
	//			Name     *string
	//			Provider *string
	//			Version  *string
	//		}
	//	}
}

type vagrantMachines struct {
	List []vagrantMachine
}

// satisfying sort.Interface interface
func (ms vagrantMachines) Len() int {
	return len(ms.List)
}

// satisfying sort.Interface interface
func (ms vagrantMachines) Swap(i, j int) {
	ms.List[i], ms.List[j] = ms.List[j], ms.List[i]
}

// satisfying sort.Interface interface
func (ms vagrantMachines) Less(i, j int) bool {
	return ms.List[i].UUID < ms.List[j].UUID
}

func vagrantmachines() (*vagrantMachines, error) {
	currentUser, _ := user.Current()
	lockFilename := currentUser.HomeDir + "/.vagrant.d/data/machine-index/index.lock"
	indexFilename := currentUser.HomeDir + "/.vagrant.d/data/machine-index/index"

	lock_file, err := os.Open(lockFilename)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(lock_file.Fd()), syscall.LOCK_EX); err != nil {
		return nil, err
	}
	defer syscall.Flock(int(lock_file.Fd()), syscall.LOCK_UN)

	text, err := ioutil.ReadFile(indexFilename)
	if err != nil {
		return nil, err
	}

	status := new(struct {
		Machines *map[string]vagrantMachine // the key is UUID
		// Version int // unused
	})
	if err := json.Unmarshal(text, status); err != nil {
		return nil, err
	}
	machines := new(vagrantMachines)
	if status.Machines != nil {
		for uuid, machine := range *status.Machines {
			machine.UUID = uuid
			machine.UUIDHTML = tooltipable(7, uuid)
			machine.Vagrantfile_pathHTML = tooltipable(50, machine.Vagrantfile_path)
			// (*status.Machines)[uuid] = machine
			machines.List = append(machines.List, machine)
		}
	}
	sort.Stable(machines)
	return machines, nil
}

var vgmachineindexFilename string

func init() {
	currentUser, _ := user.Current()
	vgmachineindexFilename = currentUser.HomeDir + "/.vagrant.d/data/machine-index/index"
}

func vgdispatch() { // (*fsnotify.FileEvent)
	machines, err := vagrantmachines()
	if err != nil { // an inconsistent write by vagrant? (although not with the flock)
		return // ignoring the error
	}
	pu := pageUpdate{}
	if err != nil {
		pu.VagrantError = err.Error()
		pu.VagrantErrord = true
	} else {
		pu.VagrantMachines = machines
	}
	pUPDATES <- &pu
}

func vgchange() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := watcher.Watch(vgmachineindexFilename); err != nil {
		return err
	}

	stop := make(chan struct{}, 1)
	go func() {
		<-watcher.Event
		go vgdispatch()
		stop <- struct{}{}
	}()
	<-stop
	time.Sleep(time.Second) // to overcome possible fsnotify data races
	watcher.Close()
	return nil
}

func vgwatch() error {
	for {
		if err := vgchange(); err != nil {
			return err
		}
	}
	return nil
}
