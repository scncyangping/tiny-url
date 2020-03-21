/*
@date : 2020/03/16
@author : YaPi
@desc :
*/
package groupdao

import (
	"tinyUrl/domain/dao/groupdao/mongo"
	"tinyUrl/domain/entity"
)

// 添加组
func AddTinyGroup(group *entity.Group) error {
	return mongo.AddGroup(group)
}

func DeleteTinyGroup(groupId string) (e error) {
	q := map[string]interface{}{
		"id": groupId,
	}
	set := map[string]interface{}{
		"$set": map[string]interface{}{
			"status": 1,
		},
	}

	return mongo.UpdateGroup(q, set)
}

func GetGroupList() (t []entity.Group, e error) {
	m := make(map[string]interface{})
	m["status"] = 0
	return mongo.GetGroupByMap(m)
}
