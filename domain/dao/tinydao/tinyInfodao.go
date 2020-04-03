/*
@Time : 2019-06-14 17:28
@Author : yangping
@File : tinyInfoDao
@Desc :
*/
package tinydao

import (
	"tinyUrl/domain/dao/tinydao/mongo"
	"tinyUrl/domain/entity"
)

func AddTinyInfo(tiny *entity.TinyInfo) (err error) {
	return mongo.AddTinyInfo(tiny)
}

func UpdateTinyById(id string, tiny *entity.TinyInfo) error {
	q := map[string]interface{}{
		"_id": id,
	}

	set := map[string]interface{}{
		"$set": map[string]interface{}{
			"longUrl":     tiny.LongUrl,
			"createTime":  tiny.CreateTime,
			"expireTime":  tiny.ExpireTime,
			"count":       tiny.Count,
			"urlId":       tiny.UrlId,
			"tinyUrl":     tiny.TinyUrl,
			"tinyUrlName": tiny.TinyUrlName,
		},
	}

	return mongo.UpdateTinyInfo(q, set)
}

func GetTinyInfoById(id string) (t entity.TinyInfo, e error) {

	return mongo.GetTinyInfoById(id)
}

func GetTinyById(id string) (t entity.TinyInfo, e error) {
	m := make(map[string]interface{})
	m["_id"] = id
	m["status"] = 0
	return mongo.GetTinyByMap(m)
}

func GetTinyByUrlId(urlId string) (t entity.TinyInfo, e error) {
	m := make(map[string]interface{})
	m["urlId"] = urlId
	m["status"] = 0
	return mongo.GetTinyByMap(m)
}

func GetTinyByGroupId(fData map[string]interface{}, skip, limit int) (t []entity.TinyInfo, e error) {
	fData["status"] = 0
	// return mongo.GetListTinyByMap(m)
	return mongo.GetListTinyLimit(fData, skip, limit)
}

func DelteTinyByUrlId(urlId string) (e error) {
	q := map[string]interface{}{
		"urlId": urlId,
	}
	set := map[string]interface{}{
		"$set": map[string]interface{}{
			"status": 1,
		},
	}

	return mongo.UpdateTinyInfo(q, set)
}

func AddAccessCount(id string) error {
	return mongo.AddAccessCount(id)
}
