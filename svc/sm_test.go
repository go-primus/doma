package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestXX(t *testing.T) {
	bb, _ := json.Marshal(12)
	fmt.Println("ssss: ", "xxx", string(bb))

	var str string
	json.Unmarshal([]byte(`"xxx"`), &str)

	fmt.Println("str:", str)
}

// var sm *ServiceRegistry

func init() {
	// sm = NewServiceRegistry()
	sm.RegisterService(newDemoService("Mebox Primus"))
	sm.RegisterService(newBaseService("svc_base", nil))
	sm.RegisterService(newBaseService("svc_base1", nil))
	sm.RegisterService(newBaseService("svc_base2", nil))
	sm.RegisterService(newBaseService("svc_base2", nil))

	// sm.Init(context.Background())
	sm.Config(context.Background())
	sm.Start(context.Background())
}

func TestSm(t *testing.T) {
	svc := sm.GetSystemService("svc_demo").(*DemoService)
	// fmt.Println(svc.Demo())

	svc.Demo()
	sm.DumpServices()
}
