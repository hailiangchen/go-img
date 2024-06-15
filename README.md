# go-img

#### 介绍
一个轻量级的图片服务器， 使用go语言翻译了zimg (https://github.com/buaazp/zimg)；
实现了和zimg一样的效果

#### 软件架构

1. gin go项目的一个web服务框架
2. redis 使用redis缓存
3. [imagick](gopkg.in/gographics/imagick.v3/imagick) go版本的处理图像工具，项目地址 gopkg.in/gographics/imagick.v3/imagick



#### 安装教程

#### 使用说明

上传接口： http://127.0.0.1:8080/upload

参数：Files 类型：文件

返回结果： [{"success":true,"message":"OK","version":"v0.1.1","data":{"size":49160,"mime":"image/jpeg","fileId":"5781339b809d5f18132f5c4fbe9df2fe","fileName":"gss0.baidu.jpg"}}]

1. 地址：http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?w=300&h=300&g=1&x=0&y=0&r=45&q=85&f=jpeg </br>
   格式组成： 服务器IP+端口/图片md5</br>
   w:宽，h:高，g:灰白化，x y:坐标点，r:旋转角度，q:压缩比，f:转换格式
2. 地址：http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?p=0  </br>
   请求原图使用p=0 默认：压缩质量为75%
3. 地址：http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?w=500&h=300&p=2  </br>
   p=2 按照目标分辨率提取图像中心部分
4. 地址：http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?w=500&h=300&p=3 </br>
   p=3 按照宽度或者高度提供的百分比，调整图像大小，参数范围1~100
5. 地址：http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?w=500&h=300&p=0 </br>
   p=0 按照提供的尺寸进行调整大小，图像会被拉伸
6. 地址：http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?w=500&h=300 </br>
   按照图像比例，调整图像大小
7. 地址：http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?w=500 </br>
   按照图像宽度或者高度进行等比例调整大小

