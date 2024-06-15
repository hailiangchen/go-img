# goimg

## 介绍
一个轻量型的图片服务器


## 软件架构
上传接口：
http://127.0.0.1:8080/upload

参数：Files 类型：文件


返回结果：
[{"success":true,"message":"OK","version":"v0.1.1","data":{"size":49160,"mime":"image/jpeg","fileId":"5781339b809d5f18132f5c4fbe9df2fe","fileName":"gss0.baidu.jpg"}}]

使用说明：
http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe  默认：压缩质量为75%

访问原图：
http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?p=1   //p=1 查看原始图片

下载图片
http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?d=1  //d=1 下载图片，浏览器不再展示图片

灰阶图
http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?g=1  //g=1 灰阶图

缩放
http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?w=100&h=100  //w宽度 h高度, 只传递w或h 等比例缩放，同时传递两个值，等宽等高裁切

压缩
http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?q=75     //q 压缩质量 

转换格式
http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?f=png    //f 转换格式 ，默认jpg

旋转
http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?r=90   //r 旋转图像

裁切
http://127.0.0.1:8080/5781339b809d5f18132f5c4fbe9df2fe?x=10&y=10&w=100&h=100  //四个值同时传递时，x和y起始坐标点，w h 要裁切的宽度和高度


## 软件运行所需要的环境

### Linux

软件用到了ImageMagick库,需要安装该库

下载 https://github.com/ImageMagick/ImageMagick/archive/refs/tags/7.1.0-49.tar.gz
```bash
yum install libtool-ltdl-devel install libjpeg-devel libpng-devel libwebp-devel libtiff-devel zlib-devel freetype-devel openjpeg2-devel giflib-devel

./configure --prefix=/usr/local/ImageMagick --with-modules --enable-shared

make & make install

export PKG_CONFIG_PATH=/usr/local/ImageMagick/lib/pkgconfig/
export LD_LIBRARY_PATH=/usr/local/ImageMagick/lib
```

