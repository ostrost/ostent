package ostent

import (
	"encoding/json"
	"os"
	"os/user"
	"syscall"
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/ostrost/ostent/types"
)

type VagrantMachines struct {
	List types.VgmachineSlice
}

func LessVgmachine(a, b types.Vgmachine) bool { return a.UUID < b.UUID }

func vagrantmachines() (*VagrantMachines, error) {
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

	open, err := os.Open(indexFilename) // text, err := ioutil.ReadFile(indexFilename)
	if err != nil {
		return nil, err
	}

	status := new(struct {
		Machines *map[string]types.Vgmachine // the key is UUID
		// Version int // unused
	})
	if err := json.NewDecoder(open).Decode(status); err != nil { // json.Unmarshal(text, status)
		return nil, err
	}
	machines := new(VagrantMachines)
	if status.Machines != nil {
		for uuid, machine := range *status.Machines {
			machine.UUID = uuid
			machine.UUIDHTML = tooltipable(7, uuid)
			machine.VagrantfilePathHTML = tooltipable(50, machine.VagrantfilePath)
			machine.StateHTML = tooltipable(8, machine.State)
			// (*status.Machines)[uuid] = machine
			machines.List = append(machines.List, machine)
		}
	}
	machines.List.StableSortBy(LessVgmachine)
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
	iu := IndexUpdate{}
	if err != nil {
		iu.VagrantError = err.Error()
		iu.VagrantErrord = true
	} else {
		iu.VagrantMachines = machines
	}
	iUPDATES <- &iu
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

func vgwatch() (err error) {
	for {
		if err = vgchange(); err != nil {
			return err
		}
	}
}
