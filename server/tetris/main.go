package main

func main() {
    s := NewTetrisServer(54545)
    s.Listen()
}