package demo

import (
	"GoRestServer/model/auth"
	"GoRestServer/model"
	"github.com/fatih/structs"
	"github.com/rs/zerolog/log"
)

func TestSql() {
	functionModel := auth.FunctionModelInstance(model.SQL)
	//err := functionModel.Insert(&auth.Function{Name:"角色管理", Url:"/auth/role", IsMenu: false})
	//if err != nil {
	//	log.Warn().Err(err)
	//}

	//function := functionModel.FindSingle("fun_name=?", "用户管理")
	//log.Info().Fields(structs.Map(function)).Msg("find result")
	functions := functionModel.FindMore("1=1").([]*auth.Function)
	for i, fun := range functions  {
		log.Info().Fields(structs.Map(fun)).Msgf("find function result: %d", i)
	}

	roleModel := auth.RoleModelInstance(model.SQL)
	err := roleModel.Insert(&auth.Role{Name: "管理员6", RoleKey: "Administrator6", Functions: functions})
	if err != nil {
		log.Error().Err(err).Msg("failed.")
	}
	//roles := roleModel.FindMore("1=1").(*[]*auth.Role)
	//for i, role := range *roles  {
	//	log.Info().Fields(structs.Map(role)).Msgf("find result: %d", i)
	//}
	//role := roleModel.FindSingle("1=1")
	//log.Info().Fields(structs.Map(role)).Msg("find role result")
}
