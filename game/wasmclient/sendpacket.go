// Copyright 2014,2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wasmclient

import (
	"sync/atomic"
	"time"

	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
	"github.com/kasworld/gowasmlib/wrapspan"
)

func (app *WasmClient) reqAIPlay(onoff bool) error {
	return app.ReqWithRspFnWithAuth(
		c2t_idcmd.AIPlay,
		&c2t_obj.ReqAIPlay_data{onoff},
		func(hd c2t_packet.Header, rsp interface{}) error {
			return nil
		},
	)
}

func (app *WasmClient) reqAchieveInfo() error {
	return app.ReqWithRspFnWithAuth(
		c2t_idcmd.AchieveInfo,
		&c2t_obj.ReqAchieveInfo_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			rpk := rsp.(*c2t_obj.RspAchieveInfo_data)
			app.systemMessage.Append(wrapspan.ColorText("Gold",
				"== Achievement == "))
			for i, v := range rpk.Achieve {
				app.systemMessage.Append(wrapspan.ColorTextf("Gold",
					"%v : %v ", achievetype.AchieveType(i).String(), v))
			}
			return nil
		},
	)
}

func (app *WasmClient) reqHeartbeat() error {
	return app.ReqWithRspFnWithAuth(
		c2t_idcmd.Heartbeat,
		&c2t_obj.ReqHeartbeat_data{
			Time: time.Now(),
		},
		func(hd c2t_packet.Header, rsp interface{}) error {
			rpk := rsp.(*c2t_obj.RspHeartbeat_data)
			pingDur := time.Now().Sub(rpk.Time)
			app.PingDur = (app.PingDur + pingDur) / 2
			return nil
		},
	)
}

func (app *WasmClient) sendPacket(cmd c2t_idcmd.CommandID, arg interface{}) {
	if cmd.NeedTurn() != 0 { // is act?
		atomic.AddInt32(&app.actPacketPerTurn, 1)
	}
	app.ReqWithRspFnWithAuth(
		cmd, arg,
		func(hd c2t_packet.Header, rsp interface{}) error {
			return nil
		},
	)
}

func (app *WasmClient) sendMovePacketByInput(tryDir way9type.Way9Type) bool {
	cf := app.currentFloor()
	playerX, playerY := app.GetPlayerXY()
	moveDir := cf.FindMovableDir(playerX, playerY, tryDir)
	if moveDir != way9type.Center {
		atomic.AddInt32(&app.movePacketPerTurn, 1)
		go app.sendPacket(c2t_idcmd.Move,
			&c2t_obj.ReqMove_data{Dir: moveDir},
		)
		return true
	} else if tryDir != way9type.Center {
		app.systemMessage.Appendf(
			"Cannot move to %v", tryDir.String())
	}
	return false
}
