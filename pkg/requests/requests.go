package requests

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/xmlpath.v2"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)


const Timeout = 20

type Options struct {
	Client    *http.Client
	Method    string
	URL       string
	Timeout   int
	Headers   map[string]string
	Retry     int
}

type Response struct {
	*http.Response
	Selector      *Selector
	Bytes         []byte
	History       []*http.Request
}

var DefaultHeaders = map[string]string{
	"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
	"Accept-Encoding": "",
	"Accept-Language": "zh-CN,zh;q=0.9",
	"Connection": "keep-alive",
	"User-Agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36",
}

type Selector struct {
	body    []byte
	reader  io.Reader
	root    *xmlpath.Node
}

func Node2Selector(node *xmlpath.Node) *Selector {
	return &Selector{
		body: node.Bytes(),
		root: node,
	}
}

func (selector *Selector) Gets(xpath interface{}) ([]*xmlpath.Node, error) {
	if selector.reader == nil {
		selector.reader = bytes.NewReader(selector.body)
	}

	var err error
	root := selector.root
	if root == nil {
		root, err = xmlpath.ParseHTML(selector.reader)
		if err != nil {
			panic(err)
		}
		selector.root = root
	}

	var xpaths []string
	switch xpath.(type) {
	case string:
		xpaths = []string{xpath.(string)}
	case []string:
		xpaths = xpath.([]string)
	}

	for _, xpath := range xpaths {
		result := []*xmlpath.Node{}
		_path := xmlpath.MustCompile(xpath)
		iter := _path.Iter(root)
		for iter.Next() {
			result = append(result, iter.Node())
		}

		if len(result) > 0 {
			return result, nil
		}

	}

	return []*xmlpath.Node{}, NodeNotFound(xpaths)
}

func (selector Selector) Get(xpath interface{}) (*xmlpath.Node, error) {
	node, err := selector.Gets(xpath)
	if err != nil {
		return nil, err
	} else {
		return node[0], err
	}
}

func (selector *Selector) GetNode(xpaths interface{}) (*Selector, error) {
	node, err := selector.Get(xpaths)
	if err != nil {
		return nil, err
	} else {
		return Node2Selector(node), nil
	}
}

func (selector *Selector) GetNodes(xpaths interface{}) ([]*Selector, error) {
	selectors := []*Selector{}
	nodes, err := selector.Gets(xpaths)

	if err != nil {
		return selectors, err
	} else {
		for _, item := range nodes {
			selectors = append(selectors, Node2Selector(item))
		}

		return selectors, nil
	}
}


func NodeNotFound(xpaths []string) error {
	return fmt.Errorf("HTML node not found. xpaths: %s", strings.Join(xpaths,","))
}


func InitCookie(client *http.Client) (*http.Client,error ){
	jar,err := cookiejar.New(nil)
	client.Jar = jar
	return client,err
}

func InitClient() (*http.Client, error) {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:true,
		},
	}
	return InitCookie(client)
}

func Request(options Options) (*Response, error) {
	method := options.Method
	method = strings.ToUpper(method)

	//set client
	client := options.Client
	if client == nil {
		client,_ = InitClient()
	}

	//set history
	var history []*http.Request
	if client.CheckRedirect == nil {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			history = via
			if len(via) > 10 {
				return errors.New("stopped after 10 redirects")
			}
			history = append(history, req)
			return nil
		}
	}


	//set headers
	headers := http.Header{}

	for name, value := range DefaultHeaders {
		headers.Set(name, value)
	}

	if options.Headers != nil {
		for name, value := range options.Headers {
			headers.Set(name,value)
		}
	}

	//new request
	request, err := http.NewRequest(method, options.URL, nil)
	if err != nil {
		panic(err)
	}

	//set timeout
	timeout := Timeout
	if options.Timeout != 0 {
		timeout = options.Timeout
	}

	client.Timeout = time.Duration(timeout) * time.Second

	request.Header = headers

	retry := options.Retry
	if retry == 0 {
		retry = 5
	}

	var response *http.Response
	var response_body []byte

	for i := 0; i < retry; i++ {
		response, err = client.Do(request)
		if err == nil && response.StatusCode <  500 {
			response_body, err = ioutil.ReadAll(response.Body)
			if err == nil {
				break
			}
		}

		if err == nil && i+1 < retry {
			response.Body.Close()
		}

		time.Sleep(time.Second * 1)
	}

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	//defer response.Body.Close()
	//fmt.Println(response.StatusCode)

	//if response.StatusCode == 200 {
	//	body, err := ioutil.ReadAll(response.Body)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	fmt.Println(string(body))
	//
	//}
	//body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(body))
	return &Response{response,&Selector{body:response_body},response_body,history}, err
}
