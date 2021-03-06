/*
@date : 2020/03/03
@author : YaPi
@desc :
*/
package business

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"tinyUrl/common/constants"
	"tinyUrl/common/http"
	"tinyUrl/common/util"
	"tinyUrl/common/util/snowflake"
	"tinyUrl/domain/dao/userdao"
	"tinyUrl/domain/dto"
	"tinyUrl/domain/entity"
)

/**
 * @api {post} /v1/api/register 注册
 * @apiUse Header
 * @apiVersion 0.0.1
 * @apiGroup userGroup
 * @apiPermission anyone
 * @apiParamExample {json} 请求示例:
 * {
	"userName":"admin",
	"password":"adminadmin"
}
 * @apiParam (body) {String} userName 	 用户名称.
 * @apiParam (body) {String} password 	 密码.
 * @apiUse FailResponse
 * @apiUse SuccessResponse
*/
func Register(ctx *gin.Context) {
	var (
		auth *dto.AuthDto
		err  error
		// 初始化返回结构体
		result = http.Instance()
		id     = strconv.Itoa(int(snowflake.NextId()))
	)
	// 请求参数校验
	if err = ctx.Bind(&auth); err != nil || auth.UserName == constants.EmptyStr ||
		auth.Password == constants.EmptyStr {
		result.Code = http.ParameterConvertError
		http.SendFailureRep(ctx, result)
		return
	}

	user := &entity.User{
		Id:         id,
		UserName:   auth.UserName,
		Password:   auth.Password,
		Status:     "standard",
		Role:       "admin",
		CreateTime: util.GetNowTimeStap(),
	}

	if err = userdao.AddUser(user); err != nil {
		result.Code = http.AddUserError
		http.SendFailureRep(ctx, result)
		return
	}
	http.SendSuccess(ctx)
}

/**
 * @api {post} /v1/api/login 登录
 * @apiUse Header
 * @apiVersion 0.0.1
 * @apiGroup userGroup
 * @apiPermission anyone
 * @apiParamExample {json} 请求示例:
 * {
	"userName":"admin",
	"password":"adminadmin"
}
 * @apiParam (body) {String} userName 	 用户名称.
 * @apiParam (body) {String} password 	 密码.
 * @apiUse FailResponse
 * @apiUse SuccessResponse
*/
func Login(ctx *gin.Context) {
	var (
		auth *dto.AuthDto
		err  error
		// 初始化返回结构体
		result = http.Instance()
	)
	// 请求参数校验
	if err = ctx.Bind(&auth); err != nil || auth.UserName == constants.EmptyStr ||
		auth.Password == constants.EmptyStr {
		result.Code = http.ParameterConvertError
		http.SendFailureRep(ctx, result)
		return
	}

	if user, err := userdao.GetByNameAndPassword(auth.UserName, auth.Password); err == nil && user != nil {
		if token, err := util.GenerateToken(auth.UserName, auth.Password); err == nil {
			result.Data = token
			http.SendSuccessRep(ctx, result)
		} else {
			result.Code = http.CreateTokenError
			http.SendFailureRep(ctx, result)
			return
		}
	} else {
		result.Code = http.UserNotFound
		http.SendFailureRep(ctx, result)
	}
}
