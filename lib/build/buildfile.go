package build

// import (
// 	"encoding/json"
// 	"io/ioutil"
// )

// // read:
// // https://medium.com/@matryer/golang-advent-calendar-day-eleven-persisting-go-objects-to-disk-7caf1ee3d11d

// // Buildfile ...
// type Buildfile struct {
// 	ContainerPort string
// 	HostPort      string
// }

// // Load ...
// func (bf *Buildfile) Load(fpath string) error {
// 	f, e := ioutil.ReadFile(fpath)
// 	if e != nil {
// 		return e
// 	}

// 	e = json.Unmarshal(f, bf)
// 	if e != nil {
// 		return e
// 	}

// 	return nil
// }
