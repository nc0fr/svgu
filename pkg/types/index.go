// Copyright Nicolas Paul (2023)
//
// * Nicolas Paul
//
// This software is a computer program whose purpose is to allow the hosting
// and sharing of Go modules using a personal domain.
//
// This software is governed by the CeCILL license under French law and
// abiding by the rules of distribution of free software.  You can  use,
// modify and/ or redistribute the software under the terms of the CeCILL
// license as circulated by CEA, CNRS and INRIA at the following URL
// "http://www.cecill.info".
//
// As a counterpart to the access to the source code and  rights to copy,
// modify and redistribute granted by the license, users are provided only
// with a limited warranty  and the software's author,  the holder of the
// economic rights,  and the successive licensors  have only  limited
// liability.
//
// In this respect, the user's attention is drawn to the risks associated
// with loading,  using,  modifying and/or developing or reproducing the
// software by the user in light of its specific status of free software,
// that may mean  that it is complicated to manipulate,  and  that  also
// therefore means  that it is reserved for developers  and  experienced
// professionals having in-depth computer knowledge. Users are therefore
// encouraged to load and test the software's suitability as regards their
// requirements in conditions enabling the security of their systems and/or
// data to be ensured and,  more generally, to use and operate it in the
// same conditions as regards security.
//
// The fact that you are presently reading this means that you have had
// knowledge of the CeCILL license and that you accept its terms.

package types

import (
	"go.nc0.fr/svgu/pkg/templates"
	"os"
	"path"
	"sync"
)

// Index is the global object representing the Starlark configuration.
type Index struct {
	Domain  string
	Modules map[string]*Module
	// internal
	lock sync.Mutex
}

// SetDomain sets the domain of the index.
func (i *Index) SetDomain(d string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.Domain = d
}

// AddModule adds a module to the index.
func (i *Index) AddModule(n string, m *Module) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.Modules[n] = m
}

// GetModule returns a module from the index.
func (i *Index) GetModule(n string) *Module {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.Modules[n]
}

// RemoveModule removes a module from the index.
func (i *Index) RemoveModule(n string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	delete(i.Modules, n)
}

// CheckModule checks if a module is in the index.
func (i *Index) CheckModule(n string) bool {
	i.lock.Lock()
	defer i.lock.Unlock()
	_, ok := i.Modules[n]
	return ok
}

// GenerateFile generates the index file.
func (i *Index) GenerateFile(out string) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	f := path.Join(out, "index.html")

	// Create the file.
	fd, err := os.Create(f)
	if err != nil {
		return err
	}
	defer func(fd *os.File) {
		err := fd.Close()
		if err != nil {
			panic(err)
		}
	}(fd)

	// Execute the template and write the output to the file.
	if err := templates.ExecIndex(fd,
		"https://pkg.go.dev", 2); err != nil {
		return err
	}

	return nil
}
