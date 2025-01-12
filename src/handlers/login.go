package handlers

import (
    "encoding/json"
    "fmt"
    "strings"
    "github.com/valyala/fasthttp"
)

type LoginFetcher struct{}

type Session struct {
    PostResponse struct {
        StatusCode int `json:"status_code"`
        Lookup     struct {
            Identifier string `json:"identifier"`
            Digest     string `json:"digest"`
        } `json:"lookup"`
    } `json:"postResponse"`
    PassResponse struct {
        StatusCode int `json:"status_code"`
    } `json:"passResponse"`
    Cookies string `json:"Cookies"`
    Message string `json:"message"`
    Errors  string `json:"errors"`
}

func (lf *LoginFetcher) Logout(token string) (map[string]interface{}, error) {
    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    req.SetRequestURI("https://campuswebapi.up.railway.app/api/auth/logoutuser/")
    req.Header.SetMethod("GET")
    req.Header.Set("Accept", "*/*")
    req.Header.Set("Accept-Language", "en-US,en;q=0.5")
    req.Header.Set("x-csrf-token", token)
    req.Header.Set("Sec-Fetch-Site", "cross-site")
    req.Header.Set("Cache-Control", "private, max-age=120, stale-while-revalidate=1200, must-revalidate")
    req.Header.Set("Referer", "https://campusweb.vercel.app/")

    if err := fasthttp.Do(req, resp); err != nil {
        return nil, err
    }

    bodyText := resp.Body()
    
    result := map[string]interface{}{
        "status": resp.StatusCode(),
        "result": string(bodyText),
    }
    return result, nil
}

func (lf *LoginFetcher) CampusLogin(username, password string) (map[string]interface{}, error) {
    user := strings.Replace(username, "@srmist.edu.in", "", 1)
    body := fmt.Sprintf(`{"username":"%s@srmist.edu.in","password":"%s"}`, user, password)
    
    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    req.SetRequestURI("https://campuswebapi.up.railway.app/api/auth/login/")
    req.Header.SetMethod("POST")
    req.Header.Set("accept", "*/*")
    req.Header.Set("priority", "u=1, i")
    req.Header.Set("Referer", "https://campusweb.vercel.app/")
    req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
    req.Header.Set("Content-Type", "application/json")
    req.SetBody([]byte(body))

    if err := fasthttp.Do(req, resp); err != nil {
        return nil, err
    }

    if resp.StatusCode() != fasthttp.StatusOK {
        return lf.Login(user, password)
    }

    var session Session
    if err := json.Unmarshal(resp.Body(), &session); err != nil {
        return nil, err
    }

    statusCodes := map[string]int{
        "password": session.PassResponse.StatusCode,
        "lookup":   session.PostResponse.StatusCode,
    }

    if !strings.HasPrefix(fmt.Sprint(statusCodes["password"]), "2") || !strings.HasPrefix(fmt.Sprint(statusCodes["lookup"]), "2") {
        return lf.Login(user, password)
    }

    return map[string]interface{}{
        "authenticated": true,
        "session":       session,
        "lookup": map[string]string{
            "identifier": session.PostResponse.Lookup.Identifier,
            "digest":     session.PostResponse.Lookup.Digest,
        },
        "cookies": session.Cookies,
        "status":  session.PassResponse.StatusCode,
    }, nil
}

func (lf *LoginFetcher) Login(username, password string) (map[string]interface{}, error) {
    user := strings.Replace(username, "@srmist.edu.in", "", 1)
    
    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    req.SetRequestURI(fmt.Sprintf("https://academia.srmist.edu.in/accounts/p/10002227248/signin/v2/lookup/%s@srmist.edu.in", user))
    req.Header.SetMethod("POST")
    req.Header.Set("accept", "*/*")
    req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
    req.Header.Set("sec-fetch-site", "same-origin")
    req.Header.Set("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
    req.Header.Set("x-zcsrf-token", "iamcsrcoo=884b99c7-829b-4ddf-8344-ce971784bbe8")
    req.Header.Set("cookie", "f0e8db9d3d=7ad3232c36fdd9cc324fb86c2c0a58ad; bdb5e23bb2=3fe9f31dcc0a470fe8ed75c308e52278; zccpn=221349cd-fad7-4b4b-8c16-9146078c40d5; ZCNEWUIPUBLICPORTAL=true; cli_rgn=IN; iamcsr=884b99c7-829b-4ddf-8344-ce971784bbe8; _zcsr_tmp=884b99c7-829b-4ddf-8344-ce971784bbe8; 74c3a1eecc=d06cba4b90fbc9287c4162d01e13c516")

    if err := fasthttp.Do(req, resp); err != nil {
        return nil, err
    }

    var data map[string]interface{}
    if err := json.Unmarshal(resp.Body(), &data); err != nil {
        return nil, err
    }

    exists := strings.Contains(data["message"].(string), "User exists")

    if !exists {
        return map[string]interface{}{
            "session": false,
            "exists":  exists,
            "status":  data["status_code"],
        }, nil
    }

    session, err := lf.GetSession(password, data["lookup"].(map[string]interface{}))
    if err != nil {
        return nil, err
    }

    redir := session["passwordauth"].(map[string]interface{})["redirect_uri"].(string)
    sessionBody := map[string]interface{}{
        "success": strings.Contains(redir, "redirectFromLogin"),
        "code":    session["passwordauth"].(map[string]interface{})["code"],
        "message": session["message"],
    }

    if strings.Contains(redir, "sessions-reminder") || strings.Contains(redir, "block-sessions") {
        lf.Cleanup(session["cookies"].(string))
        return map[string]interface{}{
            "authenticated": true,
            "exists":        exists,
            "session":       sessionBody,
            "lookup": map[string]string{
                "identifier": data["lookup"].(map[string]interface{})["identifier"].(string),
                "digest":     data["lookup"].(map[string]interface{})["digest"].(string),
            },
            "cookies": session["cookies"],
            "status":  session["status_code"],
            "message": session["message"],
            "errors":  session["errors"],
        }, nil
    }

    if strings.Contains(strings.ToLower(session["message"].(string)), "invalid") || strings.Contains(session["cookies"].(string), "undefined") {
        return map[string]interface{}{
            "authenticated": false,
            "exists":        exists,
            "lookup": map[string]string{
                "identifier": data["lookup"].(map[string]interface{})["identifier"].(string),
                "digest":     data["lookup"].(map[string]interface{})["digest"].(string),
            },
            "session": sessionBody,
            "status":  data["status_code"],
            "message": session["message"],
            "errors":  session["errors"],
        }, nil
    }

    if strings.Contains(strings.ToLower(session["message"].(string)), "hip") || strings.Contains(strings.ToLower(session["localized_message"].(string)), "captcha") || session["cdigest"] != nil {
        return map[string]interface{}{
            "authenticated": false,
            "exists":        exists,
            "status":        session["status_code"],
            "message":       session["message"],
            "lookup": map[string]string{
                "identifier": data["lookup"].(map[string]interface{})["identifier"].(string),
                "digest":     data["lookup"].(map[string]interface{})["digest"].(string),
            },
        }, nil
    }

    return map[string]interface{}{
        "authenticated": true,
        "session":       sessionBody,
        "lookup": map[string]string{
            "identifier": data["lookup"].(map[string]interface{})["identifier"].(string),
            "digest":     data["lookup"].(map[string]interface{})["digest"].(string),
        },
        "cookies": session["cookies"],
        "status":  data["status_code"],
    }, nil
}

func (lf *LoginFetcher) GetSession(password string, lookup map[string]interface{}) (map[string]interface{}, error) {
    identifier := lookup["identifier"].(string)
    digest := lookup["digest"].(string)
    body := fmt.Sprintf(`{"passwordauth":{"password":"%s"}}`, password)

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    req.SetRequestURI(fmt.Sprintf("https://academia.srmist.edu.in/accounts/p/10002227248/signin/v2/primary/%s/password?digest=%s&cli_time=1713533853845&servicename=ZohoCreator&service_language=en&serviceurl=https://academia.srmist.edu.in/portal/academia-academic-services/redirectFromLogin", identifier, digest))
    req.Header.SetMethod("POST")
    req.Header.Set("accept", "*/*")
    req.Header.Set("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
    req.Header.Set("x-zcsrf-token", "iamcsrcoo=884b99c7-829b-4ddf-8344-ce971784bbe8")
    req.Header.Set("cookie", "f0e8db9d3d=7ad3232c36fdd9cc324fb86c2c0a58ad; bdb5e23bb2=3fe9f31dcc0a470fe8ed75c308e52278; zccpn=221349cd-fad7-4b4b-8c16-9146078c40d5; ZCNEWUIPUBLICPORTAL=true; cli_rgn=IN; iamcsr=884b99c7-829b-4ddf-8344-ce971784bbe8; _zcsr_tmp=884b99c7-829b-4ddf-8344-ce971784bbe8; 74c3a1eecc=d06cba4b90fbc9287c4162d01e13c516;")
    req.SetBody([]byte(body))

    if err := fasthttp.Do(req, resp); err != nil {
        return nil, err
    }

    var data map[string]interface{}
    if err := json.Unmarshal(resp.Body(), &data); err != nil {
        return nil, err
    }

    cookies := resp.Header.Peek("Set-Cookie")
    data["cookies"] = string(cookies)
    return data, nil
}

func (lf *LoginFetcher) Cleanup(cookie string) (int, error) {
    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    req.SetRequestURI("https://academia.srmist.edu.in/accounts/p/10002227248/webclient/v1/account/self/user/self/activesessions")
    req.Header.SetMethod("DELETE")
    req.Header.Set("accept", "*/*")
    req.Header.Set("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
    req.Header.Set("x-zcsrf-token", "iamcsrcoo=8cbe86b2191479b497d8195837181ee152bcfd3d607f5a15764130d8fd8ebef9d8a22c03fd4e418d9b4f27a9822f9454bb0bf5694967872771e1db1b5fbd4585")
    req.Header.Set("Referer", "https://academia.srmist.edu.in/accounts/p/10002227248/announcement/sessions-reminder?servicename=ZohoCreator&serviceurl=https://academia.srmist.edu.in/portal/academia-academic-services/redirectFromLogin&service_language=en")
    req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
    req.Header.Set("cookie", cookie)

    if err := fasthttp.Do(req, resp); err != nil {
        return 0, err
    }

    return resp.StatusCode(), nil
}
