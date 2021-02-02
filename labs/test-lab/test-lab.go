package main

import (
    "fmt"
    "os"
)

func function() (string ) {

    name := ""
    if len(os.Args) == 1 {
        return "Error, Missing string"
    }

    for _,word := range os.Args[1:] {
        name = fmt.Sprintf("%v %v", name, word)
    }

    return fmt.Sprintf("Hello %s, Welcome to the Jungle", name)
}
func main() {
    str := function()
    fmt.Println(str)
}