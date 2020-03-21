/*
@Time : 2019-06-15 12:22
@Author : yangping
@File : tinyUrl
@Desc :
*/
package vo

type TinyVO struct {
	Id      string `form:"id" json:"id" xml:"id"`
	LongUrl string `form:"longUrl" json:"longUrl" xml:"longUrl"`
	TinyUrl string `form:"tinyUrl" json:"tinyUrl" xml:"tinyUrl"`
	// 计数
	Count       int    `form:"count" json:"count" xml:"count"`
	CreateTime  string `form:"createTime" json:"createTime" xml:"createTime"`
	ExpireTime  string `form:"expireTime" json:"expireTime" xml:"expireTime"`
	TinyUrlName string `form:"tinyUrlName" json:"tinyUrlName" xml:"tinyUrlName"`
}
