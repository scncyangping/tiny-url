/*
@Time : 2019-06-14 16:41
@Author : yangping
@File : tinyInfoDao
@Desc :
*/
package mongo

import (
	"tinyUrl/common/constants"
	"tinyUrl/config"
	"tinyUrl/config/db/mongo"
	"tinyUrl/domain/entity"
)

func AddTinyInfo(tinyInfo *entity.TinyInfo) (e error) {
	var (
		err error
	)
	if err = mongo.DbInsert(config.Base.Mongo.DbName, constants.TinyInfo, tinyInfo); err != nil {
		return err
	} else {
		return nil
	}
}

func GetTinyInfoById(id string) (t entity.TinyInfo, e error) {
	var (
		tinyInfo entity.TinyInfo
		err      error
	)
	if err = mongo.DbFindById(config.Base.Mongo.DbName, constants.TinyInfo, id, &tinyInfo); err != nil {
		return tinyInfo, err
	} else {
		return tinyInfo, nil
	}
}

func GetTinyByMap(m map[string]interface{}) (t entity.TinyInfo, e error) {
	var (
		tinyInfo entity.TinyInfo
		err      error
	)
	if err = mongo.DbFindOne(config.Base.Mongo.DbName, constants.TinyInfo, m, &tinyInfo); err != nil {
		return tinyInfo, err
	} else {
		return tinyInfo, nil
	}
}

func GetListTinyByMap(m map[string]interface{}) (t []entity.TinyInfo, e error) {
	var (
		tinyInfo []entity.TinyInfo
		err      error
	)
	if err = mongo.DbFind(config.Base.Mongo.DbName, constants.TinyInfo, m, &tinyInfo); err != nil {
		return tinyInfo, err
	} else {
		return tinyInfo, nil
	}
}

func GetListTinyLimit(FData map[string]interface{}, skip, limit int) ([]entity.TinyInfo, error) {
	var (
		tinyList []entity.TinyInfo
		err      error
	)
	if err = mongo.DBFindPageSort(config.Base.Mongo.DbName,
		constants.TinyInfo, FData, skip, limit, "-createTime", &tinyList); err != nil {
		return nil, err
	}
	return tinyList, nil
}

func AddAccessCount(id string) error {
	var (
		err error
	)

	if err = mongo.DBUpdateById(config.Base.Mongo.DbName,
		constants.TinyInfo,
		id,
		mongo.B{"$inc": mongo.B{"count": 1}}); err != nil {
		return err
	}

	return err
}

func UpdateTinyInfo(q map[string]interface{}, set map[string]interface{}) (e error) {
	var (
		err error
	)
	if err = mongo.DbUpdateOne(config.Base.Mongo.DbName, constants.TinyInfo, q, set); err != nil {
		return err
	} else {
		return nil
	}
}
