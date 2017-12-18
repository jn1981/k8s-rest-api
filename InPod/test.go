package main

import (
    //"io"
    //"log"
    "crypto/tls"
    "crypto/x509"
    //"encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    //"strings"
    "os"
)

func main() {
    //x509.Certificate.
    pool := x509.NewCertPool()
    
    caCertPath := "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"

    caCrt, err := ioutil.ReadFile(caCertPath)
    if err != nil {
        fmt.Println("ReadFile err:", err)
        return
    }
    pool.AppendCertsFromPEM(caCrt)
   

    tr := &http.Transport{
        TLSClientConfig:    &tls.Config{RootCAs: pool},
        DisableCompression: true,
    }
    client := &http.Client{Transport: tr}
    reqest, _ := http.NewRequest("POST", "https://kubernetes.default/api/v1/namespaces", nil)
    reqest.Header.Set("Authorization","Bearer "+ readFile("/var/run/secrets/kubernetes.io/serviceaccount/token"))

    resp,err := client.Do(reqest)
    if err != nil {
        fmt.Println(err.Error())
    }



    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
    fmt.Println(resp.Status)
}

func readFile(path string)string{
    fi,err := os.Open(path)
    if err != nil{panic(err)}
    defer fi.Close()
    fd,err := ioutil.ReadAll(fi)
    fmt.Println(string(fd))
    return string(fd)
}

