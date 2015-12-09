# letv-cloub-sdk

乐视云点播SDK

视频上传（Web方式）

视频上传（Web方式）
接口地址：由video.upload.init返回的upload_url确定
功能描述：视频上传
请求参数：由video.upload.init返回的upload_url确定，请不要进行任何修改
返回参数data说明：

名称	类型	描述
code	int	状态值：0表示操作成功；其它值表示失败，具体含义见message说明
message	string	状态说明

注意事项：
1、该接口只支持POST参数传递方式
2、上传时form表单提交地址由video.uploadinit返回的upload_url确定，需要添加参数enctype="multipart/form-data"，form表单中上传文件的参数名称为video_file
示例如下：
选择文件
提交

3、支持以下视频格式，请上传时做一下后缀名校验：
微软视频：.wmv .avi .dat .asf
Real Player：.rm .rmvb .ram
MPEG视频：.mpg .mpeg
其他常见视频： .mp4 .mov. m4v .mkv .flv .vob .qt .divx .cpk .fli .flc .mod .dvix .dv .ts

视频上传初始化（Web方式）

api名称：video.upload.init

功能描述：视频上传前调用，获取正式上传时需要的一些信息

应用参数说明：


名称	类型	必选	描述
video_name	string(200)	Y	视频名称
client_ip	string(15)	N	用户IP地址。为了保证用户上传速度，建议将用户公网IP地址写入该参数
file_size	int	N	文件大小，单位为字节
uploadtype	int	N	是否分片上传，0不分片，1分片；默认0。如用js编写上传功能须使用分片模式。
isdownload	int	N	是否支持缓存上传标志,默认0，0 不支持缓存 1支持缓存。（注：离线下载为移动端功能）
isdrm	int	N	是否支持DRM上传标志,默认0，不支持drm。（注：html5 不支持播放加密视频，不加密可进行播放，加密之后不支持离线下载）
ispay	int	N	是否付费标志,默认0，不付费。（注：需要客户配置回调地址。）回调地址说明

返回参数data说明：


名称	类型	描述
video_id	int	视频ID
video_unique	string	视频唯一标识码
upload_url	string	视频上传地址，视频上传时提交地址
progress_url	string	视频上传进度查询地址
token	string	视频上传标识，用于断点续传和上传进度查询
uploadtype	int	是否分片上传，0不分片，1分片 。如用js编写上传功能须使用分片模式。

视频断点续传（Web方式）

视频断点续传（Web方式）
api名称：video.upload.resume

功能描述：视频文件断点续传

应用参数说明：


名称	类型	必选	描述
token	string(500)	Y	视频上传标识
uploadtype	int	N	是否分片上传，0不分片，1分片；默认0。如用js编写上传功能须使用分片模式。

返回参数data说明：


名称	类型	描述
upload_url	string	视频续传地址，视频续传时提交地址
progress_url	string	视频上传进度查询地址
upload_size	int	已经上传的文件大小，单位为字节
uploadtype	int	是否分片上传，0不分片，1分片 。如用js编写上传功能须使用分片模式。

视频上传进度查询（Web方式）

视频上传进度查询（Web方式）
接口地址：由video.upload.init或video.upload.resume返回的progress_url确定

功能描述：视频上传进度查询，只有文件在上传过程中调用才有意义

请求参数：由video.upload.init或video.upload.resume返回的progress_url确定，请不要进行任何修改

返回参数data说明：


名称	类型	描述
total_size	int	视频文件总大小，单位为字节
upload_size	int	已经上传的数据大小，单位为字节
