package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fastProxy/app/common"
	"fastProxy/app/config"
	"fastProxy/app/models"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"
)

type Checker struct {
	data    sync.Map
	timeout int
}

type HTTPRequest struct {
	HeadBuf     []byte
	conn        *net.Conn
	Host        string
	Method      string
	URL         string
	hostOrURL   string
	isBasicAuth bool
	basicAuth   *BaseAuth
}

type BaseAuth struct {
	data        sync.Map
	authURL     string
	authOkCode  int
	authTimeout int
	authRetry   int
}

func NewHTTPRequest(inConn *net.Conn, bufSize int, basicAuth *BaseAuth) (req *HTTPRequest, client models.User, err error) {
	buf := make([]byte, bufSize)
	n := 0

	req = &HTTPRequest{
		conn:      inConn,
		basicAuth: basicAuth,
	}

	n, err = (*inConn).Read(buf[:])
	if err != nil {
		if err != io.EOF {
			if strings.Contains(err.Error(), "first record does not look like a TLS handshake") {
				common.CloseConn(inConn)
				return
			}
		}
		common.CloseConn(inConn)
		return
	}

	req.HeadBuf = buf[:n]

	index := bytes.IndexByte(req.HeadBuf, '\n')
	if index == -1 {
		common.CloseConn(inConn)
		return
	}
	fmt.Sscanf(string(req.HeadBuf[:index]), "%s%s", &req.Method, &req.hostOrURL)

	if req.Method == "" || req.hostOrURL == "" {
		common.CloseConn(inConn)
		return
	}
	req.Method = strings.ToUpper(req.Method)

	client, err = req.HTTP()
	return
}

func (req *HTTPRequest) HTTP() (subClient models.User, err error) {

	subClient, err = req.BasicAuth()
	if err != nil {
		return
	}

	req.URL = req.getHTTPURL()
	var u *url.URL
	u, err = url.Parse(req.URL)
	if err != nil {
		return
	}

	req.Host = u.Host
	req.addPortIfNot()
	return
}

func (req *HTTPRequest) BasicAuth() (subClient models.User, err error) {

	user, subClient, err := req.GetAuthDataStr()
	if err != nil {
		return
	}

	_subClient, authOk := (*req.basicAuth).Check(string(user))

	if !authOk {
		common.CloseConn(req.conn)
		return
	}
	subClient = *_subClient

	return subClient, nil
}

func (req *HTTPRequest) getHTTPURL() (URL string) {
	if !strings.HasPrefix(req.hostOrURL, "/") {
		return "http://" + req.hostOrURL
	}
	_host := req.getHeader("host")
	if _host == "" {
		return
	}
	URL = fmt.Sprintf("http://%s%s", _host, req.hostOrURL)
	return
}

func (req *HTTPRequest) getHeader(key string) (val string) {
	key = strings.ToUpper(key)
	lines := strings.Split(string(req.HeadBuf), "\r\n")
	//log.Println(lines)
	for _, line := range lines {
		hline := strings.SplitN(strings.Trim(line, "\r\n "), ":", 2)
		if len(hline) == 2 {
			k := strings.ToUpper(strings.Trim(hline[0], " "))
			v := strings.Trim(hline[1], " ")
			if key == k {
				val = v
				return
			}
		}
	}
	return
}

func (req *HTTPRequest) GetAuthDataStr() (user string, subClient models.User, err error) {
	authorization := req.getHeader("Proxy-Authorization")
	authorization = strings.Trim(authorization, " \r\n\t")
	if authorization == "" {
		fmt.Fprintf((*req.conn), "HTTP/1.1 %s Proxy Authentication Required\r\nProxy-Authenticate: Basic realm=\"\"\r\n\r\nProxy Authentication Required", "407")
		common.CloseConn(req.conn)
		err = errors.New("require auth header data")
		return
	}

	basic := strings.Fields(authorization)
	if len(basic) != 2 {
		err = fmt.Errorf("authorization data error,ERR:%s", authorization)
		common.CloseConn(req.conn)
		return
	}
	_user, err := base64.StdEncoding.DecodeString(basic[1])
	if err != nil {
		err = fmt.Errorf("authorization data parse error,ERR:%s", err)
		common.CloseConn(req.conn)
		return
	}
	user = string(_user)
	return
}

func (ba *BaseAuth) Check(userpass string) (*models.User, bool) {
	u := strings.Split(strings.Trim(userpass, " "), ":")

	if len(u) == 2 {
		if u[0] != config.GlobalConfig.ProxyCnf.UserName && u[1] != config.GlobalConfig.ProxyCnf.Password {
			return nil, false
		}
		_clientUser := &models.User{
			UserName: u[0],
			Password: u[1],
			Ok:       true,
		}
		return _clientUser, true
	}
	return nil, false
}

func (req *HTTPRequest) addPortIfNot() (newHost string) {
	//newHost = req.Host
	port := "80"
	if req.IsHTTPS() {
		port = "443"
	}
	if (!strings.HasPrefix(req.Host, "[") && strings.Index(req.Host, ":") == -1) || (strings.HasPrefix(req.Host, "[") && strings.HasSuffix(req.Host, "]")) {
		//newHost = req.Host + ":" + port
		//req.headBuf = []byte(strings.Replace(string(req.headBuf), req.Host, newHost, 1))
		req.Host = req.Host + ":" + port
	}
	return
}

func (req *HTTPRequest) IsHTTPS() bool {
	return req.Method == "CONNECT"
}
