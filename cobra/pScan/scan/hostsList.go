package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	ErrEXists    = errors.New("host alraedy in the list")
	ErrNotEXists = errors.New("host not in the list")
)

type HostsList struct {
	Hosts []string
}

// Search searches for hosts in the list
func (hl *HostsList) search(host string) (bool, int) {
	sort.Strings(hl.Hosts)
	i := sort.SearchStrings(hl.Hosts, host)

	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}

	return false, -1
}

// Add adds a host to the list
func (hl *HostsList) Add(host string) error {
	if found, _ := hl.search(host); found {
		return fmt.Errorf("%w: %s", ErrEXists, host)
	}

	hl.Hosts = append(hl.Hosts, host)

	return nil
}

// Remove deletes a host from the list
func (hl *HostsList) Remove(host string) error {
	if found, i := hl.search(host); found {
		hl.Hosts = append(hl.Hosts[:i], hl.Hosts[i+1:]...)
		return nil
	}
	return fmt.Errorf("%w: %s", ErrNotEXists, host)
}

// Load obtains hosts from a hosts file
func (hl *HostsList) Load(hostsFile string) error {
	f, err := os.Open(hostsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}

	return nil
}

// Save saves hosts to a hosts file
func (hl *HostsList) Save(hostsFile string) error {
	output := ""

	for _, h := range hl.Hosts {
		output += fmt.Sprintln(h)
	}

	return os.WriteFile(hostsFile, []byte(output), 0664)
}
