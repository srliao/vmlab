package defaults

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func RelDirPath(skip int) (string, error) {
	_, filename, _, _ := runtime.Caller(skip)
	s := strings.SplitAfterN(filename, gitRepoBasePath, 2)
	if len(s) != 2 {
		return "", fmt.Errorf("could not determine directory from caller: %s", filename)
	}
	return filepath.Dir(s[1]), nil
}
