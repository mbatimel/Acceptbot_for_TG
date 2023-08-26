package telegram

import (
	"encoding/json"
	"example/main/src/lib/e"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)
const(
	 errMsg = "can't do request"
	sendMassage = "sendMassage"
)

type Client struct {
	host string
	basepPath string
	client http.Client
}

func New(host string, token string) Client {
	return Client{
		host: host,
		basepPath: newBasepPath(token),
		client: http.Client{},
	}

}
func newBasepPath( token string) string{
	return "bot" + token
}
func (c *Client) Updates(offset int , limit int) ([]Update, error) {
 q := url.Values{}
 q.Add("offset", strconv.Itoa(offset))
 q.Add("limit", strconv.Itoa(limit))
 data, err := c.doRequest("getUpdate", q)
 if err != nil {
	return nil,e.WrapIfErr(errMsg,err)
 }
 var res UpdateResponse
 if err := json.Unmarshal(data, &res); err != nil{
	return nil, err
 }
 return res.Result, nil

}
func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {

	u:=url.URL{
		Scheme: "https",
		Host: c.host,
		Path: path.Join(c.basepPath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.WrapIfErr(errMsg,err)
	}
	req.URL.RawQuery = query.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil,e.WrapIfErr(errMsg,err)
	}
	body , err := io.ReadAll(resp.Body)
	if err != nil {
		return nil,e.WrapIfErr(errMsg,err)
	}
	return body, nil



}
func (c *Client) SendMassage(chatId int, text string) error {
	q:= url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)
	_,err:=c.doRequest(sendMassage,q)
	if err != nil {
		return e.Wrap("can't send massage", err)
	}
	return nil



}