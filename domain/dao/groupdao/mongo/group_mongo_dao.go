/*
@date : 2020/03/16
@author : YaPi
@desc :
*/
package mongo

import (
	"tinyUrl/common/constants"
	"tinyUrl/config"
	"tinyUrl/config/db/mongo"
	"tinyUrl/domain/entity"
)

func AddGroup(group *entity.Group) (e error) {
	var (
		err error
	)
	if err = mongo.DbInsert(config.Base.Mongo.DbName, constants.TinyGroup, group); err != nil {
		return err
	} else {
		return nil
	}
}

func UpdateGroup(q map[string]interface{}, set map[string]interface{}) (e error) {
	var (
		err error
	)
	if err = mongo.DbUpdateOne(config.Base.Mongo.DbName, constants.TinyGroup, q, set); err != nil {
		return err
	} else {
		return nil
	}
}

func GetGroupByMap(m map[string]interface{}) (t []entity.Group, e error) {
	var (
		groups []entity.Group
		err    error
	)
	if err = mongo.DbFind(config.Base.Mongo.DbName, constants.TinyGroup, m, &groups); err != nil {
		return groups, err
	} else {
		return groups, nil
	}
}
