### Oss Uploader

Oss上传工具可将本地文件夹中的文件全部上传到OSS对应的桶中，上传之前会自动清空桶中的所有文件。

**使用此工具可配合jenkins将打包好的网页自动上传到OSS的对应桶中，实现网页的自动发布**

1. 下载依赖库

   ```sh
   go get
   ```

2. 阿里提供的库中引用了 golang.org/x/time/rate 包，可自行进行下载，并放到对应的文件夹里

   ```sh
   go clone git@github.com:golang/time.git
   ```

3. 修改OSS配置 oss.config

   ```sh
   {
   	"EndPoint":"https://oss-cn-beijing.aliyuncs.com",
   	"AccessKeyId":"xxxxxxxxxxx",
   	"AccessKeySecret":"xxxxxxxxx"
   }
   ```

4. 编译/跨平台编译

   ```sh
   //windows
   GOOS=windows GOARCH=amd64 go build
   
   //linux
   GOOS=linux GOARCH=amd64 go build
   
   //mac
   GOOS=darwin GOARCH=amd64 go build
   ```

5. 启动参数说明

   ```
   -h 帮助
   bucket 为桶名称
   dir 为本地要上传的文件夹路径
   config 为配置文件路径，默认为 ./.oss.config
   ```

6. 更新桶文件

   ```sh
   $ ./ossuploader.exe -bucket oss-test -dir e:/test -config oss.config
   Delete File: [aaa.txt]
   获取的文件为[e:/test/aaa.txt], key:[aaa.txt]
   上传成功，文件数为[1]
   ```

   

