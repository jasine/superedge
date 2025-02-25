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

package tcp

import (
	uuid "github.com/satori/go.uuid"
	"github.com/superedge/superedge/pkg/tunnel/conf"
	"github.com/superedge/superedge/pkg/tunnel/context"
	"github.com/superedge/superedge/pkg/tunnel/module"
	"github.com/superedge/superedge/pkg/tunnel/proxy/common"
	"github.com/superedge/superedge/pkg/tunnel/proxy/handlers"
	"github.com/superedge/superedge/pkg/tunnel/util"
	"k8s.io/klog/v2"
	"net"
)

type TcpProxy struct {
}

func (tcp *TcpProxy) Register(ctx *context.Context) {
	ctx.AddModule(tcp.Name())
}

func (tcp *TcpProxy) Name() string {
	return util.TCP
}

func (tcp *TcpProxy) Start(mode string) {
	context.GetContext().RegisterHandler(util.TCP_BACKEND, tcp.Name(), handlers.DirectHandler)
	context.GetContext().RegisterHandler(util.TCP_FRONTEND, tcp.Name(), handlers.FrontendHandler)
	context.GetContext().RegisterHandler(util.TCP_CONTROL, tcp.Name(), handlers.DirectHandler)
	if mode == util.CLOUD {
		for front, backend := range conf.TunnelConf.TunnlMode.Cloud.Tcp {
			go func(front, backend string) {
				ln, err := net.Listen(util.TCP, front)
				if err != nil {
					klog.Errorf("cloud proxy start %s fail ,error = %s", front, err)
					return
				}
				defer ln.Close()
				klog.Infof("the tcp server of the cloud tunnel listen on %s\n", front)
				for {
					rawConn, err := ln.Accept()
					if err != nil {
						klog.Errorf("cloud proxy accept error!")
						return
					}
					nodes := context.GetContext().GetNodes()
					if len(nodes) == 0 {
						rawConn.Close()
						klog.Errorf("len(nodes)==0")
						continue
					}
					uuid := uuid.NewV4().String()
					node := nodes[0]
					tcpConn := context.GetContext().AddConn(uuid)
					nodeConn := context.GetContext().GetNode(node)
					nodeConn.BindNode(uuid)
					go common.Read(rawConn, nodeConn, util.TCP, util.TCP_FRONTEND, uuid, backend)
					go common.Write(rawConn, tcpConn)
				}
			}(front, backend)
		}
	}
}

func (tcp *TcpProxy) CleanUp() {
	context.GetContext().RemoveModule(tcp.Name())
}

func InitTcp() {
	module.Register(&TcpProxy{})
	klog.Infof("init module: %s success !", util.TCP)
}
