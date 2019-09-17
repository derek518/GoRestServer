package demo

import (
	"io/ioutil"
	"github.com/rs/zerolog/log"
	"encoding/json"
	"GoRestServer/helper"
)

func JsonTest() error {
	content, err := ioutil.ReadFile("query.json")
	if err != nil {
		log.Error().Err(err)
		return err
	}

	var result interface{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	log.Debug().Fields(result.(map[string]interface{})).Msg("json test result")

	where := helper.BuildSqlWhere(result.(map[string]interface{}))

	log.Debug().Fields(where).Msg("sql where")
	
	return nil
}
