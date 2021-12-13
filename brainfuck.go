package main

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Program struct {
	pos      int
	prog     []byte
	memSize  int
	mem      []uint64
	brackets *list.List
}

func checkSyntax(prog []byte) bool {
	check := 0
	for i := 0; i < len(prog); i++ {
		if prog[i] == '[' {
			check++
		} else if prog[i] == ']' {
			if check <= 0 {
				return false
			} else {
				check--
			}
		}
	}
	return check == 0
}

func getInput() (uint64, error) {
	var res uint64
	fmt.Print("input : ")
	r := bufio.NewReader(os.Stdin)
	_, err := fmt.Fscan(r, &res)
	r.Discard(r.Buffered())
	if err != nil {
		return 0, err
	}
	return res, nil
}

func NewProg(filename string, memSize int) (*Program, error) {
	prog, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if !checkSyntax(prog) {
		return nil, errors.New("SYNTAX ERROR")
	}
	r := new(Program)
	r.pos = 0
	r.prog = prog
	r.memSize = memSize
	r.mem = make([]uint64, memSize)
	r.brackets = list.New()
	return r, nil
}

func (p *Program) run() error {
	for i := 0; i < len(p.prog); i++ {
		switch p.prog[i] {
		case '+':
			p.mem[p.pos]++
		case '-':
			p.mem[p.pos]--
		case '>':
			p.pos++
			if p.pos >= p.memSize {
				return errors.New("POINTER OUT OF BOUNDS")
			}
		case '<':
			p.pos--
			if p.pos < 0 {
				return errors.New("POINTER OUT OF BOUNDS")
			}
		case '.':
			fmt.Print(string(p.mem[p.pos]))
		case ',':
			val, err := getInput()
			if err != nil {
				return err
			}
			p.mem[p.pos] = val
		case '[':
			p.brackets.PushFront(i)
		case ']':
			if p.brackets.Len() == 0 {
				// should not happen since we check syntax at the beginning, but you never know
				return errors.New("SYNTAX ERROR")
			}
			bracket := p.brackets.Front()
			if p.mem[p.pos] != 0 {
				i, _ = (*bracket).Value.(int)
			} else {
				p.brackets.Remove(bracket)
			}
		default:
		}
	}
	if p.brackets.Len() != 0 {
		// should not happen since we check syntax at the beginning, but you never know
		return errors.New("SYNTAX ERROR")
	}
	return nil
}

func usage() {
	fmt.Println("Usage : brainfuck PROGRAMM [memsize]")
}

func main() {
	var filename string
	var memsize int

	switch len(os.Args) {
	case 2:
		filename = os.Args[1]
		memsize = 30000
	case 3:
		filename = os.Args[1]
		_, err := fmt.Sscan(os.Args[2], &memsize)
		if err != nil {
			usage()
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}

	pgm, err := NewProg(filename, memsize)
	if err != nil {
		log.Fatal(err)
	}
	err = pgm.run()
	if err != nil {
		log.Fatal(err)
	}
}
