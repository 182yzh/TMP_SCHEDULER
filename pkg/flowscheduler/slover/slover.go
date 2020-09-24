package main

import "fmt"
import "os/exec"
import "io/ioutil"
import "os"

func main(){
    file,err := os.Open("temtest.txt")
	defer file.Close()
	if err != nil {
		fmt.Println("can not open file\n")
    } 
    input_bytes,err := ioutil.ReadAll(file)
    if err != nil{
        fmt.Println(err)
        return
    }
    input := string(input_bytes)
	fmt.Println("__\n"+input+"__\n")
    cmd := exec.Command("./cs2.exe")
    stdout, err := cmd.StdoutPipe()
    //defer stdout.Close()
	if err != nil{
        return 
	}
    
    fmt.Println("++++++")
    stdin, err := cmd.StdinPipe()
    
	if err != nil{
        return 
    }
    cmd.Start()
    stdin.Write([]byte(input))
    stdin.Close()

    out_bytes, _ := ioutil.ReadAll(stdout)
    
    //out_bytes,_ := cmd.Output()
    output := string(out_bytes)
    fmt.Println(output)
    cmd.Wait()
    for i:=0;i!= 10000000000;i++{
        fmt.Sprintln("test")
    }
    
    return 
}