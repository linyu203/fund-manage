
package main

import (
    "fmt"
)


func main(){
    InitTable();
    /*err := removeEmptyFunds()
    if err != nil {
	fmt.Printf("%v", err)
	return 
    }*/

    CreateFund("testttFund2", "A test fund")
    //InsertBond("testttFund2", "IBM 4.6 04/20/2020 Corp")
    funds,err := GetAllfunds();
    if err != nil {
	fmt.Printf("error: %v", err)
	return 
    }
    for idx,fund := range(funds) {
	fmt.Println( idx, fund);
    }

}
