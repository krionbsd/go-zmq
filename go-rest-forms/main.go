// forms.go
// krion - Apr. 2019
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type RunnerCall struct {
	Command     string    `json:"Command"`
	CommandArgs VmDetails `json:"CommandArgs"`
	NodeID      string    `json:"Node_Id"`
}

type VmDetails struct {
	JName        string `json:"jname"`
	CiIp4Addr    string `json:"ci_ip4_addr"`
	CiGw4        string `json:"ci_gw4"`
	Interface    string `json:"interface"`
	VmCpus       string `json:"vm_cpus"`
	ImgSize      string `json:"imgsize"`
	VmRam        string `json:"vm_ram"`
	VmOsType     string `json:"vm_os_type"`
	VmOsProfile  string `json:"vm_os_profile"`
	CiUserAdd    string `json:"ci_user_add"`
	CiUserPubKey string `json:"ci_user_pubkey"`
	RunAsap      string `json:"runasap"`
}

func main() {
	tmpl := template.Must(template.ParseFiles("forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		vmdetails := VmDetails{
			JName:        r.FormValue("jname"),
			VmRam:        r.FormValue("vm_ram"),
			VmCpus:       r.FormValue("vm_cpus"),
			ImgSize:      r.FormValue("imgsize"),
			VmOsType:     r.FormValue("vm_os_type"),
			VmOsProfile:  r.FormValue("vm_os_profile"),
			CiUserAdd:    r.FormValue("ci_user_add"),
			CiUserPubKey: r.FormValue("ci_user_pubkey"),
			CiIp4Addr:    r.FormValue("ci_ip4_addr"),
			CiGw4:        r.FormValue("ci_gw4"),
			RunAsap:      r.FormValue("runasap"),
		}

		vmdetails.Interface = "auto"

		bcreate := RunnerCall{
			Command:     "bcreate",
			CommandArgs: vmdetails,
			Node_id:     "auto",
		}

		// do something with details
		_ = vmdetails

		b, err := json.Marshal(vmdetails)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		fmt.Println(string(b))
		fmt.Println("---\n")

		c, err := json.Marshal(bcreate)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		length := len(string(c))

		if length > 10 {
			fmt.Println(string(c))
			cmdstr := fmt.Sprintf("/root/go-send-mq/go-send-mq %s", string(c))
			cmdstr = strings.TrimSuffix(cmdstr, "\n")
			arrCommandStr := strings.Fields(cmdstr)
			fmt.Printf("EXEC [%s]\n", cmdstr)
			cmd := exec.Command(arrCommandStr[0], arrCommandStr[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				fmt.Printf("%v\n", err)
			}

		}

		data := struct {
			RunnerStr string
		}{
			RunnerStr: string(c),
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
