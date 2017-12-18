package main

import (
	"io"
	"log"
	//"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	// "errors"
	"net"
	"time"
)

type NodeList struct {
	ApiVersion string        `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`
	Items      []interface{} `json:"items,omitempty" yaml:"items,omitempty"`
	Kind       string        `json:"kind,omitempty" yaml:"kind,omitempty"`
}

type ReturnJson struct {
	Node   string `json:"node,omitempty" yaml:"node,omitempty"`
	Deploy string `json:"deploy,omitempty" yaml:"deploy,omitempty"`
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	pool := x509.NewCertPool()

	//caCertPath := "/etc/cfc/conf/ca.crt"   //ICP ca file
	caCertPath := "/etc/kubernetes/pki/ca.crt"     // k8s ca file, kubeadm set up
	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)
	/*
		   tr := &http.Transport{
			   TLSClientConfig:    &tls.Config{RootCAs: pool},
			   DisableCompression: true,
		   }

		   client := &http.Client{Transport: tr}
	*/
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{RootCAs: pool},
			DisableCompression: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*5) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(10 * time.Second)) //设置发送接收数据超时
				return c, nil
			},
		},
	}

	reqest, _ := http.NewRequest("GET", "https://9.111.96.16:8001/api/v1/nodes", nil)
	/*kubectl describe secret $(kubectl get secrets |grep default-token-|awk '{print $1}')
	 *
	 *use token
	 *
	 */
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tcHJ4amMiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjZlMmY2OTM1LWQ4YjktMTFlNy1hMzc5LTAwNTA1NmIyMDdjNSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.U4tk1Nly4-pWnbcPnh4-w_1i1Ou4dVGAen3LjFbN7Di0q7ucVWOOr3Nop7Br4bD8_ngvu-pGGjYR_HYUSpMCfheEU3GX9XvJ7a9MoAI_2vqnv_6DGUw1I7S89S7HbIRisD6wv6HVaw5JNOmlzOBcG92E_pxqVRR_u9SKwz-xzoU6IPz7LwySKT3OLQg6NwcyENMnQSnFmEyzBemimE0SEl1s0scB6uUrFvVsGqOu8d1-yM-5OE4DWzWD3MbJMGTWf3qpAw_yuzcu4sxAyuGqSMmq9US9am9vwHhMxwKeMTI40gpQO_YYEdJeFzfu4NXCFBMI9VxMaokGq6Jd8t_8CQ"
	reqest.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(reqest)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	respObject := &NodeList{}
	json.Unmarshal(body, respObject)

	//------------------
	reqest1, _ := http.NewRequest("GET", "https://9.111.96.16:8001/apis/apps/v1beta1/deployments", nil)
	reqest1.Header.Set("Authorization", "Bearer "+token)
	resp1, err1 := client.Do(reqest1)
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	defer resp1.Body.Close()
	body1, _ := ioutil.ReadAll(resp1.Body)
	//fmt.Println(string(body1))
	respObject1 := &NodeList{}
	json.Unmarshal(body1, respObject1)
	countNode := strconv.Itoa(len(respObject.Items))
	countDeployment := strconv.Itoa(len(respObject1.Items))
	returnJson := ReturnJson{}
	returnJson.Node = countNode
	returnJson.Deploy = countDeployment
	byteContent, _ := json.Marshal(returnJson)
	io.WriteString(w, "vmsuccess("+string(byteContent)+")")

	//io.WriteString(w, strconv.Itoa(len(respObject.Items))+"-"+strconv.Itoa(len(respObject1.Items)))
}

func main() {
	http.HandleFunc("/getCount", HelloServer)
	err := http.ListenAndServe(":8094", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func readFile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	fmt.Println(string(fd))
	return string(fd)
}
