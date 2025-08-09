package update

import (
	"fmt"
	"os"

	selfupdate "github.com/rhysd/go-github-selfupdate/selfupdate"
)

const repo = "BT-Pro-Solutions/textonly-cli"

func CheckAndApply(checkOnly bool, currentVersion string) (string, bool, error) {
	rel, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		return "", false, err
	}
	if !found || rel == nil || rel.Version.String() == "" || rel.Version.String() == currentVersion {
		return "", false, nil
	}
	if checkOnly {
		return rel.Version.String(), true, nil
	}
	exe, err := os.Executable()
	if err != nil {
		return "", false, fmt.Errorf("cannot determine executable: %w", err)
	}
	if err := selfupdate.UpdateTo(rel.AssetURL, exe); err != nil {
		return "", false, fmt.Errorf("update failed: %w", err)
	}
	return rel.Version.String(), true, nil
}
