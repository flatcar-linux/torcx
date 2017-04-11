// Copyright 2017 CoreOS Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"io/ioutil"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/coreos/torcx/pkg/torcx"
)

var (
	cmdProfile = &cobra.Command{
		Use:   "profile [command]",
		Short: "Operate on local profile(s)",
		Long:  `This subcommand operates on local profile(s).`,
	}
)

func init() {
	TorcxCmd.AddCommand(cmdProfile)
}

// fillProfileRuntime generates the runtime config for profile subcommands,
// starting from system-wide state and config
func fillProfileRuntime(commonCfg *torcx.CommonConfig) (*torcx.ProfileConfig, error) {
	var (
		curProfileName string
		curProfilePath string
		nextProfile    string
	)

	if commonCfg == nil {
		return nil, errors.New("missing common configuration")
	}

	cpn, err := torcx.CurrentProfileName()
	if err == nil {
		curProfileName = cpn
	}
	cpp, err := torcx.CurrentProfilePath()
	if err == nil {
		curProfilePath = cpp
	}

	fc, err := ioutil.ReadFile(commonCfg.ConfProfile())
	if err == nil {
		nextProfile = strings.TrimSpace(string(fc))
	}
	if nextProfile == "" {
		nextProfile = "vendor"
		logrus.Debug("no next profile configured, assuming default")
	}

	logrus.WithFields(logrus.Fields{
		"current profile": curProfileName,
		"next profile":    nextProfile,
	}).Debug("profile configuration parsed")

	return &torcx.ProfileConfig{
		CommonConfig:       *commonCfg,
		CurrentProfileName: curProfileName,
		CurrentProfilePath: curProfilePath,
		NextProfile:        nextProfile,
	}, nil
}
