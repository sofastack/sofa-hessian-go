// nolint
// Copyright 20xx The Alipay Authors.
//
// @authors[0]: bingwu.ybw(bingwu.ybw@antfin.com|detailyang@gmail.com)
// @authors[1]: robotx(robotx@antfin.com)
//
// *Legal Disclaimer*
// Within this source code, the comments in Chinese shall be the original, governing version. Any comment in other languages are for reference only. In the event of any conflict between the Chinese language version comments and other language version comments, the Chinese language version shall prevail.
// *法律免责声明*
// 关于代码注释部分，中文注释为官方版本，其它语言注释仅做参考。中文注释可能与其它语言注释存在不一致，当中文注释与其它语言注释存在不一致时，请以中文注释为准。
//
//

package javaobject

type SofaRPCRequest struct {
	TargetAppName           string                 `hessian:"targetAppName"`
	TargetServiceUniqueName string                 `hessian:"targetServiceUniqueName"`
	MethodName              string                 `hessian:"methodName"`
	MethodArgSigs           []string               `hessian:"methodArgSigs"`
	RequestProps            map[string]interface{} `hessian:"requestProps"`
}

func (s *SofaRPCRequest) GetJavaClassName() string {
	return "com.alipay.sofa.rpc.core.request.SofaRequest"
}

type SofaRPCResponse struct {
	IsError       bool              `hessian:"isError"`
	ErrorMsg      string            `hessian:"errorMsg"`
	AppResponse   interface{}       `hessian:"appResponse"`
	ResponseProps map[string]string `hessian:"responseProps"`
}

func (s *SofaRPCResponse) GetJavaClassName() string {
	return "com.alipay.sofa.rpc.core.response.SofaResponse"
}

type SofaRPCServerException struct {
	DetailMessage        string                     `hessian:"detailMessage"`
	Cause                interface{}                `hessian:"cause"`
	StackTrace           JavaLangStackTraceElements `hessian:"stackTrace"`
	SuppressedExceptions interface{}                `hessian:"suppressedExceptions"`
}

func (s *SofaRPCServerException) GetJavaClassName() string {
	return "com.alipay.remoting.rpc.exception.RpcServerException"
}
