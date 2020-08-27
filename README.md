# GRPC测试例子

[TOC]

## GRPC基础知识
### GO编译参数

- -I 参数：指定import路径，可以指定多个-I参数，编译时按顺序查找，不指定时默认查找当前目录,就是如果多个proto文件之间有互相依赖，生成某个proto文件时，需要import其他几个proto文件，这时候就要用-I来指定搜索目录
- --go_out ：golang编译支持，支持以下参数
    - plugins=plugin1+plugin2 - 指定插件，目前只支持grpc，即：plugins=grpc
    - M 参数 - 指定导入的.proto文件路径编译后对应的golang包名(不指定本参数默认就是.proto文件中import语句的路径)
    - import_prefix=xxx - 为生成的go文件的所有import路径添加前缀
    - import_path=foo/bar - 用于指定proto文件中未声明package或go_package的文件的包名，最右面的斜线前的字符会被忽略
    - 末尾 :编译文件路径 .proto文件路径(支持通配符)

### 示例
**完整示例**
```shell script
protoc -I . --go_out=plugins=grpc,Mfoo/bar.proto=bar,import_prefix=foo/,import_path=foo/bar:. ./*.proto
```

**示例含参数说明**
```shell script
protoc -I . --go_out=plugins=grpc,\   # 指定输出go代码(也可以指定其它语言代码)
            Mhelloworld.proto=dpy,\   # M参数指定proto文件中import其它proto对应的报名,如此处指示import helloworld.proto在go代码中包名为dpy
            import_prefix=dpy,\   # 指定生成的go代码中import路径增加前缀,通常用于修正go项目引入路径
            import_path=dpy:\  # 强制指定编译后的包名
            ./*.proto # 编译的目标文件，支持通配符
```

**-I参数**

如果多个proto文件之间有互相依赖，生成某个proto文件时，需要import其他几个proto文件，这时候就要用-I来指定搜索目录

如果没有指定 –I 参数，则在当前目录进行搜索

**M参数作用:通常用于修正.proto引入包转化为go包后的路径**
```shell script
# protoA文件中引入了protoB
import "path/to/protoB.proto";

# 此时protoB文件中定义了go_package选项
option go_package = "github.com/golang/protobuf/ptypes/timestamp";

# 此时编译protoA时使用了如下M参数
protoc -I ./ --go_out=plugins=grpc,Mpath/to/protoB.proto=dpy:. ./protoA.proto

# 则编译后的golang源文件中protoB包名就为dpy
import dpy "dpy"; 

# 如果编译protoA时不带M参数,则默认使用引入的protoB中的go_package字段路径
protoc -I ./ --go_out=plugins=grpc:. ./protoA.proto

# 编译后的golang包中的import选项为
import timestamp "github.com/golang/protobuf/ptypes/timestamp"; 
```

**import_prefix**
```shell script
# 在编译时加上import_prefix=path/prefix/
protoc -I ./ --go_out=plugins=grpc,import_prefix=path/prefix/:. ./protoA.proto

# 则在生成出来的golang代码中的所有的非系统包的，import都会加上path/prefix/，如：
import (
	context "dpy/context"
	proto "dpy/github.com/golang/protobuf/proto"
	timestamp "dpy/github.com/golang/protobuf/ptypes/timestamp"
	grpc "dpy/google.golang.org/grpc"
	codes "dpy/google.golang.org/grpc/codes"
	status "dpy/google.golang.org/grpc/status"
    # 系统包不会加上路径
	fmt "fmt"
	math "math"
)
```

**import_path**
```shell script
# 在编译时加上import_path=dpy
protoc -I ./ --go_out=plugins=grpc,import_path=dpy:. ./protoA.proto

# 生成出来的golang包名即为dpy，且无论proto文件中定义的option go_package或package定义如何
package dpy
```

## 测试例目录
### 1.简单的RPC通信
1. 在simple-proto文件夹中新建hello.proto，并定义hello服务
2. 使用`protoc -I ./ --go_out=plugins=grpc:./simple-proto-out ./simple-proto/*.proto`生成go代码文件
3. `simple-proto-out`目录中保存了proto编译后的代码文件
4. `simple-use`目录演示了简单的一元交互
5. `simple-use-interceptor`目录演示了在简单一元交互基础上如何使用拦截器

### 2.流模式RPC通信
#### 客户端流模式
1. 在stream-proto文件夹中新建stream_hello.proto,并定义StreamHello服务
2. 使用`protoc -I ./ --go_out=plugins=grpc:./stream-proto-out ./stream-proto/*.proto`生成go代码文件
