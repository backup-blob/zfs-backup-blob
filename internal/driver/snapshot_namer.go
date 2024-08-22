package driver

import (
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"regexp"
	"strings"
	"time"
)

type defaultNamer struct {
	nower domain.Nower
}

const SnapPrefix = "backup_blob_"

func NewDefaultNamer(nower domain.Nower) domain.SnapshotNamestrategy {
	return &defaultNamer{nower: nower}
}

// IsGreater returns true if snapNameA is greater than snapNameB.
func (d *defaultNamer) IsGreater(snapNameA, snapNameB string) bool {
	return snapNameA > snapNameB
}

// dateFormatted returns a RFC3339 dated with : replaced by - to make it uri compliant.
func (d *defaultNamer) dateFormatted() string {
	return strings.ReplaceAll(d.nower.Now().UTC().Format(time.RFC3339), ":", "-")
}

// GetName returns a snapshot name with prefix (example: backup_blob_2024-04-12T16-28-26Z).
func (d *defaultNamer) GetName() string {
	return fmt.Sprintf(SnapPrefix+"%s", d.dateFormatted())
}

// IsMatching checks if a snapshotName is matching the naming schema.
func (d *defaultNamer) IsMatching(snapshotName string) bool {
	pattern := `^` + SnapPrefix + `\d{4}\-\d{2}\-\d{2}T\d{2}\-\d{2}\-\d{2}Z$`
	matched, err := regexp.MatchString(pattern, snapshotName)

	return err == nil && matched
}
