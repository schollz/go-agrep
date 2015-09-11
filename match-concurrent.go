package main

import (
    "fmt"
    "github.com/arbovm/levenshtein"
    "bufio"
    "os"
    "strings"
	"runtime"
	"sync"
)


func abs(x int) (int) {
    if x < 0 {
        return -x
    } else if x == 0 {
        return 0 // return correctly abs(-0)
    }
    return x
}

func generateHash(path string) {
  inFile, _ := os.Open(path)
  defer inFile.Close()
  scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines) 
  
  for scanner.Scan() {
    s := strings.Replace(scanner.Text(),"/","",-1)
    partials, num := getPartials(s)
    for i := 0; i < num; i ++ {
        addToCache(partials[i],s)
    }
  }
}


func addToCache(spartial string, s string) {
    f, err := os.OpenFile("cache/" + spartial, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        fmt.Println("%v",spartial)
        panic(err)
    }

    defer f.Close()

    if _, err = f.WriteString(s+"\n"); err != nil {
        fmt.Println("%v",spartial)
        panic(err)
    }
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
            break
        }
    }
    return false
}

func getPartials(s string) ([]string, int) { 
    partials := make([]string,1000)
    num := 0
    s = strings.Replace(s,"/","",-1)
    slen := len(s)
    if slen <= 3 {
       partials[num] = "asdf"
       num = num + 1
    } else {
        for i := 0; i <= slen-3; i ++ {
            partials[num] = s[i:i+3]
            num = num + 1
        }
    }
    return partials,num
}

func getMatch(s string) (string) {
    match := "No match"
    partials, num := getPartials(s)
    matches := make([]string,10000)
    numm := 0


    for i := 0; i < num; i ++ {
    
      inFile, _ := os.Open("cache/"+partials[i])
      defer inFile.Close()
      scanner := bufio.NewScanner(inFile)
        scanner.Split(bufio.ScanLines) 
      
      for scanner.Scan() {
        //if stringInSlice(scanner.Text(),matches) == false { ITS NOT WORTH LOOKING THROUGH DUPLICATES
            matches[numm] = scanner.Text()
            numm = numm + 1
       // }
      }
    
    }


    bestLevenshtein := 1000
    
	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup
	wg.Add(2)
    
	go func() {
		defer wg.Done()

    for i := 0; i < numm; i ++ {
        if (i%2==0) {
              d := levenshtein.Distance(s, matches[i])
              if (d < bestLevenshtein) {
                bestLevenshtein = d
                match = matches[i]
              } 
        }
    }
    
	}()
    
	go func() {
		defer wg.Done()

    for i := 0; i < numm; i ++ {
        if (i%2==1) {
              d := levenshtein.Distance(s, matches[i])
              if (d < bestLevenshtein) {
                bestLevenshtein = d
                match = matches[i]
              } 
        }
    }
    
	}()
    
    
	wg.Wait()
    
    return match
}


func main() {
   //generateHash("wordlist")
    match := getMatch(os.Args[1])
    fmt.Printf("Match: %v\n",match)
}
