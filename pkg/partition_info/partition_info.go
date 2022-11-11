package partition_info

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jaypipes/ghw/pkg/block"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// PartitionInfo maps a partition label to a partition UUID.
// It's used in order to be able to ask the kcrypt-challenger for the passphrase
// using the partition label, even when the label is not accessible (e.g. before
// decrypting the partition). The UUID can be used to lookup the partition label
// and make the request.
type PartitionInfo map[string]string

func PartitionToString(p *block.Partition) string {
	return fmt.Sprintf("%s:%s:%s", p.Label, p.Name, p.UUID)
}

// Takes a partition info string (as returned by PartitionToString) and return
// the partition label and the UUID
func PartitionDataFromString(partitionStr string) (string, string) {
	parts := strings.Split(partitionStr, ":")

	return parts[0], parts[2]
}

// UpdatePartitionLabelMapping takes partition information as a string argument
// the the form: `label:name:uuid` (that's what the `kcrypt encrypt` command returns
// on success. This function stores it in the PartitionInfoFile yaml file for
// later use.
func UpdatePartitionLabelMapping(partitionData, file string) error {
	partitionInfo, err := ParsePartitionInfoFile(file)
	if err != nil {
		return err
	}

	label, uuid := PartitionDataFromString(partitionData)
	partitionInfo[label] = uuid

	return UpdatePartitionInfoFile(partitionInfo, file)
}

func ParsePartitionInfoFile(file string) (PartitionInfo, error) {
	var result PartitionInfo

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return result, errors.Wrap(err, "reading the partition info file")
	}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return result, errors.Wrap(err, "unmarshalling partition info file")
	}

	return result, nil
}

func UpdatePartitionInfoFile(partitionInfo PartitionInfo, file string) error {
	data, err := yaml.Marshal(&partitionInfo)
	if err != nil {
		return errors.Wrap(err, "marshalling the new partition info to yaml")
	}
	err = ioutil.WriteFile(file, data, 0)
	if err != nil {
		return errors.Wrap(err, "writing back the partition info file")
	}
	return nil
}
