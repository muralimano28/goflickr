package main
 
import (
    "fmt"
    "net/http"
    "io/ioutil"
    "os"
    "encoding/json"
    )

type Details struct{
	Id string	`json:"id"`
	Owner string	`json:"owner"`
	Secret string	`json:"secret"`
	Server string	`json:"server"`
	Farm int	`json:"farm"`
	Title string 	`json:"title"`
	Ispublic int	`json:"ispublic"`
	Isfriend int	`json:"isfriend"`
	Isfamily int	`json:"isfamily"`
}
type Packet struct{
	Page int	`json:"page"`
	Pages int	`json:"pages"`
	Perpage int	`json:"perpage"`
	Total string	`json:"total"`
	Photo []Details	`json:"photo"`
} 
type Message struct{
	Photos Packet	`json:"photos"`
	Stat string	`json:"stat"`
}
func main() {

	//getting the data from flickr api using get method.....

    response, err := http.Get("https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=c2ef6776bc7946440b01cc9070d55ac0&tags=cute+puppy&per_page=1&format=json&nojsoncallback=1&auth_token=72157654809259522-833770d09d600311&api_sig=164482f059ef2a032b504551e720b66c")
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }

	m := Message{}
	errn := json.Unmarshal([]byte(contents),&m)
	if errn != nil {
		fmt.Println(errn);
	}
	arrayofjsondata := m.Photos.Photo;
	//fmt.Printf("%+v",arrayofjsondata);
	for i,v := range arrayofjsondata {
		fmt.Printf("arr[%d] --> %+v\n",i,v);	
	}
   }
}
