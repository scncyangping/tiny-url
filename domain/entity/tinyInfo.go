/*
@Time : 2019-06-14 16:18
@Author : yangping
@File : tinyInfo
@Desc :
*/
package entity

type TinyInfo struct {
	Id          string `form:"id" json:"id" xml:"id" bson:"_id"`
	UrlId       string `form:"urlId" json:"urlId" xml:"urlId" bson:"urlId"`
	UserName    string `form:"username" json:"id" username:"id" bson:"username"`
	LongUrl     string `form:"longUrl" json:"longUrl" xml:"longUrl" bson:"longUrl"`
	TinyUrl     string `form:"tinyUrl" json:"tinyUrl" xml:"tinyUrl" bson:"tinyUrl"`
	RedirectUrl string `form:"redirectUrl" json:"redirectUrl" xml:"redirectUrl" bson:"redirectUrl"`
	TinyUrlName string `form:"tinyUrlName" json:"tinyUrlName" xml:"tinyUrlName" bson:"tinyUrlName"`
	Count       int    `form:"count" json:"count" xml:"count"`
	Type        string `form:"type" json:"type" xml:"type" bson:"type"`
	CreateTime  int64  `form:"createTime" json:"createTime" xml:"createTime" bson:"createTime"`
	ExpireTime  int64  `form:"expireTime" json:"expireTime" xml:"expireTime" bson:"expireTime"`
	Status      uint8  `form:"status" json:"status" xml:"status"`
	GroupId     string `form:"groupId" json:"groupId" xml:"groupId" bson:"groupId"`
}
