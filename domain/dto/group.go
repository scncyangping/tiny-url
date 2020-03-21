/*
@date : 2020/03/16
@author : YaPi
@desc :
*/
package dto

type GroupDto struct {
	Id        string `form:"id" json:"id" xml:"id" bson:"id"`
	GroupName string `form:"groupName" json:"groupName" xml:"groupName"`
}
