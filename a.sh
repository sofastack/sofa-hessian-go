#! /bin/bash

x=$(bin/hessian fromjson '{
  "$class": "com.alipay.sofa.rpc.core.request.SofaRequest",
  "$": {
    "requestProps": {
      "mesh_vip_enforce": "false",
      "rpc_trace_context": {
        "sofaRpcId": "0.1.2",
        "zproxyUID": "",
        "zproxyTargetZone": "",
        "sofaCallerIp": "100.88.185.222",
        "sofaPenAttrs": "group4=alipay.net&abskey=999&group=alipay.net&",
        "zproxyTimeout": "3000",
        "Elastic": "F",
        "sysPenAttrs": "",
        "sofaTraceId": "6458a54e15847884117178038e1f4c",
        "sofaCallerApp": "ifcriskmatrixus"
      },
      "mesh_vip_only": "false"
    },
    "targetAppName": "-",
    "targetServiceUniqueName": "com.alipay.isecuritycore.service.analyze.facade.SecurityPolicyService:1.0",
    "methodName": "querySecurityPolicy",
    "methodArgSigs": [
      "com.alipay.isecuritycore.service.analyze.model.SecurityBusiness"
    ]
  }
}' \
'{
    "signatures": [
        "com.alipay.isecuritycore.service.analyze.model.SecurityBusiness"
    ],
    "arguments": [
        {
            "$": {
                "eventInfo": {
                    "$class": "com.alibaba.fastjson.JSONObject",
                    "$": {
                        "ipayUserId": "2188211939800093"
                    }
                },
                "params": {
                    "API_CODE_LIST": [
                        {
                            "$class": "com.alipay.isecuritycore.service.analyze.model.PayChannelApiSummary",
                            "$": {
                                "apiCode": "CYBSOCBC_MIXEDCARD_ASYNC_000_001",
                                "isAvailable": false,
                                "reason": null
                            }
                        }
                    ]
                },
                "gmtOccur": {
                    "$class": "java.util.Date",
                    "$": 1584794520
                },
                "securityId": null,
                "sceneId": "PAY_AUTH",
                "serverId": "ifcriskmatrixus-eu95-2.gz00b.stable.alipay.net",
                "eventId": null
            },
            "$class": "com.alipay.isecuritycore.service.analyze.model.SecurityBusiness"
        }
    ]
}')

echo $x

bin/sofa curl bolt://11.166.128.220:12200 \
-H "service: com.alipay.isecuritycore.service.analyze.facade.SecurityPolicyService:1.0" \
-d "$x" --request-body-format hex
