package terminal

import "k8s.io/client-go/tools/remotecommand"

type TerminalSizeQueue interface {
	Next() *TerminalSize
}

type TerminalSize struct {
	// 终端的宽度
	Width uint16 `json:"width"`
	// 终端高度
	Height uint16 `json:"height"`
}

func NewTerminalSize() *TerminalSize {
	return &TerminalSize{}
}

type TerminalResizer struct {
	sizeChan chan remotecommand.TerminalSize
	dongChan chan struct{}
}

func NewTerminalSizer() *TerminalResizer {
	return &TerminalResizer{
		sizeChan: make(chan remotecommand.TerminalSize, 10),
		dongChan: make(chan struct{}),
	}
}

func (i *TerminalResizer) SetSize(ts TerminalSize) {
	i.sizeChan <- remotecommand.TerminalSize{Width: ts.Width, Height: ts.Height}
}

func (i *TerminalResizer) Next() *remotecommand.TerminalSize {
	select {
	case <-i.dongChan:
		return nil
	case size := <-i.sizeChan:
		return &size
	}
}
