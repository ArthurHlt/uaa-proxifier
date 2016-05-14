package main

import (
	"io/ioutil"
	"bytes"
	"strconv"
	"net/http"
	"encoding/json"
	"regexp"
	"log"
	"strings"
)

type TransportUserInfo struct {
	http.RoundTripper
}

func (t *TransportUserInfo) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	var uaaUser UserFromUaa
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if req.URL.Path != "/userinfo" {
		t.loadResponse(b, resp)
		return resp, err
	}

	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &uaaUser)
	if err != nil {
		t.loadResponse(b, resp)
		return resp, nil
	}
	id, err := t.generateId(uaaUser.UserID)
	if err != nil {
		t.loadResponse(b, resp)
		return resp, nil
	}
	userGitlab := &GitLabUser{
		Id: id,
		Username:uaaUser.UserName,
		Login:uaaUser.UserName,
		Email:uaaUser.Email,
		Name:uaaUser.UserName,
	}
	finalResp, _ := json.Marshal(userGitlab)
	t.loadResponse(finalResp, resp)
	return resp, nil
}
func (t *TransportUserInfo) loadResponse(b []byte, resp *http.Response) {
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
}
func (t *TransportUserInfo) generateId(guid string) (int64, error) {
	reg, err := regexp.Compile("([A-Za-z]|-)*")
	if err != nil {
		log.Print(err)
		return 0, err
	}
	safe := reg.ReplaceAllString(guid, "")
	safe = strings.ToLower(strings.Trim(safe, ""))
	safe = safe[:9]
	id, err := strconv.ParseInt(safe, 10, 64)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return id, err
}