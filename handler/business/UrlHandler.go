/*
@Time : 2019-06-14 10:17
@Author : yangping
@File : UrlHandler
@Desc :
*/
package business

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"tinyUrl/common/constants"
	"tinyUrl/common/http"
	"tinyUrl/common/util"
	"tinyUrl/common/util/snowflake"
	"tinyUrl/config"
	"tinyUrl/config/db/mongo"
	"tinyUrl/config/db/redis"
	"tinyUrl/domain/dao/groupdao"
	"tinyUrl/domain/dao/tinyDao"
	"tinyUrl/domain/dto"
	"tinyUrl/domain/entity"
	"tinyUrl/domain/vo"
)

/**
* @api {get} /v1/api/tiny/url/info 对应短链基础信息
* @apiName 对应短链获取计数
* @apiUse Header
* @apiVersion 0.0.1
* @apiGroup urlGroup
* @apiPermission anyone
* @apiParamExample {http} 请求示例:
	http://localhost/v1/api/tiny/url/info?tinyId=123
* @apiParam  {String} tinyId 	 锻炼id.
* @apiSuccessExample {json} 返回示例:
	{
    "code": 200,
    "msg": "OK",
    "data": {
        "tinyId": "6651482131138191360",
        "groupId": "6651482343403528192",
        "longUrl": "http://baidu.com",
        "tinyUrl": "7VlReWQLV2o",
        "expireTime": 0,
        "tinyUrlName": "aaaa"
    }
}
* @apiUse FailResponse
* @apiUse SuccessResponse
*/
func UrlBaseInfo(ctx *gin.Context) {
	var (
		tinyDto dto.TinyDto
		// convert = util.NewBinaryConvert(config.Base.Convert.BinaryStr)
		err error
		// 初始化返回结构体
		result = http.Instance()
	)
	// 请求参数校验
	if err = ctx.Bind(&tinyDto); err != nil ||
		tinyDto.TinyId == constants.EmptyStr {
		result.Code = http.ParameterConvertError
		http.SendFailureRep(ctx, result)
		return
	}

	// 若 Redis 不存在此key, 查询DB内是否有对应key
	// urlId := strconv.Itoa(convert.AnyToDecimal(tinyDto.TinyUrl))
	if t, err := tinydao.GetTinyById(tinyDto.TinyId); err != nil {
		result.Code = http.QueryDBError
		http.SendFailureRep(ctx, result)
	} else {
		result.Data = &dto.TinyDto{
			TinyId:      t.Id,
			GroupId:     t.GroupId,
			LongUrl:     t.LongUrl,
			TinyUrl:     t.TinyUrl,
			TinyUrlName: t.TinyUrlName,
		}
		http.SendSuccessRep(ctx, result)
	}
}

/**
* @api {put} /v1/api/tiny/url/transform 更新短链接
* @apiUse Header
* @apiVersion 0.0.1
* @apiGroup urlGroup
* @apiPermission anyone
* @apiParamExample {json} 请求示例:
	{
		"longUrl" : "www.douyu.com",
		"expireTime" : 600,
		"tinyUrlName":"短链接名称2",
		"groupId":"6645249103407292416",
		"tinyId":"6645249103407292412"
	}
* @apiParam  {String} longUrl 	 	 原链接.
* @apiParam  {String} tinyUrlName 	 短链名称.
* @apiParam  {String} groupId 	 	 组ID.

* @apiParam  {Number} expireTime 	 过期时间(单位秒).
* @apiUse FailResponse
* @apiUse SuccessResponse
*/
func UpdateUrlTransform(ctx *gin.Context) {
	var (
		tinyInfo entity.TinyInfo
		tinyDto  dto.TinyDto
		err      error
		// 短连接
		tinyUrl string
		// 雪花算法生成ID
		id = int(snowflake.NextId())
		// 获取进制转换工具
		convert = util.NewBinaryConvert(config.Base.Convert.BinaryStr)
		// 初始化返回结构体
		result = http.Instance()
		//session = util.GetSession(ctx)
	)

	// 请求参数校验
	if err = ctx.Bind(&tinyDto); err != nil || tinyDto.LongUrl == constants.EmptyStr {
		result.Code = http.RequestParameterError
		http.SendFailureRep(ctx, result)
		return
	}

	// 查询此短链是否存在 存在直接返回 -- Redis
	//if isExist, _ := checkLongUrl(tinyDto.LongUrl, session.UserName); isExist {
	//	result.Data = "此长链已存在"
	//	http.SendSuccessRep(ctx, result)
	//	return
	//}
	// 相同长链对应多个短链, 若需要 1对1, 单独校重处理
	// 将ID转化为62进制
	tinyUrl = convert.DecimalToAny(id)
	tinyInfo.LongUrl = util.ConvertHttpUrl(tinyDto.LongUrl)
	tinyInfo.CreateTime = util.GetNowTimeStap()
	tinyInfo.ExpireTime = int64(tinyDto.ExpireTime*60) + tinyInfo.CreateTime
	tinyInfo.Count = constants.ZERO
	tinyInfo.UrlId = strconv.Itoa(id)
	tinyInfo.TinyUrl = tinyUrl
	tinyInfo.TinyUrlName = tinyDto.TinyUrlName

	if err = tinydao.UpdateTinyById(tinyDto.TinyId, &tinyInfo); err != nil {
		http.SendFailureError(ctx, result, err)
	} else {
		// 放在Redis中
		// 长链Redis中
		addLongUrlRedisKey(tinyInfo)

		// 短链放Redis中
		addTinyUrlRedisKey(tinyInfo)

		result.Data = map[string]string{
			"longUrl": tinyInfo.LongUrl,
			"tinyUrl": tinyInfo.TinyUrl,
		}
		http.SendSuccessRep(ctx, result)
	}
}

/**
* @api {post} /v1/api/tiny/url/transform 生成随机短链接
* @apiUse Header
* @apiVersion 0.0.1
* @apiGroup urlGroup
* @apiPermission anyone
* @apiParamExample {json} 请求示例:
	{
		"longUrl" : "www.douyu.com",
		"expireTime" : 600,
		"tinyUrlName":"短链接名称2",
		"groupId":"6645249103407292416"
	}
* @apiParam  {String} longUrl 	 	 原链接.
* @apiParam  {String} tinyUrlName 	 短链名称.
* @apiParam  {String} groupId 	 	 组ID.

* @apiParam  {Number} expireTime 	 过期时间(单位秒).
* @apiUse FailResponse
* @apiUse SuccessResponse
*/
func UrlTransform(ctx *gin.Context) {
	var (
		tinyInfo entity.TinyInfo
		tinyDto  dto.TinyDto
		err      error
		// 短连接
		tinyUrl string
		// 雪花算法生成ID
		id = int(snowflake.NextId())
		// 获取进制转换工具
		convert = util.NewBinaryConvert(config.Base.Convert.BinaryStr)
		// 初始化返回结构体
		result  = http.Instance()
		session = util.GetSession(ctx)
	)

	// 请求参数校验
	if err = ctx.Bind(&tinyDto); err != nil || tinyDto.LongUrl == constants.EmptyStr {
		result.Code = http.RequestParameterError
		http.SendFailureRep(ctx, result)
		return
	}

	// 查询此短链是否存在 存在直接返回 -- Redis
	if isExist, _ := checkLongUrl(tinyDto.LongUrl, session.UserName); isExist {
		result.Data = "此长链已存在"
		http.SendSuccessRep(ctx, result)
		return
	}
	// 相同长链对应多个短链, 若需要 1对1, 单独校重处理
	// 将ID转化为62进制
	tinyUrl = convert.DecimalToAny(id)
	tinyInfo.Id = strconv.Itoa(int(snowflake.NextId()))
	tinyInfo.LongUrl = util.ConvertHttpUrl(tinyDto.LongUrl)
	tinyInfo.UserName = session.UserName
	tinyInfo.CreateTime = util.GetNowTimeStap()
	tinyInfo.ExpireTime = int64(tinyDto.ExpireTime*60) + tinyInfo.CreateTime
	tinyInfo.Count = constants.ZERO
	tinyInfo.UrlId = strconv.Itoa(id)
	tinyInfo.TinyUrl = tinyUrl
	tinyInfo.Type = constants.ConvertDefault
	tinyInfo.Status = constants.ZERO
	tinyInfo.TinyUrlName = tinyDto.TinyUrlName
	tinyInfo.GroupId = tinyDto.GroupId
	if err = tinydao.AddTinyInfo(&tinyInfo); err != nil {
		http.SendFailureError(ctx, result, err)
	} else {
		// 放在Redis中
		// 长链Redis中
		addLongUrlRedisKey(tinyInfo)

		// 短链放Redis中
		addTinyUrlRedisKey(tinyInfo)

		result.Data = map[string]string{
			"longUrl": tinyInfo.LongUrl,
			"tinyUrl": tinyInfo.TinyUrl,
		}
		http.SendSuccessRep(ctx, result)
	}
}

/**
* @api {post} /v1/api/tiny/url/custom 自定义短链接
* @apiUse Header
* @apiVersion 0.0.1
* @apiGroup urlGroup
* @apiPermission anyone
* @apiParamExample {json} 请求示例:
	{
	"longUrl" : "www.baidu.com",
	"tinyUrl" : "test",
	"expireTime" : 60
	}
* @apiParam  {String} longUrl 	 	 原链接.
* @apiParam  {String} tinyUrl 	 	 希望生成的短链接.
* @apiParam  {Number} expireTime 	 过期时间(单位秒).
* @apiUse FailResponse
* @apiUse SuccessResponse
*/
func UrlTransformCustom(ctx *gin.Context) {
	var (
		tinyInfo entity.TinyInfo
		tinyDto  dto.TinyDto
		err      error
		// 初始化返回结构体
		result  = http.Instance()
		convert = util.NewBinaryConvert(config.Base.Convert.BinaryStr)
		session = util.GetSession(ctx)
	)

	// 请求参数校验
	if err = ctx.Bind(&tinyDto); err != nil ||
		tinyDto.LongUrl == constants.EmptyStr ||
		tinyDto.TinyUrl == constants.EmptyStr {
		result.Code = http.ParameterConvertError
		http.SendFailureRep(ctx, result)
		return
	}

	// 查询此长链是否存在 存在直接返回 -- Redis
	if isExist, _ := checkLongUrl(tinyDto.LongUrl, session.UserName); isExist {
		result.Data = "此长链已存在"
		http.SendSuccessRep(ctx, result)
		return
	}

	// 查询此短链是否存在 存在直接返回 -- Redis
	// 自定义短链接会校验DB
	if isExist, _ := checkTinyUrl(tinyDto.TinyUrl); isExist {
		result.Data = "此短链已存在"
		http.SendSuccessRep(ctx, result)
		return
	}

	// 不存在就新增
	tinyInfo.Id = strconv.Itoa(int(snowflake.NextId()))
	tinyInfo.UrlId = strconv.Itoa(convert.AnyToDecimal(tinyDto.TinyUrl))
	tinyInfo.UserName = session.UserName
	tinyInfo.CreateTime = util.GetNowTimeStap()
	tinyInfo.ExpireTime = int64(tinyDto.ExpireTime*60) + tinyInfo.CreateTime
	tinyInfo.LongUrl = util.ConvertHttpUrl(tinyDto.LongUrl)
	tinyInfo.Count = constants.ZERO
	tinyInfo.TinyUrl = tinyDto.TinyUrl
	tinyInfo.Type = constants.ConvertCustom
	tinyInfo.Status = constants.ZERO

	if err = tinydao.AddTinyInfo(&tinyInfo); err != nil {
		http.SendFailureError(ctx, result, err)
	} else {
		// 长链Redis中
		addLongUrlRedisKey(tinyInfo)

		// 短链放Redis中
		addTinyUrlRedisKey(tinyInfo)

		result.Data = &dto.TinyDto{
			LongUrl: tinyInfo.LongUrl,
			TinyUrl: tinyInfo.TinyUrl,
		}

		http.SendSuccessRep(ctx, result)
	}
}

/**
* @api {get} /v1/api/go 跳转对应链接
* @apiName 对应短链获取计数
* @apiUse Header
* @apiVersion 0.0.1
* @apiGroup urlGroup
* @apiPermission anyone
* @apiParamExample {http} 请求示例:
	http://localhost/v1/api/go?tinyUrl=Nqstssd
* @apiParam  {String} tinyUrl 	 短链名称.
* @apiUse FailResponse
* @apiUse SuccessResponse
*/
func Redirect4TinyUrl(ctx *gin.Context) {
	var (
		tinyDto dto.TinyDto
		err     error
		// 初始化返回结构体
		result = http.Instance()
	)
	// 请求参数校验
	if err = ctx.Bind(&tinyDto); err != nil ||
		tinyDto.TinyUrl == constants.EmptyStr {
		result.Code = http.ParameterConvertError
		http.SendFailureRep(ctx, result)
		return
	}
	// 查询此短链是否存在 存在直接返回 -- Redis
	// 自定义短链接会校验DB
	if isExist, longUrl := checkTinyUrl(tinyDto.TinyUrl); isExist && longUrl != constants.EmptyStr {
		// 若对应长链接存在,需要统计访问信息
		// 可以放在消息队列里面去做 便于更多样的统计
		// 这里直接单开线程 同步信息到DB中
		array := strings.Split(longUrl, constants.UnderLine)

		go tinydao.AddAccessCount(array[len(array)-1])

		ctx.Redirect(http.StatusFound, array[constants.ZERO])
		return
	}
	result.Data = "url not found"
	http.SendFailureRep(ctx, result)
}

/**
* @api {post} /v1/api/tiny/group 添加分组
* @apiUse Header
* @apiVersion 0.0.1
* @apiGroup urlGroup
* @apiPermission anyone
* @apiParamExample {json} 请求示例:
	{
		"groupName":"测试分组01"
	}
* @apiParam  {String} groupName 	 	 分组名称.
* @apiUse FailResponse
* @apiUse SuccessResponse
*/
func AddTinyGroup(ctx *gin.Context) {
	var (
		group    entity.Group
		groupDto dto.GroupDto
		err      error
		// 雪花算法生成ID
		id = int(snowflake.NextId())
		// 获取进制转换工具
		// 初始化返回结构体
		result = http.Instance()
		// session = util.GetSession(ctx)
	)

	// 请求参数校验
	if err = ctx.Bind(&groupDto); err != nil || groupDto.GroupName == constants.EmptyStr {
		result.Code = http.RequestParameterError
		http.SendFailureRep(ctx, result)
		return
	}

	group.Id = strconv.Itoa(int(id))
	group.GroupName = groupDto.GroupName
	group.CreateTime = util.GetNowTimeStap()
	group.Status = 0

	if err = groupdao.AddTinyGroup(&group); err != nil {
		http.SendFailureError(ctx, result, err)
	} else {
		result.Data = &dto.GroupDto{
			GroupName: groupDto.GroupName,
			Id:        group.Id,
		}
		http.SendSuccessRep(ctx, result)
	}
}

/**
* @api {post} /v1/api/tiny/list 短链列表
* @apiUse Header
* @apiVersion 0.0.1
* @apiGroup urlGroup
* @apiPermission anyone
* @apiParamExample {json} 请求示例:
{
	"groupId":"6645249103407292416",
	"page":1,
	"pageSize":1
}
* @apiParam  {String} groupId 	 	 分组Id.
* @apiSuccessExample {json} 返回示例:
	{
    "code": 200,
    "msg": "OK",
    "data": {
        "count": 2,
        "list": [
            {
                "id": "6647134219821666304",
                "longUrl": "http://www.douyu.com",
                "tinyUrl": "7V1WBwYQUnK",
                "count": 0,
                "createTime": "2020-03-21 22:18:15",
                "expireTime": "2020-03-21 22:28:15",
                "tinyUrlName": "短链接名称2"
            }
        ]
    }
}
* @apiUse FailResponse
*/
func TinyGroupList(ctx *gin.Context) {
	var (
		// 雪花算法生成ID
		// 获取进制转换工具
		// 初始化返回结构体
		result = http.Instance()
		fData  = http.QueryData(ctx)
		voList = make([]*vo.TinyVO, 0)
		// session = util.GetSession(ctx)
	)

	// 请求参数校验
	if fData.Data["groupId"] == nil {
		result.Code = http.RequestParameterError
		http.SendFailureRep(ctx, result)
		return
	}

	if resLis, err := tinydao.GetTinyByGroupId(fData.Data, fData.Skip, fData.Limit); err != nil {
		http.SendFailureError(ctx, result, err)
	} else {
		for _, v := range resLis {
			vo := &vo.TinyVO{
				Id:          v.Id,
				LongUrl:     v.LongUrl,
				TinyUrl:     v.TinyUrl,
				TinyUrlName: v.TinyUrlName,
				CreateTime:  util.GetTimeFormat(v.CreateTime),
				ExpireTime:  util.GetTimeFormat(v.ExpireTime),
			}
			voList = append(voList, vo)
		}
		fData.Data["status"] = 0
		countAll, err := mongo.GetCountByDb(config.Base.Mongo.DbName, constants.TinyInfo, fData.Data)

		if err != nil {
			result.Code = http.QueryDBError
			http.SendFailureRep(ctx, result)
			return
		}
		// 查询总数
		result.Data = map[string]interface{}{
			"list":  voList,
			"count": countAll,
		}
		http.SendSuccessRep(ctx, result)
	}
}

/**
* @api {get} /v1/api/tiny/group/list 查询分组列表
* @apiUse Header
* @apiVersion 0.0.1
* @apiGroup urlGroup
* @apiPermission anyone
* @apiParamExample {http} 请求示例:
		http://localhost:9069/v1/api/tiny/group/list
* @apiSuccessExample {json} 返回示例:
{
    "code": 200,
    "msg": "OK",
    "data": [
        {
            "id": "6645249103407292416",
            "groupName": "测试分组01"
        }
    ]
}
* @apiUse FailResponse
*/
func GroupList(ctx *gin.Context) {
	var (
		// 雪花算法生成ID
		// 获取进制转换工具
		// 初始化返回结构体
		result = http.Instance()
		voList = make([]*dto.GroupDto, 0)
		// session = util.GetSession(ctx)
	)
	// 请求参数校验
	if resLis, err := groupdao.GetGroupList(); err != nil {
		http.SendFailureError(ctx, result, err)
	} else {
		for _, v := range resLis {
			vo := &dto.GroupDto{
				Id:        v.Id,
				GroupName: v.GroupName,
			}
			voList = append(voList, vo)
		}
		result.Data = voList
		http.SendSuccessRep(ctx, result)
	}
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 校验Redis是否存在此长链接对应key 可设置是否校验DB
 */
func checkLongUrl(longUrl, userName string) (bool, string) {
	var (
		redisKey string
		str      string
		err      error
	)

	redisKey = fmt.Sprintf("%s:%s:%s:%s", constants.URL, constants.LongUrl, userName, longUrl)

	// 查询此短链是否存在 存在直接返回 -- Redis
	if str, err = getRedisKey(redisKey, false); err == nil {
		return true, str
	}
	return false, constants.EmptyStr
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 同一长链接在设置时间周期内,不允许重复生成短链接,防止攻击
 */
func addLongUrlRedisKey(tinyInfo entity.TinyInfo) (bool, error) {
	var (
		redisKey string
		str      string
		err      error
	)

	redisKey = fmt.Sprintf("%s:%s:%s:%s",
		constants.URL,
		constants.LongUrl,
		tinyInfo.UserName,
		tinyInfo.LongUrl)
	// 将这一条记录放在Redis当中

	str = tinyInfo.TinyUrl +
		constants.UnderLine +
		tinyInfo.UrlId +
		constants.UnderLine + tinyInfo.UserName
	err = redis.SetByTtl(redisKey, str, config.Base.Convert.LongUrlExpire)

	if err != nil {
		return false, err
	}
	return true, nil
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 添加短链缓存
 */
func addTinyUrlRedisKey(tinyInfo entity.TinyInfo) (bool, error) {
	var (
		redisKey string
		str      string
		err      error
	)

	redisKey = fmt.Sprintf("%s:%s:%s",
		constants.URL,
		constants.TinyUrl,
		tinyInfo.TinyUrl)
	// 将这一条记录放在Redis当中

	str = tinyInfo.LongUrl + constants.UnderLine + tinyInfo.UrlId

	now := util.GetNowTimeStap()
	// 若查出来的短链对应数据未过期
	if tinyInfo.ExpireTime > now {
		expireTime := int64(tinyInfo.ExpireTime) - now
		err = redis.SetByTtl(redisKey, str, expireTime)
		if err != nil {
			return false, err
		}
	} else {
		// 已过期
		// 将原数据进行删除
		go tinydao.DelteTinyByUrlId(tinyInfo.UrlId)
	}

	return true, nil
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 校验Redis是否存在此短链接对应key 可设置是否校验DB
 */
func checkTinyUrl(tinyUrl string) (bool, string) {
	var (
		convert  = util.NewBinaryConvert(config.Base.Convert.BinaryStr)
		redisKey string
	)

	redisKey = fmt.Sprintf("%s:%s:%s",
		constants.URL,
		constants.TinyUrl,
		tinyUrl)

	// 查询此短链是否存在 存在直接返回 -- Redis
	if str, err := getRedisKey(redisKey, false); err == nil {
		return true, str
	} else {
		if config.Base.Global.CheckDbOnRedisNotFound {
			// 若 Redis 不存在此key, 查询DB内是否有对应key
			urlId := strconv.Itoa(convert.AnyToDecimal(tinyUrl))
			t, error := tinydao.GetTinyByUrlId(urlId)
			// 将这一条记录放在Redis当中
			if error == nil {
				if _, err := addTinyUrlRedisKey(t); err != nil {
					return false, constants.EmptyStr

				} else {
					return true, constants.EmptyStr
				}
			}
		}
	}
	return false, constants.EmptyStr
}

/*
 * date : 2019-06-14
 * author : yangping
 * desc : 获取Redis的值,若存在则更新过期时间
 */
func getRedisKey(redisKey string, upExpire bool) (string, error) {
	var (
		tinyUrl string
	)

	tinyUrl = redis.Get(redisKey)

	if tinyUrl != constants.EmptyStr {

		if upExpire {
			if err := redis.Expire(redisKey, constants.ExpireTime); err != nil {
				return constants.EmptyStr, errors.New("update expire time error")
			}
		}
		return tinyUrl, nil

	} else {
		return constants.EmptyStr, errors.New("value does not exist")
	}
}
