package cloudflare

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Identity struct {
	Name   string
	Token  string
	ZoneId string
	Email  string
	Domain string
}

type Identities []Identity

func (ids Identities) MarshalJSON() ([]byte, error) {
	// We want a structure like this
	// "name": {
	// 		"token": <str>
	//		"zone_id": <str>
	// }
	result := make(map[string]map[string]string)
	for _, id := range ids {
		result[id.Name] = map[string]string{
			"token":   id.Token,
			"zone_id": id.ZoneId,
			"email":   id.Email,
			"domain":  id.Domain,
		}
	}

	return json.Marshal(result)
}

func (ids *Identities) UnmarshalJSON(data []byte) error {
	// Unmarshal the weird structure we got
	var tempMap map[string]map[string]string
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	// Remember that we need to modify the pointer
	*ids = make(Identities, 0, len(tempMap))
	// Iterate over the key, value pairs
	for name, details := range tempMap {
		id := Identity{
			Name:   name,
			Token:  details["token"],
			ZoneId: details["zone_id"],
			Email:  details["email"],
			Domain: details["domain"],
		}
		*ids = append(*ids, id)
	}

	return nil
}

func (ids *Identities) Get(name string) (*Identity, error) {
	for _, identity := range *ids {
		if strings.EqualFold(identity.Name, name) {
			return &identity, nil
		}
	}
	return nil, errors.New("Identity not found")
}

func (ids *Identities) Add(id *Identity) {
	*ids = append(*ids, *id)
}

func (ids *Identities) Remove(name string) error {
	for i, identity := range *ids {
		if strings.EqualFold(identity.Name, name) {
			// Attach the two "slices" together
			*ids = append((*ids)[:i], (*ids)[i+1:]...)
			return nil
		}
	}
	return errors.New("identity not found")
}

func (ids *Identities) Save() {
	out, err := json.Marshal(*ids)
	if err != nil {
		fmt.Printf("Could not save identities to disk: %v", err)
		return
	}

	identiesFile, err := getUserDataFile()
	if err != nil {
		fmt.Printf("Could not retrieve user data location: %v", err)
		return
	}

	dir := filepath.Dir(identiesFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Could not create user data directory: %v", err)
		return
	}

	if err = os.WriteFile(identiesFile, out, 0644); err != nil {
		fmt.Printf("Could not write identities to disk: %v", err)
	}
}

func ReadIdentities() (*Identities, error) {
	identiesFile, err := getUserDataFile()
	if err != nil {
		return &Identities{}, nil
	}

	var identities Identities
	bytes, err := os.ReadFile(identiesFile)
	if err != nil {
		identities = Identities{}
		if os.IsNotExist(err) {
			identities.Save()
		} else {
			return &identities, err
		}
	} else if err := json.Unmarshal(bytes, &identities); err != nil {
		return &Identities{}, err
	}

	return &identities, nil
}

func getUserDataFile() (string, error) {
	userDataDir, err := getUserDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userDataDir, "smokescreen", "identities.json"), nil
}

func getUserDataDir() (string, error) {
	// Apparently this env variable is the standard location for user data
	// (in my WSL Debian instance is not present though)
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome != "" {
		return dataHome, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".local", "share"), nil
}
