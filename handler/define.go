/*
@date : 2020/03/04
@author : YaPi
@desc :
*/
package handler

/**
 * @apiDefine userGroup 用户模块
 */

/**
 * @apiDefine urlGroup 链接模块
 */

/**
 * @apiDefine Header
 * @apiHeader (header) {String} API_TOKEN 鉴权参数,登录获取.
 */

/**
 * @apiDefine FailResponse
 * @apiError (error) {Number} code 错误码
 * @apiError (error) {String} msg 错误内容
 * @apiErrorExample  {json} Response (fail):
 *     {
 *       "code":"error",
 *       "msg" : "错误内容",
 *		 "data" : null
 *     }
 */

/**
 * @apiDefine SuccessResponse
 * @apiSuccess (success) {Object} data 响应实体内容.
 * @apiSuccessExample  {json} Response (success):
 *     {
 *       "code":"200",
 *       "msg": "",
 *       "data":"响应实体内容"
 *     }
 */
