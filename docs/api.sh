# GET    /ping                     --> PackageServer/controller.CommonApi.Ping
curl -X GET \
"http://localhost:8080/ping"


# GET    /timestamp                --> PackageServer/controller.CommonApi.TimeStamp
curl -X GET \
"http://localhost:8080/timestamp?user=manager"


# POST   /api/v1/digestToken       --> PackageServer/controller.RpcLoginApi.RetrieveRpcToken
curl -X POST \
-H "Content-Type: application/json" \
-d '{
  "User": "wuji",
  "LoginToken":"secret"
}' "http://localhost:8080/api/v1/digestToken"


# GET    /api/v1/rpcLogin          --> PackageServer/controller.RpcLoginApi.RpcLogin
curl -X GET \
"http://localhost:8080/api/v1/rpcLogin?area=spd-local"


# POST   /api/v1/auth/rpcToken/refreshExpire --> PackageServer/controller.RpcLoginApi.RefreshTokenExipireTime
curl -X POST \
-d '{
  "User": "manager",
  "ExpireTime": 3600
}' \
"http://localhost:8002/api/v1/auth/rpcToken/refreshExpire"


# GET    /api/v1/auth/rpcToken/del --> PackageServer/controller.RpcLoginApi.TokenDelete
curl -X GET \
"http://localhost:8002/api/v1/auth/rpcToken/del?area=manager"


#POST   /api/v1/auth/areaInfo    --> PackageServer/controller.AreaInfoApi.Create
curl -X POST \
-H "Content-Type: application/json" \
-d '{
  "AreaName" :  "spd-local",
  "IpAddr":   "127.0.0.1",
  "Port": "8002",
  "UrlPath":  "/",
  "Password": "95588"
}' \
"http://localhost:8080/api/v1/auth/areaInfo"


#GET    /api/v1/auth/areaInfo    --> PackageServer/controller.AreaInfoApi.List
curl -X GET \
-H "Content-Type: application/json" \
"http://localhost:8080/api/v1/auth/areaInfo"


#DELETE /api/v1/auth/areaInfo/:areaName --> PackageServer/controller.AreaInfoApi.DeleteByName
curl -X DELETE \
"http://localhost:8080/api/v1/auth/areaInfo/spd_demo"


#GET    /api/v1/auth/areaInfo/:areaName --> PackageServer/controller.AreaInfoApi.GetByName
curl -X GET \
-H "Content-Type: application/json" \
"http://localhost:8080/api/v1/auth/areaInfo/spd-local"


#PUT    /api/v1/auth/areaInfo/:areaName --> PackageServer/controller.AreaInfoApi.UpdateByName
curl -X PUT \
-H "Content-Type: application/json" \
-d '{
  "IpAddr":   "127.0.0.1",
  "Port": "8002",
  "UrlPath":  "/"
}' \
"http://localhost:8080/api/v1/auth/areaInfo/spd-local"


#GET    /api/v1/auth/areaInfo/:areaName/services --> PackageServer/controller.AreaInfoApi.ListServices
curl -X GET \
"http://localhost:8080/api/v1/auth/areaInfo/spd-local/services"


#POST   /api/v1/auth/areaInfo/:areaName/services/add --> PackageServer/controller.AreaInfoApi.AddServices
curl -X POST \
-H "Content-Type: application/json" \
-d '{
  "ServiceList": ["spd-web", "spd-scm","spd-wms"]
}' \
"http://localhost:8080/api/v1/auth/areaInfo/spd-local/services/add"


#POST   /api/v1/auth/areaInfo/:areaName/services/delete --> PackageServer/controller.AreaInfoApi.DelServices
curl -X POST \
-H "Content-Type: application/json" \
 -d '{
  "ServiceList": ["spd-web"]
}' \
"http://localhost:8080/api/v1/auth/areaInfo/spd_uat/services/delete"


#POST   /api/v1/auth/serviceOnline --> PackageServer/controller.ServiceOnlineApi.Create
curl -X POST \
-H "Content-Type: application/json" \
-d '{
  "ServiceName": "spd-wms",
  "ArchiveType": "jar"
}' \
"http://localhost:8080/api/v1/auth/serviceOnline"


#GET    /api/v1/auth/serviceOnline/ --> PackageServer/controller.ServiceOnlineApi.List
curl -X GET \
-H "Content-Type: application/json" \
"http://localhost:8080/api/v1/auth/serviceOnline"


#DELETE /api/v1/auth/serviceOnline/:serviceName --> PackageServer/controller.ServiceOnlineApi.DeleteByName
curl -X DELETE \
"http://localhost:8080/api/v1/auth/serviceOnline/spd-third-join"


#GET    /api/v1/auth/serviceOnline/:serviceName --> PackageServer/controller.ServiceOnlineApi.GetByName
curl -X GET \
-H "Content-Type: application/json" \
"http://localhost:8080/api/v1/auth/serviceOnline/spd-scm"


#PUT    /api/v1/auth/serviceOnline/:serviceName --> PackageServer/controller.ServiceOnlineApi.UpdateByName
curl -X POST \
-H "Content-Type: application/json" \
-d '{
  "ServiceName": "spd-wms",
  "ArchiveType": "jar"
}' \
"http://localhost:8080/api/v1/auth/serviceOnline/spd-wms"



#GET    /api/v1/auth/serviceOnline/:serviceName/areas --> PackageServer/controller.ServiceOnlineApi.ListAreas
curl -X GET \
-H "Content-Type: application/json" \
"http://localhost:8080/api/v1/auth/serviceOnline/spd-scm/areas"


#POST   /api/v1/auth/publishPlan/ --> PackageServer/controller.PublishPlanApi.PublishVersion
curl -X POST \
-H "Content-Type: application/json" \
-d '{
  "AreaName": "spd-local",
  "ServiceName": "spd-wms", 
  "Version": "20220825-1110"
}' \
"http://localhost:8080/api/v1/auth/publishPlan"


# POST   /api/v1/auth/package/info --> PackageServer/controller.PackageReceiveApi.PackInfo
curl -X POST \
-H "Content-Type: application/json" \
-d '{
  "FileName": "spd-wms_tag_spd-local_20220825-1110.jar",
  "Md5": "3ee7fd273b5a1a37f72e7b35a15a7e40",
  "ChunkNum": 48
}' \
"http://localhost:8002/api/v1/auth/package/info"


# POST   /api/v1/auth/package/chunkUpload --> PackageServer/controller.PackageReceiveApi.ChunkUpload
curl -X POST \
-H "packageName: spd-wms_tag_spd-local_20220825-1110.jar" \
-H "Content-Type: multipart/form-data" \
-F "chunk=@/Users/wujimaster/tmp/package/manager/spd-wms_tag_spd-local_20220825-1110.jar.chunk_1.tmp" \
"http://localhost:8002/api/v1/auth/package/chunkUpload"



# GET    /api/v1/auth/package/check --> PackageServer/controller.PackageReceiveApi.PackCheck
curl -X GET \
-H "Content-Type: application/json" \
"http://localhost:8002/api/v1/auth/package/check?packageName=spd-wms_tag_spd-local_20220825-1110.jar"

