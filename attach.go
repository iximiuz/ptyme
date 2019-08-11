package main

import (
	"io"
	"net"
	"os"
	"sync"

	"golang.org/x/sys/unix"
)

// TODO: signal handling
// TODO: resizing
// TODO: wordwrap?

func main() {
	saved, err := tcget(os.Stdin.Fd())
	if err != nil {
		panic(err)
	}
	defer func() {
		tcset(os.Stdin.Fd(), saved)
	}()

	raw := makeraw(*saved)
	tcset(os.Stdin.Fd(), &raw)

	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		io.Copy(conn, os.Stdin)
		wg.Done()
	}()

	go func() {
		io.Copy(os.Stdout, conn)
		wg.Done()
	}()

	wg.Wait()
}

func tcget(fd uintptr) (*unix.Termios, error) {
	termios, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	if err != nil {
		return nil, err
	}
	return termios, nil
}

func tcset(fd uintptr, p *unix.Termios) error {
	return unix.IoctlSetTermios(int(fd), unix.TCSETS, p)
}

func makeraw(t unix.Termios) unix.Termios {
	t.Iflag &^= (unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON)
	t.Oflag &^= unix.OPOST
	t.Lflag &^= (unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN)
	t.Cflag &^= (unix.CSIZE | unix.PARENB)
	t.Cflag &^= unix.CS8
	t.Cc[unix.VMIN] = 1
	t.Cc[unix.VTIME] = 0
	return t
}
