package LFS

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	Client *http.Client = &http.Client{}
)

// Resoonse : Struct to store the http response in orginzed way
type Response struct {
	Body, nextURL    string
	StatusCode       int
	Cookies, Headers map[string]string
}

// getCookiesAsString : Return the cookies as string in form of key1=value1; key2=value2...
func (c *Response) getCookiesAsString() string {
	cookie2 := []string{}
	for k, v := range c.Cookies {
		cookie2 = append(cookie2, k+"="+v)
	}
	return strings.Join(cookie2, "; ")
}

// InitClient : set the client up
// 0 or less for timeout will lead to the default settings for the client
func InitClient(Timeout int) {
	transporter := http.Transport{}
	transporter.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	if Timeout > 0 {
		Client.Timeout = time.Second * time.Duration(Timeout)
	}
	Client.Transport = &transporter
}

// SendRequest Send http request and return the the detials as response struct
func SendRequest(method string, url string, data string, headers map[string]string, cookies map[string]string, setNextURL bool) (Response, error) {
	body := func() io.Reader {
		if strings.ToUpper(method) == "POST" {
			return bytes.NewBufferString(data)
		} else {
			return nil
		}
	}()
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return Response{}, err
	}
	for header := range headers {
		request.Header.Add(header, headers[header])
	}
	if strings.ToUpper(method) == "POST" {
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
	for cookie := range cookies {
		request.AddCookie(&http.Cookie{Name: cookie, Value: cookies[cookie]})
	}
	httpresponse, err := Client.Do(request)
	if err != nil {
		return Response{}, err
	}
	var reader io.ReadCloser
	switch httpresponse.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(httpresponse.Body)
		defer reader.Close()
	default:
		reader = httpresponse.Body
	}
	bodyR, _ := ioutil.ReadAll(reader)
	resbody := string(bodyR)
	response := Response{}
	response.Body = string(resbody)
	response.StatusCode = httpresponse.StatusCode
	response.Cookies = make(map[string]string)
	response.Headers = make(map[string]string)
	if setNextURL {
		response.nextURL = httpresponse.Request.URL.String()
	}
	cookie := httpresponse.Cookies()
	for _, Cook := range cookie {
		response.Cookies[Cook.Name] = Cook.Value
	}
	for k := range httpresponse.Header {
		response.Headers[k] = httpresponse.Header.Get(k)
	}
	return response, nil
}
