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

package activeobject

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/enum/achievetype_stats"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype_stats"
	"github.com/kasworld/goguelike/enum/potiontype_stats"
	"github.com/kasworld/goguelike/enum/scrolltype_stats"
	"github.com/kasworld/goguelike/game/activeobject/serverai2"
	"github.com/kasworld/goguelike/game/visitarea"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd_stats"
)

// function for web

func (ao *ActiveObject) GetBornFaction() factiontype.FactionType {
	return ao.bornFaction
}

func (ao *ActiveObject) GetSP() float64 {
	return ao.sp
}

func (ao *ActiveObject) GetDeath() int {
	return int(ao.achieveStat.Get(achievetype.Death))
}

func (ao *ActiveObject) GetKill() int {
	return int(ao.achieveStat.Get(achievetype.Kill))
}

func (ao *ActiveObject) GetActStats() *c2t_idcmd_stats.CommandIDStat {
	return &ao.aoActionStat
}

func (ao *ActiveObject) GetPotionStat() *potiontype_stats.PotionTypeStat {
	return &ao.potionStat
}

func (ao *ActiveObject) GetAIObj() *serverai2.ServerAI {
	return ao.ai
}

func (ao *ActiveObject) Web_ActiveObjInfo(w http.ResponseWriter, r *http.Request) {
	tplIndex, err := template.New(
		"index").Funcs(
		c2t_idcmd_stats.IndexFn).Funcs(
		potiontype_stats.IndexFn).Funcs(
		scrolltype_stats.IndexFn).Funcs(
		fieldobjacttype_stats.IndexFn).Funcs(
		achievetype_stats.IndexFn).Parse(`
	<html> <head>
	<title>ActiveObject</title>
	</head>
	<body>
	{{.}}
	</br>
	<a href= "/KickActiveObj?aoid={{.GetUUID}}" >
		KickActiveObj
	</a>
	</br>
	Level : {{.GetTurnData.Level}}
	</br>
	Exp : {{.GetTurnData.TotalExp}} 
	</br>
	NonBattle {{.GetTurnData.NonBattleExp}}
	</br>
	Kill {{.GetKill}}
	</br>
	Death {{.GetDeath}}
	</br>
	HP : {{.GetHP}}
	</br>
	SP : {{.GetSP}}
	</br>
	Sight : {{.GetTurnData.Sight}}
	</br>
	Bias : {{.GetBias}}
	</br>
	BornFaction : {{.GetBornFaction}}
	</br>
	AtkBias {{.GetTurnData.AttackBias}} 
	</br>
	DefBias {{.GetTurnData.DefenceBias}} 
	<br/>
	{{if .GetAIObj }} 
		AI : {{.GetAIObj}} 
		<br/>
		AI Dur : {{.GetAIObj.GetAIDur}} 
		<br/>
		AI Plans : {{.GetAIObj.GetPlanNameList}} 
	{{end}}
	<br/>
	LoadRate {{.GetTurnData.LoadRate}}
	<br/>
	TotalWeight {{.GetInven.GetTotalWeight}}
	<br/>
	Wallet {{.GetInven.GetWalletValue}}
	<hr/>
	Equipped 
	<br/>
	{{range $i,$v := .GetInven.GetEquipSlot}}
		{{if $v}}
			{{$i}} {{$v}}
		<br/>
		{{end}}
	{{end}}

	EquipBag
	<br/>
	{{range $i,$v := .GetInven.GetEquipList}}
		{{if $v}}
			{{$i}} {{$v}}
		<br/>
		{{end}}
	{{end}}

	PotionBag
	<br/>
	{{range $i,$v := .GetInven.GetPotionList}}
		{{if $v}}
			{{$i}} {{$v}}
		<br/>
		{{end}}
	{{end}}

	ScrollBag
	<br/>
	{{range $i,$v := .GetInven.GetScrollList}}
		{{if $v}}
			{{$i}} {{$v}}
		<br/>
		{{end}}
	{{end}}

	<br/>
	{{with .GetAchieveStat}}
		<table border=1 style="border-collapse:collapse;">
		` + achievetype_stats.HTML_tableheader + `
		{{range $i, $v := .}}
		` + achievetype_stats.HTML_row + `
		{{end}}
		` + achievetype_stats.HTML_tableheader + `
		</table>
	{{end}}
	{{with .GetPotionStat}}
		<table border=1 style="border-collapse:collapse;">
		` + potiontype_stats.HTML_tableheader + `
		{{range $i, $v := .}}
		` + potiontype_stats.HTML_row + `
		{{end}}
		` + potiontype_stats.HTML_tableheader + `
		</table>
	{{end}}
	{{with .GetScrollStat}}
		<table border=1 style="border-collapse:collapse;">
		` + scrolltype_stats.HTML_tableheader + `
		{{range $i, $v := .}}
		` + scrolltype_stats.HTML_row + `
		{{end}}
		` + scrolltype_stats.HTML_tableheader + `
		</table>
	{{end}}
	{{with .GetFieldObjActStat}}
		<table border=1 style="border-collapse:collapse;">
		` + fieldobjacttype_stats.HTML_tableheader + `
		{{range $i, $v := .}}
		` + fieldobjacttype_stats.HTML_row + `
		{{end}}
		` + fieldobjacttype_stats.HTML_tableheader + `
		</table>
	{{end}}
	{{with .GetActStats}}
		<table border=1 style="border-collapse:collapse;">
		` + c2t_idcmd_stats.HTML_tableheader + `
		{{range $i, $v := .}}
		` + c2t_idcmd_stats.HTML_row + `
		{{end}}
		` + c2t_idcmd_stats.HTML_tableheader + `
		</table>
	{{end}}
	{{range $i, $v := .GetVisitFloorList}}
		{{if $v}}
			{{$i}} {{$v}}
			<br/>
			<img src="/ActiveObjVisitImgae?aoid={{$.GetUUID}}&floorid={{$v.GetUUID}}" >			
			<br/>
		{{end}}
	{{end}}
	<hr/>
	</body> </html> 
	`)
	if err != nil {
		ao.log.Error("%v", err)
		fmt.Println(err)
	}
	if err := tplIndex.Execute(w, ao); err != nil {
		ao.log.Error("%v", err)
		fmt.Println(err)
	}
}

func (ao *ActiveObject) GetVisitFloorList() []*visitarea.VisitArea {
	return ao.uuid2VisitArea.GetList()
}

func (ao *ActiveObject) GetVisitFloor(floorid string) *visitarea.VisitArea {
	r, _ := ao.uuid2VisitArea.GetByID(floorid)
	return r
}
