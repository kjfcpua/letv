package sdk

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

//乐视云SDK

const (
	ALL     int = 0
	PLAY_OK int = 10
	FAILED  int = 20
	WAIT    int = 30
)

var (
	//	post timeout时间
	PTimeOut time.Duration
	//	get timeout时间
	GTimeOut time.Duration
)

type LetvCloudV1 struct {
	userUnique string
	secretKey  string
	restUrl    string
	format     string
	apiVersion string
}

func NewLetvCloudV1(unique, key string) *LetvCloudV1 {
	return &LetvCloudV1{userUnique: unique, secretKey: key, restUrl: "http://api.letvcloud.com/open.php", format: "json", apiVersion: "2.0"}
}

/**
 * 视频上传初始化
 * @param  string video_name 视频名称
 * @param  string client_ip  用户IP地址
 * @return []byte
 */
func (this *LetvCloudV1) VideoUploadInit(video_name string) []byte {
	return this.videoUploadInit_(video_name, "", 0)
}

func (this *LetvCloudV1) SetSecretKey(secretKey string) {
	this.secretKey = secretKey
}
func (this *LetvCloudV1) SetRestUrl(restUrl string) {
	this.restUrl = restUrl
}
func (this *LetvCloudV1) SetFormat(format string) {
	this.format = format
}
func (this *LetvCloudV1) SetApiVersion(apiVersion string) {
	this.apiVersion = apiVersion
}

/**
 * 视频上传初始化
 * @param  string video_name 视频名称
 * @param  string client_ip  用户IP地址
 * @param  int file_size  文件大小，单位为字节
 * @return []byte
 */
func (this *LetvCloudV1) videoUploadInit_(video_name string, client_ip string, file_size int) []byte {
	api := "video.upload.init"
	params := make(map[interface{}]interface{})

	params["video_name"] = video_name
	if len(client_ip) > 0 {
		params["client_ip"] = client_ip
	}
	if file_size > 0 {

		params["file_size"] = strconv.Itoa(file_size)
	}
	return this.makeRequest(api, params)
}

/**
 * 视频上传 (web方式)
 * @param  string video_file 文件绝对路径
 * @param  string upload_url 视频上传地址，视频上传时提交地址
 * @return []byte
 */
func (this *LetvCloudV1) VideoUpload(video_file, upload_url string) []byte {
	return this.doUploadFile(video_file, upload_url)
}

/**
 * 视频上传（Flash方式）
 * @param  string video_name 视频名称
 * @param  string js_callback Javascript回调函数，视频上传完毕后调用
 * @param  int flash_width Flash宽度，默认值为600
 * @param  int flash_height Flash高度，默认值为450
 * @param  string client_ip 用户IP地址
 * @return  []byte
 */
func (this *LetvCloudV1) videoUploadFlash(video_name, js_callback string, flash_width, flash_height int, client_ip string) []byte {
	api := "video.upload.flash"
	params := make(map[interface{}]interface{})
	params["video_name"] = video_name
	if len(js_callback) > 0 {
		params["js_callback"] = js_callback
	}
	if flash_width > 0 {
		params["flash_width"] = strconv.Itoa(flash_width)
	}
	if flash_height > 0 {
		params["flash_height"] = strconv.Itoa(flash_height)
	}
	if len(client_ip) > 0 {
		params["client_ip"] = client_ip
	}
	return this.makeRequest(api, params)
}

/**
 * 视频上传（Flash方式）
 * @param  string video_name 视频名称
 * @param  string js_callback Javascript回调函数，视频上传完毕后调用
 * @param  int flash_width Flash宽度，默认值为600
 * @param  int flash_height Flash高度，默认值为450
 * @return []byte
 */
func (this *LetvCloudV1) videoUploadFlash_(video_name, js_callback string, flash_width, flash_height int) []byte {
	return this.videoUploadFlash(video_name, js_callback, flash_width, flash_height, "")
}

/**
 * 视频上传（Flash方式）
 * @param  string video_name 视频名称
 * @param  string js_callback Javascript回调函数，视频上传完毕后调用
 * @param  int flash_width Flash宽度，默认值为600
 * @return []byte
 */
func (this *LetvCloudV1) videoUploadFlash_1(video_name, js_callback string, flash_width int) []byte {
	return this.videoUploadFlash(video_name, js_callback, flash_width, 0, "")
}

/**
 * 视频上传（Flash方式）
 * @param  string video_name 视频名称
 * @param  string js_callback Javascript回调函数，视频上传完毕后调用
 * @return []byte
 */
func (this *LetvCloudV1) videoUploadFlash_2(video_name, js_callback string) []byte {
	return this.videoUploadFlash(video_name, js_callback, 0, 0, "")
}

/**
 * 视频上传（Flash方式）
 * @param  string video_name 视频名称
 * @return []byte
 */
func (this *LetvCloudV1) videoUploadFlash_3(video_name string) []byte {
	return this.videoUploadFlash(video_name, "", 0, 0, "")
}

/**
 * 视频断点续传
 * @param  string token 视频上传标识
 * @return String
 */
func (this *LetvCloudV1) VideoUploadResume(token string) string {
	api := "video.upload.resume"
	params := make(map[interface{}]interface{})
	params["token"] = token
	return string(this.makeRequest(api, params))
}

/**
 * 构造云视频Sign
 * @param params 业务参数
 * @return string
 */
func (this *LetvCloudV1) generateSign(params map[interface{}]interface{}) string {
	array := make([]string, 0)

	for key := range params {
		array = append(array, key.(string))
	}
	sort.Strings(array)
	keyStr := ""
	for _, v := range array {

		keyStr = keyStr + v + params[v].(string)
	}
	keyStr += this.secretKey
	//	fmt.Println("md5:", keyStr)
	return this.md5_(keyStr)
}

/**
 * 视频信息更新
 * @param  int video_id 视频ID
 * @param  string video_name 视频名称
 * @param  string video_desc 视频简介
 * @param  string tag 标签
 * @param  int is_pay 视频是否收费：0表示不收费；1表示收费（收费视频播放时会进行用户鉴权，请不要随便设置）
 * @return []byte
 */
func (this *LetvCloudV1) videoUpdate(video_id int, video_name, video_desc, tag string, is_pay int) []byte {
	api := "video.upload.init"
	params := make(map[interface{}]interface{})

	params["video_id"] = strconv.Itoa(video_id)

	if len(video_name) > 0 {
		params["video_name"] = video_name
	}
	if len(video_desc) > 0 {
		params["video_desc"] = video_desc
	}
	if len(tag) > 0 {
		params["tag"] = tag
	}

	if is_pay == 0 || is_pay == 1 {
		params["is_pay"] = strconv.Itoa(is_pay)
	}

	return this.makeRequest(api, params)
}

/**
 * 视频信息更新
 * @param  int video_id 视频ID
 * @param  string video_name 视频名称
 * @param  string video_desc 视频简介
 * @return []byte
 */
func (this *LetvCloudV1) videoUpdate_1(video_id int, video_name, video_desc string) []byte {
	return this.videoUpdate(video_id, video_name, video_desc, "", -1)
}

/**
 * 视频信息更新
 * @param  int video_id 视频ID
 * @return []byte
 */
func (this *LetvCloudV1) videoUpdate_2(video_id int) []byte {
	return this.videoUpdate(video_id, "", "", "", -1)
}

/**
 * 视频信息更新
 * @param  int video_id 视频ID
 * @param  string video_name 视频名称
 * @param  string video_desc 视频简介
 * @param  string tag 标签
 * @return []byte
 */
func (this *LetvCloudV1) videoUpdate_(video_id int, video_name, video_desc, tag string) []byte {
	return this.videoUpdate(video_id, video_name, video_desc, tag, -1)
}

/**
 * 获取视频列表
 * @param  int index 开始页索引，默认值为1
 * @param  int size 分页大小，默认值为10，最大值为100
 * @param  const status 视频状态：ALL表示全部；PLAY_OK表示可以正常播放；FAILED表示处理失败；WAIT表示正在处理过程中。默认值为ALL
 * @return []byte
 */
func (this *LetvCloudV1) videoList(index, size, status int) []byte {
	api := "video.list"
	params := make(map[interface{}]interface{})
	if index > 0 {
		params["index"] = strconv.Itoa(index)
	}
	if size > 0 {
		params["size"] = strconv.Itoa(size)
	}
	if status == ALL || status == PLAY_OK || status == FAILED || status == WAIT {
		params["status"] = strconv.Itoa(status)
	}
	return this.makeRequest(api, params)
}

/**
 * 获取视频列表
 * @param  int index 开始页索引，默认值为1
 * @param  int size 分页大小，默认值为10，最大值为100
 * @return []byte
 */
func (this *LetvCloudV1) videoList_(index, size int) []byte {
	return this.videoList(index, size, -1)
}

/**
 * 获取视频列表
 * @param  int index 开始页索引，默认值为1
 * @param  int size 分页大小，默认值为10，最大值为100
 * @return []byte
 */
func (this *LetvCloudV1) videoList_1(index int) []byte {
	return this.videoList(index, 0, -1)
}

/**
 * 获取视频列表
 * @param  int index 开始页索引，默认值为1
 * @param  int size 分页大小，默认值为10，最大值为100
 * @return []byte
 */
func (this *LetvCloudV1) videoList_2() []byte {
	return this.videoList(0, 0, -1)
}

/**
 * 获取单个视频信息
 * @param videoid 视频id
 * @return []byte
 */
func (this *LetvCloudV1) videoGet(videoid int) []byte {
	api := "video.get"
	params := make(map[interface{}]interface{})
	params["video_id"] = strconv.Itoa(videoid)
	return this.makeRequest(api, params)
}

/**
 * 删除视频
 * @param  int video_id 视频ID
 * @return []byte
 */
func (this *LetvCloudV1) videoDel(video_id int) []byte {
	api := "video.del"
	params := make(map[interface{}]interface{})
	params["video_id"] = strconv.Itoa(video_id)
	return this.makeRequest(api, params)
}

/**
 * 批量删除视频
 * @param  string video_id_list 视频ID列表，使用符号-作为间隔符，每次最多操作50条记录
 * @return String
 */
func (this *LetvCloudV1) videoDelBatch(video_id_list string) []byte {
	api := "video.del.batch"
	params := make(map[interface{}]interface{})
	params["video_id_list"] = video_id_list
	return this.makeRequest(api, params)
}

/**
 * 视频暂停
 * @param  int video_id 视频ID
 * @return String
 */
func (this *LetvCloudV1) videoPause(video_id int) []byte {
	api := "video.pause"
	params := make(map[interface{}]interface{})
	params["video_id"] = strconv.Itoa(video_id)
	return this.makeRequest(api, params)
}

/**
 * 视频恢复
 * @param  int video_id 视频ID
 * @return String
 */
func (this *LetvCloudV1) videoRestore(video_id int) []byte {
	api := "video.restore"
	params := make(map[interface{}]interface{})
	params["video_id"] = strconv.Itoa(video_id)
	return this.makeRequest(api, params)
}

/**
 * 获取视频截图
 * @param  int video_id 视频ID
 * @param  string size 截图尺寸，每种尺寸各有8张图。
 * @return []byte
 */
func (this *LetvCloudV1) imageGet(video_id int, size string) []byte {
	api := "image.get"
	params := make(map[interface{}]interface{})
	params["video_id"] = strconv.Itoa(video_id)
	params["size"] = size
	return this.makeRequest(api, params)
}

/**
 * 视频小时数据
 * @param  string date 日期，格式为：yyyy-mm-dd
 * @param  int hour 小时，0-23之间
 * @param  int video_id 视频ID
 * @param  int index 开始页索引，默认值为1
 * @param  int size 分页大小，默认值为10，最大值为100
 * @return []byte
 */
func (this *LetvCloudV1) dataVideoHour(date string, hour, video_id, index, size int) []byte {
	api := "data.video.hour"
	params := make(map[interface{}]interface{})
	params["date"] = date
	if hour >= 0 && hour <= 23 {
		params["hour"] = strconv.Itoa(hour)
	}
	if video_id > 0 {
		params["video_id"] = strconv.Itoa(video_id)
	}
	if index > 0 {
		params["index"] = strconv.Itoa(index)
	}
	if size > 0 {
		params["size"] = strconv.Itoa(size)
	}
	return this.makeRequest(api, params)
}

/**
 * 视频小时数据
 * @param  string date 日期，格式为：yyyy-mm-dd
 * @param  int hour 小时，0-23之间
 * @param  int video_id 视频ID
 * @param  int index 开始页索引，默认值为1
 * @return []byte
 */
func (this *LetvCloudV1) dataVideoHour_(date string, hour, video_id, index int) []byte {
	return this.dataVideoHour(date, hour, video_id, index, 0)
}

/**
 * 视频小时数据
 * @param  string date 日期，格式为：yyyy-mm-dd
 * @param  int hour 小时，0-23之间
 * @param  int video_id 视频ID
 * @param  int index 开始页索引，默认值为1
 * @return []byte
 */
func (this *LetvCloudV1) dataVideoHour_1(date string, hour, video_id int) []byte {
	return this.dataVideoHour(date, hour, video_id, 0, 0)
}

/**
 * 视频小时数据
 * @param  string date 日期，格式为：yyyy-mm-dd
 * @param  int hour 小时，0-23之间
 * @param  int video_id 视频ID
 * @param  int index 开始页索引，默认值为1
 * @return []byte
 */
func (this *LetvCloudV1) dataVideoHour_2(date string, hour int) []byte {
	return this.dataVideoHour(date, hour, 0, 0, 0)
}

/**
 * 视频小时数据
 * @param  string date 日期，格式为：yyyy-mm-dd
 * @param  int hour 小时，0-23之间
 * @param  int video_id 视频ID
 * @param  int index 开始页索引，默认值为1
 * @return []byte
 */
func (this *LetvCloudV1) dataVideoHour_3(date string) []byte {
	return this.dataVideoHour(date, -1, 0, 0, 0)
}

/**
 * 视频天数据
 * @param  string start_date 开始日期，格式为：yyyy-mm-dd
 * @param  string end_date 结束日期，格式为：yyyy-mm-dd
 * @param  int video_id 视频ID，不输入该参数将返回所有视频的数据
 * @param  int index 开始页索引，默认值为1
 * @param  int size 分页大小，默认值为10，最大值为100
 * @return []byte
 */
func (this *LetvCloudV1) dataVideoDate(start_date, end_date string, video_id, index, size int) []byte {
	api := "data.video.date"
	params := make(map[interface{}]interface{})
	params["start_date"] = start_date
	params["end_date"] = end_date
	if video_id > 0 {
		params["video_id"] = strconv.Itoa(video_id)
	}
	if index > 0 {
		params["index"] = strconv.Itoa(index)
	}
	if size > 0 {
		params["size"] = strconv.Itoa(size)
	}
	return this.makeRequest(api, params)
}

/**
 * 视频天数据
 * @param  string start_date 开始日期，格式为：yyyy-mm-dd
 * @param  string end_date 结束日期，格式为：yyyy-mm-dd
 * @param  int video_id 视频ID，不输入该参数将返回所有视频的数据
 * @param  int index 开始页索引，默认值为1
 * @return []byte
 */
func (this *LetvCloudV1) dataVideoDate_(start_date, end_date string, video_id, index int) []byte {
	return this.dataVideoDate(start_date, end_date, video_id, index, 0)
}

/**
 * 视频天数据
 * @param  string start_date 开始日期，格式为：yyyy-mm-dd
 * @param  string end_date 结束日期，格式为：yyyy-mm-dd
 * @param  int video_id 视频ID，不输入该参数将返回所有视频的数据
 * @param  int index 开始页索引，默认值为1
 * @return []byte
 */
func (this *LetvCloudV1) dataVideoDate_1(start_date, end_date string, video_id int) []byte {
	return this.dataVideoDate(start_date, end_date, video_id, 0, 0)
}

/**
 * 视频天数据
 * @param  string start_date 开始日期，格式为：yyyy-mm-dd
 * @param  string end_date 结束日期，格式为：yyyy-mm-dd
 * @param  int video_id 视频ID，不输入该参数将返回所有视频的数据
 * @param  int index 开始页索引，默认值为1
 * @return []byte
 */
func (this *LetvCloudV1) dataVideoDate_2(start_date, end_date string) []byte {
	return this.dataVideoDate(start_date, end_date, 0, 0, 0)
}

/**
 * 所有数据
 * @param  string start_date 开始日期，格式为：yyyy-mm-dd
 * @param  string end_date 结束日期，格式为：yyyy-mm-dd
 * @param  int index 开始页索引，默认值为1
 * @param  int size 分页大小，默认值为10，最大值为100
 * @return []byte
 */
func (this *LetvCloudV1) dataTotalDate(start_date, end_date string, index, size int) []byte {
	api := "data.total.date"
	params := make(map[interface{}]interface{})
	params["start_date"] = start_date
	params["end_date"] = end_date
	if index > 0 {
		params["index"] = strconv.Itoa(index)
	}
	if size > 0 {
		params["size"] = strconv.Itoa(size)
	}
	return this.makeRequest(api, params)
}

/**
 * 所有数据
 * @param  string start_date 开始日期，格式为：yyyy-mm-dd
 * @param  string end_date 结束日期，格式为：yyyy-mm-dd
 * @param  int index 开始页索引，默认值为1
 * @param  int size 分页大小，默认值为10，最大值为100
 * @return []byte
 */
func (this *LetvCloudV1) dataTotalDate_(start_date, end_date string, index int) []byte {
	return this.dataTotalDate(start_date, end_date, index, 0)
}

/**
 * 所有数据
 * @param  string start_date 开始日期，格式为：yyyy-mm-dd
 * @param  string end_date 结束日期，格式为：yyyy-mm-dd
 * @param  int index 开始页索引，默认值为1
 * @param  int size 分页大小，默认值为10，最大值为100
 * @return []byte
 */
func (this *LetvCloudV1) dataTotalDate_1(start_date, end_date string) []byte {
	return this.dataTotalDate(start_date, end_date, 0, 0)
}

/**
 * 获取视频播放接口
 * @param string uu 用户唯一标识码，由乐视网统一分配并提供
 * @param string vu 视频唯一标识码
 * @param string type 接口类型：url表示播放URL地址；js表示JavaScript代码；flash表示视频地址；html表示HTML代码
 * @param string pu 播放器唯一标识码
 * @param int auto_play 是否自动播放：1表示自动播放；0表示不自动播放。默认值由双方事先约定
 * @param int width 播放器宽度
 * @param int height 播放器高度
 * @return String
 */
func (this *LetvCloudV1) videoGetPlayinterface(uu, vu, types, pu string, auto_play, width, height int) string {
	params := make(map[interface{}]interface{})
	params["uu"] = uu
	params["vu"] = vu
	if len(pu) > 0 {
		params["pu"] = pu
	}
	if auto_play != -1 {
		params["auto_play"] = strconv.Itoa(auto_play)
	}
	if width > 0 {
		params["width"] = strconv.Itoa(width)
	} else {
		width = 800
	}
	if height > 0 {
		params["height"] = strconv.Itoa(height)
	} else {
		height = 450
	}
	queryString := this.mapToQueryString(params)
	jsonString := this.mapToJsonString(params)
	response := ""
	if types == "url" {
		response = "http://yuntv.letv.com/bcloud.html?" + queryString
	}
	if types == "js" {
		response = "<script type=\"text/javascript\">var letvcloud_player_conf = " + jsonString + ";</script><script type=\"text/javascript\" src=\"http://yuntv.letv.com/bcloud.js\"></script>"
	}
	if types == "flash" {
		response = "http://yuntv.letv.com/bcloud.swf?" + queryString
	}
	if types == "html" {

		response = "<embed src=\"http://yuntv.letv.com/bcloud.swf\" allowFullScreen=\"true\" quality=\"high\" width=\"" + strconv.Itoa(width) + "\" height=\"" + strconv.Itoa(height) + "\" align=\"middle\" allowScriptAccess=\"always\" flashvars=\"" + queryString + "\" type=\"application/x-shockwave-flash\"></embed>"
	}
	return response
}

/**
 * 获取视频播放接口
 * @param string uu 用户唯一标识码，由乐视网统一分配并提供
 * @param string vu 视频唯一标识码
 * @param string type 接口类型：url表示播放URL地址；js表示JavaScript代码；flash表示视频地址；html表示HTML代码
 * @param string pu 播放器唯一标识码
 * @param int auto_play 是否自动播放：1表示自动播放；0表示不自动播放。默认值由双方事先约定
 * @param int width 播放器宽度
 * @return String
 */
func (this *LetvCloudV1) videoGetPlayinterface_(uu, vu, types, pu string, auto_play, width int) string {
	return this.videoGetPlayinterface(uu, vu, types, pu, auto_play, width, 0)
}

/**
 * 获取视频播放接口
 * @param string uu 用户唯一标识码，由乐视网统一分配并提供
 * @param string vu 视频唯一标识码
 * @param string type 接口类型：url表示播放URL地址；js表示JavaScript代码；flash表示视频地址；html表示HTML代码
 * @param string pu 播放器唯一标识码
 * @param int auto_play 是否自动播放：1表示自动播放；0表示不自动播放。默认值由双方事先约定
 * @param int width 播放器宽度
 * @return String
 */
func (this *LetvCloudV1) videoGetPlayinterface_1(uu, vu, types, pu string, auto_play int) string {
	return this.videoGetPlayinterface(uu, vu, types, pu, auto_play, 0, 0)
}

/**
 * 获取视频播放接口
 * @param string uu 用户唯一标识码，由乐视网统一分配并提供
 * @param string vu 视频唯一标识码
 * @param string type 接口类型：url表示播放URL地址；js表示JavaScript代码；flash表示视频地址；html表示HTML代码
 * @param string pu 播放器唯一标识码
 * @param int auto_play 是否自动播放：1表示自动播放；0表示不自动播放。默认值由双方事先约定
 * @param int width 播放器宽度
 * @return String
 */
func (this *LetvCloudV1) videoGetPlayinterface_2(uu, vu, types, pu string) string {
	return this.videoGetPlayinterface(uu, vu, types, pu, -1, 0, 0)
}

/**
 * 获取视频播放接口
 * @param string uu 用户唯一标识码，由乐视网统一分配并提供
 * @param string vu 视频唯一标识码
 * @param string type 接口类型：url表示播放URL地址；js表示JavaScript代码；flash表示视频地址；html表示HTML代码
 * @param string pu 播放器唯一标识码
 * @param int auto_play 是否自动播放：1表示自动播放；0表示不自动播放。默认值由双方事先约定
 * @param int width 播放器宽度
 * @return String
 */
func (this *LetvCloudV1) videoGetPlayinterface_3(uu, vu, types string) string {
	return this.videoGetPlayinterface(uu, vu, types, "", -1, 0, 0)
}

/**
 * 将 int64转换为string
 * @param i int64类型
 * @return s string
 */
func Int64Tstr(i int64) string {
	s := strconv.FormatInt(i, 10)
	return s
}

//构造请求串
func (this *LetvCloudV1) makeRequest(api string, params map[interface{}]interface{}) []byte {
	params["user_unique"] = this.userUnique
	//微秒
	time := time.Now().UnixNano() / 1000000
	params["timestamp"] = Int64Tstr(time) //毫秒时间
	params["ver"] = this.apiVersion
	params["format"] = this.format
	params["api"] = api
	params["sign"] = this.generateSign(params)
	//	params["uploadtype"] = "1"
	//	params["isdownload"] = "1"

	resurl := ""
	resurl += this.restUrl + "?" + this.mapToQueryString(params)

	return doGet(resurl)
}

//将 map 中的参数及对应值转换为查询字符串
func (this *LetvCloudV1) mapToQueryString(params map[interface{}]interface{}) string {

	str := ""
	v := url.Values{}
	for key, value := range params {
		v.Add(key.(string), value.(string))
	}
	str = v.Encode()

	//url空格需要替换成%20,传输才能成功
	if strings.Contains(str, "+") {
		str = strings.Replace(str, " ", "%20", -1)
	}

	return str
}

func (this *LetvCloudV1) mapToJsonString(params map[interface{}]interface{}) string {

	if bs, err := json.Marshal(params); err != nil {
		fmt.Println(err)
		return ""
	} else {
		return string(bs)
	}

}

//GET请求
func doGet(url string) []byte {
	c := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(GTimeOut * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*GTimeOut)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	var body []byte
	resp, err := c.Get(url)

	if err != nil {
		fmt.Println("error get connection is fail")
		return nil
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error get connection ioutil.ReadAll")
		return nil
	}

	return body
}

//POST上传文件

func (this *LetvCloudV1) doUploadFile(filename, targetUrl string) []byte {

	c := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(PTimeOut * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*PTimeOut)
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		fmt.Println("error io.Copy")
		return nil
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := c.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		fmt.Println("error post connection is fail")
		return nil
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(resp_body)
		return nil
	}

	return resp_body
}

//MD5加密
func (this *LetvCloudV1) md5_(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
