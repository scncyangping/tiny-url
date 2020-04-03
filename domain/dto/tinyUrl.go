/*
@Time : 2019-06-14 17:10
@Author : yangping
@File : tinyUrl
@Desc :
*/
package dto

type TinyDto struct {
	TinyId      string `form:"tinyId" json:"tinyId" xml:"tinyId" bson:"tinyId"`
	GroupId     string `form:"groupId" json:"groupId" xml:"groupId" bson:"groupId"`
	LongUrl     string `form:"longUrl" json:"longUrl" xml:"longUrl"`
	TinyUrl     string `form:"tinyUrl" json:"tinyUrl" xml:"tinyUrl"`
	ExpireTime  int    `form:"expireTime" json:"expireTime" xml:"expireTime"`
	TinyUrlName string `form:"tinyUrlName" json:"tinyUrlName" xml:"tinyUrlName" bson:"tinyUrlName"`
}
