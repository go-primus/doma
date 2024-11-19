package svc

import (
	"context"
	"fmt"
)

type DemoService struct {
	Service
	content string
}

func newDemoService(content string) *DemoService {

	return &DemoService{
		Service: newBaseService("svc_demo", nil),
		content: content,
	}
}

func (s *DemoService) Start(ctx context.Context) error {
	fmt.Println("demo service: start")
	return s.Service.Start(ctx)
}

func (s *DemoService) Demo() {
	fmt.Println("echo demo ---> ", s.content)
}
