# 包同步服务平台 - packageServer

## 组件

* manager
* worker



#### 组件描述

##### manager

管理端（manager）位于打包构建环境（CI）。用于仓库包管理，包推送至云端，上线计划制定。


##### worker

工作端（worker）位于私有平台（local）。用于从调度端接收服务包，执行上线计划，更新本地服务包，健康检查，异常回滚。





结构图

```
项目分层
server
|
router
| + middleware
controller
| + dto
service
|		|
repo 	rpc
|		|
model	request
|
db

```

