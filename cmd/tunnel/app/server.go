/*
Copyright 2020 The SuperEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"github.com/spf13/cobra"
	"github.com/superedge/superedge/cmd/tunnel/app/options"
	"github.com/superedge/superedge/pkg/tunnel/conf"
	"github.com/superedge/superedge/pkg/tunnel/module"
	"github.com/superedge/superedge/pkg/tunnel/proxy/modules/egress"
	"github.com/superedge/superedge/pkg/tunnel/proxy/modules/https"
	"github.com/superedge/superedge/pkg/tunnel/proxy/modules/ssh"
	"github.com/superedge/superedge/pkg/tunnel/proxy/modules/stream"
	"github.com/superedge/superedge/pkg/tunnel/proxy/modules/tcp"
	"github.com/superedge/superedge/pkg/util"
	"github.com/superedge/superedge/pkg/version"
	"github.com/superedge/superedge/pkg/version/verflag"
	"k8s.io/klog/v2"
)

func NewTunnelCommand() *cobra.Command {
	option := options.NewTunnelOption()
	cmd := &cobra.Command{
		Use: "tunnel",
		Run: func(cmd *cobra.Command, args []string) {
			verflag.PrintAndExitIfRequested()

			klog.Infof("Versions: %#v\n", version.Get())
			util.PrintFlags(cmd.Flags())

			err := conf.InitConf(*option.TunnelMode, *option.TunnelConf)
			if err != nil {
				klog.Info("tunnel failed to load configuration file !")
				return
			}
			module.InitModules(*option.TunnelMode)
			stream.InitStream(*option.TunnelMode)
			egress.InitEgress()
			tcp.InitTcp()
			https.InitHttps()
			ssh.InitSSH()
			module.LoadModules(*option.TunnelMode)
			module.ShutDown()
		},
	}
	fs := cmd.Flags()
	namedFlagSets := option.Addflag()
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}
	return cmd
}
