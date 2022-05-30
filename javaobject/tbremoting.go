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

type TBRemotingRequestContext struct {
	ID   int64                        `hessian:"id"`
	This *TBRemotingConnectionRequest `hessian:"this$0"`
}

func (c *TBRemotingRequestContext) Reset() {
	c.ID = 0
	c.This = nil
}

func (c *TBRemotingRequestContext) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionRequest$RequestContext"
}

type TBRemotingConnectionRequest struct {
	Ctx                *TBRemotingRequestContext `hessian:"ctx"`
	FromAppKey         string                    `hessian:"fromAppKey"`
	ToAppKey           string                    `hessian:"toAppKey"`
	EncryptedToken     []byte                    `hessian:"encryptedToken"`
	ApplicationRequest interface{}               `hessian:"-"`
}

func (c *TBRemotingConnectionRequest) Reset() {
	if c.Ctx != nil {
		c.Ctx.Reset()
	}
	c.FromAppKey = ""
	c.ToAppKey = ""
	c.EncryptedToken = c.EncryptedToken[:0]
	c.ApplicationRequest = nil
}

func (c *TBRemotingConnectionRequest) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionRequest"
}

type TBRemotingConnectionResponseContext struct {
	ID int64 `hessian:"id"`
}

func (c *TBRemotingConnectionResponseContext) Reset() {
	c.ID = 0
}

func (c *TBRemotingConnectionResponseContext) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionResponse$ResponseContext"
}

type TBRemotingConnectionResponse struct {
	Ctx        *TBRemotingConnectionResponseContext
	Host       string `hessian:"host"`
	Result     int32  `hessian:"result"`
	ErrorMsg   string `hessian:"errorMsg"`
	ErrorStack string `hessian:"errorStack"`
	FromAppKey string `hessian:"fromAppKey"`
	ToAppKey   string `hessian:"toAppKey"`
}

func (c *TBRemotingConnectionResponse) Reset() {
	if c.Ctx != nil {
		c.Ctx.Reset()
	}
}

func (c *TBRemotingConnectionResponse) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionResponse"
}
