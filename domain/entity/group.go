/*
@date : 2020/03/16
@author : YaPi
@desc :
*/
package entity

type Group struct {
	Id         string `form:"id" json:"id" xml:"id" bson:"_id"`
	GroupName  string `form:"groupName" json:"groupName" xml:"groupName" bson:"groupName"`
	CreateTime int64  `form:"createTime" json:"createTime" xml:"createTime" bson:"createTime"`
	Status     uint8  `form:"status" json:"status" xml:"status"`
}
